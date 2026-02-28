package sites

import "net/http"

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	success := r.URL.Query().Get("subscribed") == "true"

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

	content := `{{define "content"}}
    
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
            <a href="#" class="hover:underline underline-offset-4">World</a>
            <a href="#" class="hover:underline underline-offset-4">Politics</a>
            <a href="#" class="hover:underline underline-offset-4">Business</a>
            <a href="#" class="hover:underline underline-offset-4">Tech</a>
            <a href="#" class="hover:underline underline-offset-4">Opinion</a>
        </div>
    </header>

    <main class="w-full max-w-screen-2xl mx-auto min-h-screen border-l-[2px] border-r-[2px] border-[#111] bg-[#ededed]">
        
        {{if .Query}}
        <div class="p-8 border-b-[2px] border-[#111]">
            <h2 class="text-4xl font-bold uppercase mb-2">Results for: <span class="bg-[#111] text-[#ededed] px-2">{{.Query}}</span></h2>
            <p class="text-xl font-medium mt-8 mb-8">No articles found matching your criteria. Please refine your search query.</p>
        </div>
        {{else}}
        
        <div class="grid grid-cols-1 lg:grid-cols-12 brutalist-border-b">
            
            <!-- Lead Story -->
            <article class="lg:col-span-8 brutalist-border-r border-b-[2px] lg:border-b-0 border-[#111] article-card p-6 md:p-10 flex flex-col justify-between">
                <div>
                    <div class="flex gap-4 items-center mb-6">
                        <span class="bg-red-600 text-white text-xs font-bold uppercase tracking-widest px-2 py-1 relative -rotate-2">Editor's Pick</span>
                        <span class="text-sm font-bold uppercase tracking-widest text-gray-500">10 Min Read</span>
                    </div>
                    <h2 class="text-5xl md:text-7xl font-black leading-[1.05] tracking-tight mb-8 uppercase">The Automation Age:<br>How AI is Rewriting the Rules of Global Commerce.</h2>
                    <img src="https://images.unsplash.com/photo-1485827404703-89b55fcc595e?q=80&w=1200" class="w-full aspect-video object-cover border-[2px] border-[#111] mb-8 shadow-[8px_8px_0_0_#111]">
                </div>
                <p class="text-xl md:text-2xl font-medium leading-relaxed max-w-3xl">As neural translation engines reach zero-latency, international markets converge into a single unified storefront. We examine the geopolitical impact of instantaneous global communication.</p>
            </article>

            <!-- Sidebar Stories -->
            <aside class="lg:col-span-4 flex flex-col">
                <article class="article-card p-6 md:p-8 border-b-[2px] border-[#111] flex-1">
                    <span class="text-xs font-bold uppercase tracking-widest text-[#111] border-[1px] border-[#111] px-2 py-0.5 rounded-full mb-4 inline-block">Markets</span>
                    <h3 class="text-2xl md:text-3xl font-bold leading-tight uppercase mb-4">Federal Reserve Signals Unprecedented Shift in Monetary Policy</h3>
                    <p class="text-md font-medium text-gray-500 line-clamp-3">Inflationary pressures ease as automated supply chains optimize routing globally.</p>
                </article>
                <article class="article-card p-6 md:p-8 flex-1">
                    <span class="text-xs font-bold uppercase tracking-widest text-[#111] border-[1px] border-[#111] px-2 py-0.5 rounded-full mb-4 inline-block">Culture</span>
                    <h3 class="text-2xl md:text-3xl font-bold leading-tight uppercase mb-4">The Brutalist Revival in Digital Architecture</h3>
                    <p class="text-md font-medium text-gray-500 line-clamp-3">Why modern designers are rejecting soft gradients for harsh lines and high contrast.</p>
                </article>
            </aside>

        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 brutalist-border-b bg-[#111] text-[#ededed]">
            <div class="p-8 md:p-12 md:col-span-2 border-b-[2px] md:border-b-0 md:border-r-[2px] border-[#ededed]">
                 <h3 class="text-4xl font-black uppercase mb-4 text-[#fff]">Newsletter</h3>
                 <p class="text-xl mb-8 font-medium">Get the unvarnished truth delivered twice daily.</p>
                 {{if .Success}}
                    <div class="bg-green-400 text-[#111] p-4 font-bold uppercase border-[2px] border-[#ededed] shadow-[4px_4px_0_0_#fff]">
                        Subscription confirmed. Welcome to the resistance.
                    </div>
                 {{else}}
                 <form action="/news" method="GET" class="flex max-w-md">
                     <input type="hidden" name="subscribed" value="true">
                     <input type="email" placeholder="YOUR@EMAIL.COM" class="brutal-input px-4 py-3 w-full font-bold bg-[#ededed] text-[#111]" required>
                     <button type="submit" class="btn-brutal bg-[#ededed] text-[#111] hover:bg-red-600 hover:text-white px-6 py-3 ml-2 border-[#ededed]">Join</button>
                 </form>
                 {{end}}
            </div>
            <div class="p-8 md:p-12 flex items-center justify-center">
                <div class="w-24 h-24 border-8 border-current rounded-full flex items-center justify-center font-black text-4xl">C</div>
            </div>
        </div>

        {{end}}
    </main>
    {{end}}`

	Render(w, head, content, PageData{
		Brand:       "Chronicle",
		Title:       "Global News Network",
		Description: "The latest breaking news, brutal style.",
		Query:       query,
		Success:     success,
	})
}
