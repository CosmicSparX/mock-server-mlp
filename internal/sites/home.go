package sites

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	head := `{{define "head"}}
    <style>
        body { 
            font-family: 'Inter', sans-serif; 
            background-color: #000; 
            color: #fff; 
            overflow-x: hidden;
            margin: 0;
            padding: 0;
        }
        
        .grid-bg {
            position: fixed;
            top: 0; left: 0; width: 100vw; height: 100vh;
            background-size: 50px 50px;
            background-image: 
                linear-gradient(to right, rgba(255,255,255,0.05) 1px, transparent 1px),
                linear-gradient(to bottom, rgba(255,255,255,0.05) 1px, transparent 1px);
            z-index: -2;
        }
        
        .glow {
            position: fixed;
            top: 50%; left: 50%; width: 600px; height: 600px;
            background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(0,0,0,0) 70%);
            transform: translate(-50%, -50%);
            z-index: -1;
            pointer-events: none;
        }

        .site-card {
            position: relative;
            background: rgba(20, 20, 20, 0.6);
            backdrop-filter: blur(12px);
            -webkit-backdrop-filter: blur(12px);
            border: 1px solid rgba(255,255,255,0.1);
            transition: all 0.5s cubic-bezier(0.165, 0.84, 0.44, 1);
            overflow: hidden;
            display: flex;
            flex-direction: column;
        }
        
        .site-card::before {
            content: '';
            position: absolute;
            top: 0; left: 0; width: 100%; height: 100%;
            background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
            opacity: 0;
            transition: opacity 0.5s ease;
        }

        .site-card:hover {
            transform: translateY(-8px);
            border-color: rgba(255,255,255,0.3);
            box-shadow: 0 20px 40px rgba(0,0,0,0.5);
        }
        
        .site-card:hover::before {
            opacity: 1;
        }

        .site-card::after {
            content: '↗';
            position: absolute;
            top: 1.5rem; right: 1.5rem;
            font-size: 1.5rem;
            opacity: 0;
            transform: translate(-10px, 10px);
            transition: all 0.4s ease;
        }

        .site-card:hover::after {
            opacity: 1;
            transform: translate(0, 0);
        }

        .pill {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 999px;
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            font-weight: 600;
            margin-bottom: 1rem;
        }
        
        .pill-saas { background: rgba(59,130,246,0.2); color: #60a5fa; border: 1px solid rgba(59,130,246,0.3); }
        .pill-shop { background: rgba(217,119,6,0.2); color: #fbbf24; border: 1px solid rgba(217,119,6,0.3); }
        .pill-news { background: rgba(220,38,38,0.2); color: #f87171; border: 1px solid rgba(220,38,38,0.3); }

        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .animate-up { animation: fadeInUp 1s cubic-bezier(0.16, 1, 0.3, 1) forwards; opacity: 0; }
        .delay-100 { animation-delay: 100ms; }
        .delay-200 { animation-delay: 200ms; }
        .delay-300 { animation-delay: 300ms; }
    </style>
    {{end}}`

	content := `{{define "content"}}
    <div class="grid-bg"></div>
    <div class="glow" id="mouse-glow"></div>

    <main class="min-h-screen flex flex-col items-center justify-center p-6 md:p-12 relative z-10">
        
        <header class="text-center mb-20 animate-up">
            <div class="w-16 h-16 border-2 border-white rounded-xl mx-auto flex items-center justify-center mb-8 rotate-45 transform transition hover:rotate-90 duration-500">
                <div class="w-8 h-8 border-2 border-transparent bg-white rounded-sm -rotate-45"></div>
            </div>
            <h1 class="text-6xl md:text-8xl font-black tracking-tighter mb-4">mock.csx.codes<span class="text-gray-500">.</span></h1>
            <p class="text-xl text-gray-400 font-light tracking-wide max-w-lg mx-auto">
                A staging environment 
            </p>
        </header>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full max-w-6xl mx-auto">
            
            <a href="/saas" class="site-card p-10 rounded-2xl animate-up delay-100 group">
                <span class="pill pill-saas">B2B SaaS</span>
                <h2 class="text-3xl font-bold mb-3 tracking-tight">Lumina.ai</h2>
                <p class="text-gray-400 leading-relaxed text-sm mb-12 flex-grow">
                    Enterprise dark-mode aesthetic. Features robust session-based authentication, interactive forms, and glassmorphism.
                </p>
                <div class="mt-auto flex gap-2">
                    <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                    <div class="w-2 h-2 rounded-full bg-indigo-500 opacity-50"></div>
                </div>
            </a>

            <a href="/shop" class="site-card p-10 rounded-2xl animate-up delay-200 group">
                <span class="pill pill-shop">E-Commerce</span>
                <h2 class="text-3xl font-bold mb-3 tracking-tight font-serif">Artisan Brews</h2>
                <p class="text-gray-400 leading-relaxed text-sm mb-12 flex-grow">
                    Kinfolk-inspired editorial layout. Features asymmetric product grids, session carts, and elegant serif typography.
                </p>
                <div class="mt-auto flex gap-2">
                    <div class="w-2 h-2 rounded-full bg-amber-700"></div>
                    <div class="w-2 h-2 rounded-full bg-stone-500 opacity-50"></div>
                </div>
            </a>

            <a href="/news" class="site-card p-10 rounded-2xl animate-up delay-300 group">
                <span class="pill pill-news">Media Portal</span>
                <h2 class="text-3xl font-black mb-3 tracking-tight uppercase">The Chronicle</h2>
                <p class="text-gray-400 leading-relaxed text-sm mb-12 flex-grow">
                    Swiss brutalist design. Features high-density data, dynamic category filtering, database search, and marquees.
                </p>
                <div class="mt-auto flex gap-2">
                    <div class="w-2 h-2 rounded-full bg-red-600"></div>
                    <div class="w-2 h-2 rounded-full bg-gray-500 opacity-50"></div>
                </div>
            </a>

        </div>

    </main>

    <footer class="fixed bottom-0 w-full p-6 text-center text-xs text-gray-600 font-mono tracking-widest z-10 animate-up delay-300">
        // SYSTEM ONLINE &bull; PORT 9090 &bull; SQLITE CONNECTED //
    </footer>

    <script>
        document.addEventListener('mousemove', (e) => {
            const glow = document.getElementById('mouse-glow');
            if(glow) {
                // Smooth follow
                glow.style.left = e.clientX + 'px';
                glow.style.top = e.clientY + 'px';
                glow.style.transform = 'translate(-50%, -50%)';
            }
        });
    </script>
    {{end}}`

	Render(w, head, content, PageData{
		Title:       "Directory",
		Brand:       "Origin",
		Description: "Mock origin server directory.",
	})
}
