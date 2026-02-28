package sites

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mock_server/internal/db"
)

type Article struct {
	ID          int
	Title       string
	Content     string
	Category    string
	ImageURL    string
	PublishedAt time.Time
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	success := r.URL.Query().Get("subscribed") == "true"

	if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/news/newsletter") {
		email := r.FormValue("email")
		if email != "" {
			db.DB.Exec("INSERT INTO newsletter_subs (email) VALUES (?)", email)
			http.Redirect(w, r, "/news?subscribed=true", http.StatusSeeOther)
			return
		}
	}

	head := `{{define "head"}}
    <style>
        body { font-family: 'Space Grotesk', sans-serif; background-color: #ededed; color: #111; overflow-x: hidden; }
        
        /* Brutalist Grid & Borders */
        .brutalist-container { border: 2px solid #111; }
        .brutalist-border-b { border-bottom: 2px solid #111; }
        .brutalist-border-r { border-right: 2px solid #111; }
        
        .ticker-wrap {
            width: 100%; overflow: hidden; height: 3rem; background-color: #111; 
            padding-left: 100%; box-sizing: content-box; border-bottom: 2px solid #111;
        }
        .ticker {
            display: inline-block; height: 3rem; line-height: 3rem; 
            white-space: nowrap; padding-right: 100%; box-sizing: content-box;
            animation-iteration-count: infinite; animation-timing-function: linear;
            animation-name: ticker; animation-duration: 30s;
        }
        .ticker__item { display: inline-block; padding: 0 2rem; font-weight: 700; color: #fff; text-transform: uppercase; font-size: 0.875rem; letter-spacing: 0.1em; }
        
        @keyframes ticker {
            0% { transform: translate3d(0, 0, 0); visibility: visible; }
            100% { transform: translate3d(-100%, 0, 0); }
        }

        .marquee-hover:hover .ticker { animation-play-state: paused; }
        
        .article-card { transition: background 0.2s, color 0.2s; cursor: pointer; }
        .article-card:hover { background: #111; color: #ededed; }
        .article-card:hover .text-gray-500 { color: #a1a1aa; }
        .article-card:hover img { filter: grayscale(0%); }
        .article-card img { filter: grayscale(100%); transition: filter 0.4s; }
        
        .btn-brutal {
            background: #111; color: #ededed; text-transform: uppercase; font-weight: 700;
            border: 2px solid #111; position: relative; transition: all 0.1s;
            box-shadow: 4px 4px 0 0 #000;
        }
        .btn-brutal:active { transform: translate(4px, 4px); box-shadow: 0 0 0 0 #000; }
        .btn-brutal:hover { background: #fff; color: #111; }
        
        input.brutal-input {
            background: transparent; border: 2px solid #111; border-radius: 0;
            outline: none; transition: background 0.2s;
        }
        input.brutal-input:focus { background: #fff; }
    </style>
    {{end}}`

	navHTML := `
    <div class="ticker-wrap marquee-hover">
        <div class="ticker">
            <div class="ticker__item">Breaking: Global markets rally as new tech regulations announced</div>
            <div class="ticker__item">Opinion: The future of artificial intelligence in design</div>
            <div class="ticker__item">Tech: Translation models drop latency below 50ms</div>
            <div class="ticker__item">World: Climate summit concludes with historic agreements</div>
        </div>
    </div>

    <header class="w-full border-b-[2px] border-[#111] bg-[#ededed] relative z-10 pt-8 pb-4 px-4 md:px-8">
        <div class="max-w-screen-2xl mx-auto flex flex-col md:flex-row justify-between items-end gap-6 border-b-[2px] border-[#111] pb-6 mb-2">
            <div>
                <p class="text-xs uppercase tracking-widest font-bold mb-2">Est. 2026 // Vol. 42</p>
                <h1 class="text-6xl md:text-8xl font-black uppercase tracking-tighter leading-none"><a href="/news">The Chronicle</a></h1>
            </div>
            
            <form action="/news" method="GET" class="flex w-full md:w-auto">
                <input type="text" name="q" placeholder="Type to search..." class="brutal-input px-4 py-2 w-full md:w-64 font-bold placeholder-gray-500" value="{{.Query}}">
                <button type="submit" class="btn-brutal px-6 py-2 ml-2">Search</button>
            </form>
        </div>
        <div class="max-w-screen-2xl mx-auto flex gap-6 text-sm font-bold uppercase tracking-widest overflow-x-auto whitespace-nowrap pt-2">
            <a href="/news" class="hover:underline underline-offset-4">Front Page</a>
            <a href="/news/category/World" class="hover:underline underline-offset-4">World</a>
            <a href="/news/category/Tech" class="hover:underline underline-offset-4">Tech</a>
            <a href="/news/category/Markets" class="hover:underline underline-offset-4">Markets</a>
            <a href="/news/category/Culture" class="hover:underline underline-offset-4">Culture</a>
        </div>
    </header>`

	footerHTML := `
		<div class="grid grid-cols-1 md:grid-cols-3 brutalist-border-b bg-[#111] text-[#ededed]">
            <div class="p-8 md:p-12 md:col-span-2 border-b-[2px] md:border-b-0 md:border-r-[2px] border-[#ededed]">
                 <h3 class="text-4xl font-black uppercase mb-4 text-[#fff]">Newsletter</h3>
                 <p class="text-xl mb-8 font-medium">Get the unvarnished truth delivered twice daily.</p>
                 {{if .Success}}
                    <div class="bg-green-400 text-[#111] p-4 font-bold uppercase border-[2px] border-[#ededed] shadow-[4px_4px_0_0_#fff]">
                        Subscription confirmed. Welcome to the resistance.
                    </div>
                 {{else}}
                 <form action="/news/newsletter" method="POST" class="flex max-w-md">
                     <input type="email" name="email" placeholder="YOUR@EMAIL.COM" class="brutal-input px-4 py-3 w-full font-bold bg-[#ededed] text-[#111]" required>
                     <button type="submit" class="btn-brutal bg-[#ededed] text-[#111] hover:bg-red-600 hover:text-white px-6 py-3 ml-2 border-[#ededed]">Join</button>
                 </form>
                 {{end}}
            </div>
            <div class="p-8 md:p-12 flex items-center justify-center">
                <div class="w-24 h-24 border-8 border-current rounded-full flex items-center justify-center font-black text-4xl">C</div>
            </div>
        </div>`

	var content string

	if strings.HasPrefix(r.URL.Path, "/news/article/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/news/article/")
		var a struct {
			ID       int
			Title    string
			Content  string
			Category string
			Author   string
			Date     string
			ImageURL string
		}

		err := db.DB.QueryRow("SELECT id, title, content, category, author, created_at, image_url FROM articles WHERE id = ?", idStr).Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Author, &a.Date, &a.ImageURL)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// LIMIT TEST: Fragmented Text Nodes & Massive Segment Count
		// Replace standard article content with 1000+ paragraphs of alternating <span> tags.
		var builder strings.Builder

		authenticSentences := []string{
			"This structural shift forces immediate recalibration of global supply chains.",
			"Analysts project that within the next fiscal quarter, these pressures will compound.",
			"Preparations must begin to mitigate systemic shock across all integrated sectors.",
			"Governing bodies remain deadlocked as market mechanisms take over.",
			"The implications for downstream logistics providers are unprecedented.",
			"Raw compute is rapidly becoming the most heavily contested resource.",
			"Interest rates could remain steady, defying earlier consensus estimates.",
			"A sober look at the data reveals underlying fragility in the current model.",
			"Silicon valleys are turning to drastic measures to secure manufacturing time.",
			"The sheer volume of algorithmic trading has amplified the recent volatility.",
		}

		for i := 0; i < 1500; i++ {
			sentence := authenticSentences[i%len(authenticSentences)]
			builder.WriteString("<p class=\"mb-6\">")
			for _, char := range sentence {
				builder.WriteString(fmt.Sprintf("<span>%c</span>", char))
			}
			builder.WriteString("</p>")
		}
		a.Content = builder.String()

		// Fetch 3 related articles from the same category
		rows, _ := db.DB.Query("SELECT id, title, excerpt, created_at FROM articles WHERE category = ? AND id != ? LIMIT 3", a.Category, a.ID)
		defer rows.Close()
		var relatedHTML string
		for rows.Next() {
			var rID int
			var rTitle, rExcerpt, rDate string
			rows.Scan(&rID, &rTitle, &rExcerpt, &rDate)
			relatedHTML += `
			<a href="/news/article/` + strconv.Itoa(rID) + `" class="block group border-t border-stone-200 py-6 mb-4">
				<div class="flex justify-between items-center mb-3">
					<h4 class="font-bold text-lg group-hover:text-red-700 transition-colors">` + rTitle + `</h4>
					<svg class="w-5 h-5 opacity-0 -translate-x-4 group-hover:opacity-100 group-hover:translate-x-0 transition-all text-red-700" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
				</div>
				<p class="text-stone-600 text-sm font-serif line-clamp-2">` + rExcerpt + `</p>
			</a>`
		}

		if relatedHTML == "" {
			relatedHTML = `<p class="text-stone-500 italic font-serif">No related coverage found at this time.</p>`
		}

		content = `{{define "content"}}` + navHTML + `
		<main class="pt-32 pb-20 px-4 md:px-8 max-w-5xl mx-auto flex flex-col md:flex-row gap-12 lg:gap-24 relative fade-up">
			
			<article class="md:w-2/3">
				<div class="mb-12">
					<div class="flex items-center gap-4 text-xs font-bold uppercase tracking-widest text-red-700 mb-6">
						<a href="/news/category/` + a.Category + `" class="hover:underline">` + a.Category + `</a>
						<span class="text-stone-300">&bull;</span>
						<span class="text-stone-500 font-mono">` + strings.Split(a.Date, "T")[0] + `</span>
					</div>
					<h1 class="text-5xl md:text-7xl font-black tracking-tighter leading-tight mb-8">` + a.Title + `</h1>
					<div class="flex items-center gap-4 border-y border-stone-200 py-4">
						<div class="w-12 h-12 bg-stone-200"></div>
						<div>
							<div class="font-bold text-sm">By ` + a.Author + `</div>
							<div class="font-mono text-xs text-stone-500 uppercase">Senior Correspondent</div>
						</div>
					</div>
				</div>

				<img src="` + a.ImageURL + `" class="w-full h-[500px] object-cover bg-stone-200 mb-12 grayscale">

				<div class="prose prose-stone prose-lg max-w-none font-serif text-stone-800 leading-loose text-justify">
					<p class="text-2xl font-sans font-light leading-snug text-stone-600 mb-10">We are witnessing a structural shift that demands immediate analysis and unvarnished reporting. The following brief outlines the core components of the developing situation.</p>
					` + a.Content + `
				</div>
			</article>

			<aside class="md:w-1/3 relative border-l border-stone-200 pl-8 hidden md:block">
				<div class="sticky top-40">
					<h3 class="font-black text-2xl uppercase tracking-tighter mb-8 flex items-center gap-3">
						<div class="w-3 h-3 bg-red-600"></div>
						Related Coverage
					</h3>
					
					` + relatedHTML + `

					<div class="mt-16 bg-stone-100 p-8">
						<h4 class="font-bold text-lg mb-4">Daily Briefing</h4>
						<p class="text-sm font-serif text-stone-600 mb-6 line-clamp-3">Subscribe to receive unvarnished, raw intelligence directly to your terminal every morning at 06:00 EST.</p>
						<form method="POST" action="/news/newsletter">
							<input type="email" name="email" placeholder="YOUR EMAIL" class="w-full bg-white border border-stone-300 p-3 mb-3 text-xs font-mono uppercase focus:outline-none focus:border-red-600 focus:ring-1 focus:ring-red-600">
							<button type="submit" class="w-full bg-black text-white px-6 py-3 text-xs font-bold font-mono tracking-widest uppercase hover:bg-red-700 transition-colors">Subscribe</button>
						</form>
					</div>
				</div>
			</aside>

		</main>
		{{end}}`

		// LIMIT TEST: Delayed Origin Flush (TTFB Trap)
		// We execute the template to a string buffer manually so we can flush part of it, sleep, and flush the rest.
		data := PageData{
			Brand:       "Chronicle",
			Title:       a.Title,
			Description: "The latest breaking news, brutal style.",
			Query:       query,
			Success:     success,
		}

		tmpl, _ := template.New("layout").Parse(Layout)
		tmpl.New("head").Parse(head)
		tmpl.New("content").Parse(content)

		var buf bytes.Buffer
		tmpl.Execute(&buf, data)
		fullHTML := buf.Bytes()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if len(fullHTML) > 100 {
			w.Write(fullHTML[:100])
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(500 * time.Millisecond) // Intentionally lag TTFB by 500ms
			w.Write(fullHTML[100:])
		} else {
			w.Write(fullHTML)
		}
		return

	} else {
		// Category or Home
		var rows *sql.Rows
		var err error
		var headerTitle string

		if strings.HasPrefix(r.URL.Path, "/news/category/") {
			category := strings.TrimPrefix(r.URL.Path, "/news/category/")
			headerTitle = "CATEGORY: " + category
			rows, err = db.DB.Query("SELECT id, title, excerpt, category, author, created_at, image_url FROM articles WHERE category = ? ORDER BY id DESC LIMIt 20", category)
		} else if query != "" {
			headerTitle = "SEARCH: " + query
			rows, err = db.DB.Query("SELECT id, title, excerpt, category, author, created_at, image_url FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY id DESC LIMIT 20", "%"+query+"%", "%"+query+"%")
		} else {
			headerTitle = "LATEST WIRE"
			rows, err = db.DB.Query("SELECT id, title, excerpt, category, author, created_at, image_url FROM articles ORDER BY id DESC LIMIT 20")
		}

		var articlesHTML string
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id int
				var title, excerpt, category, author, date, img string
				rows.Scan(&id, &title, &excerpt, &category, &author, &date, &img)

				formattedDate := strings.Split(date, "T")[0]

				articlesHTML += `
				<a href="/news/article/` + strconv.Itoa(id) + `" class="block group border-b-[2px] border-[#111] p-8 md:p-12 hover:bg-[#111] hover:text-[#ededed] transition-colors duration-300">
					<div class="flex flex-col md:flex-row gap-8 lg:gap-16 items-start md:items-center">
						<div class="flex-1">
							<div class="flex items-center gap-4 mb-4">
								<span class="text-xs font-bold uppercase tracking-widest text-red-600 bg-red-100 group-hover:bg-red-900 px-2 py-1">` + category + `</span>
								<span class="text-sm font-mono uppercase text-gray-500 group-hover:text-gray-400">` + formattedDate + `</span>
							</div>
							<h3 class="text-3xl md:text-5xl font-black mb-4 leading-none tracking-tight">` + title + `</h3>
							<p class="text-lg font-serif font-light line-clamp-2 md:max-w-3xl opacity-80 group-hover:opacity-100 transition-opacity">` + excerpt + `</p>
							<div class="mt-4 font-mono text-xs uppercase text-gray-500 tracking-widest group-hover:text-gray-400">
								By ` + author + `
							</div>
						</div>
						<div class="w-full md:w-64 aspect-video bg-[#111] group-hover:bg-[#ededed] border-[2px] border-[#111] overflow-hidden shrink-0">
							<img src="` + img + `" class="w-full h-full object-cover grayscale group-hover:grayscale-0 transition-all duration-500 group-hover:scale-105">
						</div>
					</div>
				</a>`
			}
		}

		if articlesHTML == "" {
			articlesHTML = `<div class="p-12 text-center text-xl font-medium font-serif italic text-gray-500">NO TRANSMISSIONS LOCATED.</div>`
		}

		content = `{{define "content"}}` + navHTML + `
		<main class="w-full max-w-screen-2xl mx-auto min-h-screen border-l-[2px] border-r-[2px] border-[#111] bg-[#ededed] mt-48">
			
			<div class="border-b-[2px] border-[#111] bg-[#111] text-white p-4 overflow-hidden relative flex items-center">
				<div class="w-2 h-2 bg-red-600 mr-4 shrink-0 animate-pulse"></div>
				<div class="font-bold uppercase tracking-widest text-sm shrink-0 mr-8">Live Feed</div>
				<marquee class="font-mono text-sm tracking-widest opacity-80" scrollamount="8">
					MARKETS RALLY ON RATE PAUSE // TECH EARNINGS EXCEED EXPECTATIONS // GLOBAL SUPPLY CHAINS REROUTE // QUANTUM BREAKTHROUGH ANNOUNCED IN GENEVA //
				</marquee>
			</div>

			<div class="p-8 md:p-12 border-b-[2px] border-[#111] flex flex-col md:flex-row justify-between md:items-end gap-8 outline-none focus:bg-yellow-100 transition-colors" tabindex="0">
				<div>
					<h2 class="text-6xl md:text-8xl font-black leading-[1.05] tracking-tighter uppercase">` + headerTitle + `</h2>
				</div>
				<form method="GET" action="/news" class="flex w-full md:w-auto">
					<input type="text" name="q" placeholder="SEARCH ARCHIVE..." class="bg-transparent border-b-[2px] border-[#111] text-xl font-bold uppercase placeholder-[#111]/50 focus:outline-none p-2 w-full md:w-64 input-focus">
					<button type="submit" class="p-2 border-b-[2px] border-[#111] hover:bg-[#111] hover:text-white transition-colors">
						<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
					</button>
				</form>
			</div>

			<div class="divide-y-[2px] divide-[#111]">
				` + articlesHTML + `
			</div>
			
			` + footerHTML + `
		</main>
		{{end}}`
	}

	Render(w, head, strings.Replace(content, "%%", "%", -1), PageData{
		Brand:       "Chronicle",
		Title:       "Global News Network",
		Description: "The latest breaking news, brutal style.",
		Query:       query,
		Success:     success,
	})
}
