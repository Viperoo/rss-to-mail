package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	rss "github.com/ungerik/go-rss"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
	"time"
)

type MessageTemplate struct {
	Website Website
	Item    rss.Item
	Email   string
	Time    string
	Subject string
}

func prepareMessage(item rss.Item, website Website) (subject, body string) {

	t, err := template.New("feed").Parse(getMailTemplate())
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer

	now := time.Now().Format(TIME_FORMAT)
	subject = "[RSS][" + website.Title + "] " + item.Title
	err = t.Execute(&doc, MessageTemplate{
		Item:    item,
		Website: website,
		Email:   Config.Default.To,
		Time:    now,
		Subject: subject,
	})
	if err != nil {
		panic(err)
	}
	body = doc.String()

	return subject, body
}
func sendEmial(subject, body string) {
	headers := make(map[string]string)
	headers["From"] = Config.SMTP.From
	headers["To"] = Config.Default.To
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", Config.SMTP.User, Config.SMTP.Password, Config.SMTP.Host)

	tlc := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         Config.SMTP.Host,
	}
	c, err := smtp.Dial(Config.SMTP.Host + ":" + strconv.Itoa(Config.SMTP.Port))
	if err != nil {
		log.Panic(err)
	}

	if err = c.StartTLS(tlc); err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(Config.SMTP.From); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(Config.Default.To); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}
