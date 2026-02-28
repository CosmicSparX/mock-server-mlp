package sites

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Brand       string
	Title       string
	Description string
	ActivePage  string
	CartCount   int
	Query       string
	Success     bool
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
    {{template "head" .}}
</head>
<body class="antialiased selection:bg-indigo-500 selection:text-white">
    {{template "content" .}}
</body>
</html>
`

func Render(w http.ResponseWriter, headTmpl, contentTmpl string, data PageData) {
	tmpl, _ := template.New("layout").Parse(Layout)
	if headTmpl != "" {
		tmpl.New("head").Parse(headTmpl)
	} else {
		tmpl.New("head").Parse("")
	}
	tmpl.New("content").Parse(contentTmpl)
	tmpl.Execute(w, data)
}
