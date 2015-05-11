package main

func getMailTemplate() string {
	return `
<html>
<head>
<title> {{.Subject}}</title>
<meta charset="UTF-8">
<style>
body {
	font-family: monospace
}
dl {
    margin-top: 0;
    margin-bottom: 20px
}

dt {
    font-weight: 700
}

dd {
    margin-left: 0
}
@media (min-width: 768px) {
    .dl-horizontal dt {
        float: left;
        width: 120px;
        overflow: hidden;
        clear: left;
        text-align: right;
        text-overflow: ellipsis;
        white-space: nowrap
    }

    .dl-horizontal dd {
        margin-left: 125px
    }
}
 .dl-horizontal dd:after, .dl-horizontal dd:before {
	  display: table;
    content: " "
}
.dl-horizontal dd:after {
	clear: both;
}
img {
	max-width: 300px;
	max-height: 150px
	}
</style>
</head>
<body>

<dl class="dl-horizontal">
	<dt>Description:</dt>
</dl>
	<p style="clear:both">
	{{ .Item.Description}}
</p>
<dl class="dl-horizontal">

	<dt>Pub date:</dt>
    <dd>{{ .Time }}</dd>
    <dt>Website:</dt>
    <dd>{{ .Website.Title }}</dd>
    <dt>Title:</dt>
	<dd>{{ .Item.Title }}</dd>
	<dt>Link:</dt>
	<dd>{{ .Item.Link }} [<a href="{{ .Item.Link }}">link</a>]</dd>
	
	<dt>See more:</dt>
    <dd><a href="{{ .Item.Link}}">Here</a> </dd>
</dl>

</body>
</html>`
}
