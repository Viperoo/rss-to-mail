package main

import (
	"flag"
	"fmt"
	"github.com/Viperoo/golog"
	"github.com/jinzhu/gorm"
	rss "github.com/ungerik/go-rss"
	"io"
	"os"
	"os/user"
)

var logger log.Logger
var workdir string

var url = flag.String("a", "", "Add website [url]")
var remove = flag.String("r", "", "Remove website [url]")
var list = flag.Bool("l", false, "List added websites")
var check = flag.Bool("c", false, "Check for updates and send mail")
var migrate = flag.Bool("m", false, "Migrate database")

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: rss-to-mail [options] paramflag>\n\n")
		flag.PrintDefaults()
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
	}
}

const (
	TIME_FORMAT = "Mon Jan 02 15:04:05 MST 2006"
)

func main() {
	flag.Parse()

	setWorkDir()
	setLogger()
	db := conn()
	db.LogMode(false)
	if *migrate != false {
		makeMigrate(db)
	}

	if *url != "" {
		addFeed(db, *url)
	}

	if *list != false {
		showList(db)
	}

	if *check != false {
		checkForUpdates(db)
	}

	if *remove != "" {
		removeWebsite(db, *remove)
	}

}

func addFeed(db gorm.DB, url string) {
	logger.Info("Adding url website: " + url)
	feed, err := rss.Read(url)

	if err != nil {
		logger.Warn("Add url failed")
	} else {
		var count int64
		db.Model(Website{}).Where("url = ?", url).Count(&count)

		if count > 0 {
			logger.Warnf("There is website with %v url", url)
		} else {
			website := Website{Title: feed.Title, Url: url, Description: feed.Description, Language: feed.Language}
			db.Save(&website)
			logger.Infof("Added with id: %v ("+feed.Title+")", website.Id)

			for _, item := range feed.Item {

				i := Feed{WebsiteId: website.Id, GUID: item.GUID}
				db.Save(&i)
			}
		}
	}
}

func checkForUpdates(db gorm.DB) {
	logger.Info("Start checking for updates:")

	var websites []Website
	db.Find(&websites)

	for _, website := range websites {
		feed, err := rss.Read(website.Url)
		if err == nil {
			for _, item := range feed.Item {
				var count int64
				db.Model(Feed{}).Where("g_u_i_d = ?", item.GUID).Count(&count)

				if count == 0 {
					logger.Info("New item (" + item.GUID + ") Send e-mail...")

					subject, body := prepareMessage(item, website)
					sendEmial(subject, body)

					i := Feed{WebsiteId: website.Id, GUID: item.GUID}
					db.Save(&i)
				}
			}
		}

	}

}

func showList(db gorm.DB) {
	logger.Info("List added websites:")

	var websites []Website
	db.Find(&websites)
	fmt.Printf("Id.\tTitle\n")
	for _, website := range websites {
		fmt.Printf("%v\t%v\n", website.Id, website.Title)
	}
}

func removeWebsite(db gorm.DB, url string) {
	var website Website
	db.Where("url = ?", url).First(&website)
	if website.Id == 0 {
		logger.Warnf("There is not website with %v url", url)
	} else {
		db.Delete(&website)
		logger.Infof("Deleted website with %v id", website.Id)
	}
}

func setWorkDir() {
	usr, err := user.Current()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	workdir = usr.HomeDir + "/.rss-to-mail/"
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		os.Mkdir(workdir, 0777)
		os.Mkdir(workdir+"db", 0777)

		if _, err := os.Stat(workdir + "configuration.conf"); err != nil {
			file, err := os.OpenFile(workdir+"configuration.conf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			defer file.Close()

			var d string = `[Default]
To = your-email@localhost
[SMTP]
Host = localhost
Port = 25
User = smtp@localhost
Password = 
From = rss@localhost
`

			if _, err = file.WriteString(d); err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			logger.Info("Default configuration loaded")
		}

	}
	loadConfig(workdir + "configuration.conf")

}

func setLogger() {
	file, err := os.OpenFile(workdir+"rss-to-mail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var multi io.Writer
	multi = io.MultiWriter(file, os.Stdout)
	logger, _ = log.NewLogger(multi,
		log.TIME_FORMAT_SEC,   // Set time writting format.
		log.LOG_FORMAT_SIMPLE, // Set log writting format.
		log.LogLevel_Debug)    // Set log level.
}
