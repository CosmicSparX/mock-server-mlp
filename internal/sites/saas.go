package sites

import "net/http"

func SaaSHandler(w http.ResponseWriter, r *http.Request) {
	// Handle Form Submission
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if email != "" {
			http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
			return
		}
	}

	success := r.URL.Query().Get("success") == "true"

	head := `{{define "head"}}
    <style>
        body { font-family: 'Inter', sans-serif; background-color: #050505; color: #ededed; overflow-x: hidden; }
        .glass-nav {
            background: rgba(10, 10, 10, 0.6);
            backdrop-filter: blur(16px);
            -webkit-backdrop-filter: blur(16px);
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        .hero-glow {
            position: absolute;
            top: -20%; left: 50%;
            width: 800px; height: 800px;
            background: radial-gradient(circle, rgba(59,130,246,0.15) 0%, rgba(0,0,0,0) 70%);
            transform: translateX(-50%);
            z-index: -1;
            pointer-events: none;
        }
        .btn-primary {
            background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
            box-shadow: 0 4px 14px 0 rgba(37, 99, 235, 0.39);
            transition: all 0.3s ease;
        }
        .btn-primary:hover {
            box-shadow: 0 6px 20px rgba(37, 99, 235, 0.6);
            transform: translateY(-2px);
        }
        .feature-card {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            transition: transform 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275), background 0.4s ease;
        }
        .feature-card:hover {
            transform: translateY(-10px) scale(1.02);
            background: rgba(255, 255, 255, 0.04);
            border-color: rgba(255, 255, 255, 0.1);
        }
        /* Simple fade-in animation for load */
        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .animate-fade-in { animation: fadeInUp 1s cubic-bezier(0.16, 1, 0.3, 1) forwards; opacity: 0; }
        .delay-100 { animation-delay: 100ms; }
        .delay-200 { animation-delay: 200ms; }
        .delay-300 { animation-delay: 300ms; }
    </style>
    {{end}}`

	content := `{{define "content"}}
    <div class="hero-glow"></div>
    
    <nav class="fixed top-0 w-full z-50 glass-nav transition-all duration-300 py-4">
        <div class="max-w-7xl mx-auto px-6 flex justify-between items-center">
            <div class="text-2xl font-bold tracking-tighter flex items-center gap-2">
                <div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center">
                    <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
                </div>
                Lumina.ai
            </div>
            <div class="hidden md:flex space-x-8 text-sm font-medium text-gray-400">
                <a href="#features" class="hover:text-white transition-colors">Features</a>
                <a href="#pricing" class="hover:text-white transition-colors">Pricing</a>
                <a href="#about" class="hover:text-white transition-colors">About</a>
            </div>
            <div class="space-x-4">
                <a href="/login" class="text-sm font-medium text-gray-300 hover:text-white transition">Log in</a>
                <a href="#demo" class="text-sm font-medium px-5 py-2.5 rounded-full bg-white text-black hover:bg-gray-200 transition">Get Started</a>
            </div>
        </div>
    </nav>

    <main class="pt-40 pb-20 px-6 max-w-7xl mx-auto flex flex-col items-center">
        {{if .Success}}
        <div class="animate-fade-in bg-blue-900/30 border border-blue-500/30 text-blue-200 px-6 py-4 rounded-xl mb-12 flex items-center gap-3 backdrop-blur-sm">
            <svg class="w-6 h-6 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            <span class="font-medium">Success! We've sent your priority demo link to your inbox.</span>
        </div>
        {{end}}

        <div class="text-center max-w-4xl animate-fade-in">
            <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-900/20 border border-blue-500/20 text-blue-400 text-sm font-medium mb-8">
                <span class="relative flex h-2 w-2">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
                </span>
                Introducing Translation Agent 2.0
            </div>
            <h1 class="text-6xl md:text-8xl font-bold tracking-tighter mb-8 leading-[1.1]">
                Automate your <br/>
                <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-600">Global Workflows.</span>
            </h1>
            <p class="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed">
                Connect your CMS, databases, and apps to our neural translation engine. Reduce localization time from weeks to seconds with unmatched accuracy.
            </p>
            
            <form method="POST" action="/" id="demo" class="flex flex-col sm:flex-row gap-3 justify-center max-w-md mx-auto delay-100 animate-fade-in">
                <input type="email" name="email" placeholder="Enter work email" class="flex-1 bg-white/5 border border-white/10 text-white px-6 py-4 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all placeholder-gray-500" required>
                <button type="submit" class="btn-primary text-white px-8 py-4 rounded-xl font-medium tracking-wide">Get Demo</button>
            </form>
        </div>

        <div id="features" class="mt-40 grid md:grid-cols-3 gap-6 w-full animate-fade-in delay-200">
            <div class="feature-card p-8 rounded-2xl">
                <div class="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center mb-6 border border-blue-500/20">
                    <svg class="w-6 h-6 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
                </div>
                <h3 class="text-xl font-semibold mb-3">Lightning Fast</h3>
                <p class="text-gray-400 leading-relaxed text-sm">Our edge network processes requests in milliseconds, ensuring your localized content is always available instantly.</p>
            </div>
            <div class="feature-card p-8 rounded-2xl">
                <div class="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center mb-6 border border-purple-500/20">
                    <svg class="w-6 h-6 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"/></svg>
                </div>
                <h3 class="text-xl font-semibold mb-3">Native Nuance</h3>
                <p class="text-gray-400 leading-relaxed text-sm">Advanced LLM pipelines understand context, tone, and brand voice, delivering translations that feel human-authored.</p>
            </div>
            <div class="feature-card p-8 rounded-2xl">
                <div class="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center mb-6 border border-emerald-500/20">
                    <svg class="w-6 h-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>
                </div>
                <h3 class="text-xl font-semibold mb-3">Enterprise Security</h3>
                <p class="text-gray-400 leading-relaxed text-sm">SOC2 Type II certified. Your data is encrypted at rest and in transit. We never use your data to train public models.</p>
            </div>
        </div>
        
        <script>
            // Add a subtle parallax effect to the hero glow based on mouse movement
            document.addEventListener('mousemove', (e) => {
                const glow = document.querySelector('.hero-glow');
                if(!glow) return;
                const x = (e.clientX / window.innerWidth - 0.5) * 40;
                const y = (e.clientY / window.innerHeight - 0.5) * 40;
                glow.style.transform = "translate(calc(-50% + " + x + "px), " + y + "px)";
            });
        </script>
    </main>
    {{end}}`

	Render(w, head, content, PageData{
		Brand:       "Lumina",
		Title:       "Enterprise AI Translation",
		Description: "Automate your global workflow with Lumina AI.",
		Success:     success,
	})
}
