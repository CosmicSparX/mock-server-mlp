package sites

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("awwwards_mock_secret_key_2026"))

type PageData struct {
	Brand       string
	Title       string
	Description string
	ActivePage  string
	CartCount   int
	Query       string
	Success     bool
	UserEmail   string
	IsLoggedIn  bool
}

const Layout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Description}}">
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400;1,600&family=Space+Grotesk:wght@400;500;700&display=swap" rel="stylesheet">
    <title>{{.Title}} | {{.Brand}}</title>
	<link href="https://mock.csx.codes/" hreflang="x-default" rel="alternate">

	<link href="https://hi.mock.csx.codes/" hreflang="hi" rel="alternate">

	<link href="https://ja.mock.csx.codes/" hreflang="ja" rel="alternate">

	<script src="https://script-cdn.multilipi.com/static/JS/page_translations.js" multilipi-key="61d1bfb7-95bd-42b6-a0ad-1b9d0308228c" mode="auto" data-pos-x="50" data-pos-y="50" crossorigin="anonymous" defer>
	</script>
    {{template "head" .}}
</head>
<body class="antialiased selection:bg-indigo-500 selection:text-white">
    {{template "content" .}}
</body>
</html>
`

func Render(w http.ResponseWriter, headTmpl, contentTmpl string, data PageData) {
	tmpl, err := template.New("layout").Parse(Layout)
	if err != nil {
		log.Printf("Layout parse error: %v\n", err)
	}
	if headTmpl != "" {
		_, err = tmpl.New("head").Parse(headTmpl)
		if err != nil {
			log.Printf("Head parse error: %v\n", err)
		}
	} else {
		tmpl.New("head").Parse("")
	}
	_, err = tmpl.New("content").Parse(contentTmpl)
	if err != nil {
		log.Printf("Content parse error: %v\n", err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execute error: %v\n", err)
	}
}
