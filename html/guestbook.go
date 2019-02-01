package html

import (
	"html/template"
)

var GuestBookTmpl = template.Must(template.New("index.html").Parse(`<!DOCTYPE html>
<title>Guestbook - {{.Env}}</title>
<style type="text/css">
html, body {
	font-family: Helvetica, sans-serif;
}
blockquote {
	font-family: cursive, Helvetica, sans-serif;
}
.banner {
	height: 125px;
	width: 250px;
}
.greeting {
	font-size: 85%;
}
.motd {
	font-weight: bold;
}
</style>
<h1>Guestbook</h1>
<div><img class="banner" src="{{.BannerSrc}}"></div>
{{with .MOTD}}<p class="motd">Admin says: {{.}}</p>{{end}}
{{range .Greetings}}
<div class="greeting">
	Someone wrote:
	<blockquote>{{.Content}}</blockquote>
</div>
{{end}}
<form action="/sign" method="POST">
	<div><textarea name="content" rows="3"></textarea></div>
	<div><input type="submit" value="Sign"></div>
</form>
`))
