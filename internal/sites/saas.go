package sites

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mock_server/internal/db"
	"net/http"
	"strings"
)

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func SaaSHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "lumina-session")
	userEmail, isLoggedIn := session.Values["email"].(string)

	success := r.URL.Query().Get("success") == "true"

	switch r.URL.Path {
	case "/saas/api/register":
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			if email != "" && password != "" {
				_, err := db.DB.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, hashPassword(password))
				if err == nil {
					session.Values["email"] = email
					session.Save(r, w)
					http.Redirect(w, r, "/saas?success=true", http.StatusSeeOther)
					return
				}
			}
		}
		http.Redirect(w, r, "/saas/register?error=true", http.StatusSeeOther)
		return

	case "/saas/api/login":
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			var dbHash string
			err := db.DB.QueryRow("SELECT password_hash FROM users WHERE email = ?", email).Scan(&dbHash)
			if err == nil && dbHash == hashPassword(password) {
				session.Values["email"] = email
				session.Save(r, w)
				http.Redirect(w, r, "/saas", http.StatusSeeOther)
				return
			}
		}
		http.Redirect(w, r, "/saas/login?error=true", http.StatusSeeOther)
		return

	case "/saas/api/logout":
		session.Values["email"] = ""
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/saas", http.StatusSeeOther)
		return

	case "/saas/api/contact":
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			if email != "" {
				http.Redirect(w, r, "/saas?success=true", http.StatusSeeOther)
				return
			}
		}
		http.Redirect(w, r, "/saas", http.StatusSeeOther)
		return

	case "/saas/assets/bloated.js":
		w.Header().Set("Content-Type", "application/javascript")
		// Generate an exact 3MB payload string
		chunk := "/* MOCK_3MB_JS_PAYLOAD_CHUNK_DATA_FILLER_STRING_1234567890 */\n"
		// chunk is 63 bytes. 3,145,728 bytes / 63 = 49932 chunks
		payload := strings.Repeat(chunk, 50000)
		w.Write([]byte(payload))
		return
	}

	// 1. LIMIT TEST: Intense Regex Trap (CPU)
	// A massive inline CSS block in the <head> forces "Regex Backtracking" in naive proxy scrapers.
	var cssTrapBuilder strings.Builder
	for i := 0; i < 3000; i++ {
		cssTrapBuilder.WriteString(fmt.Sprintf("\n.obfuscated-class-%d > div:nth-child(%d):hover::after { content: 'TEST_CHUNK_%d'; color: transparent; background-position: %dpx %dpx; opacity: 0.%d; }", i, i%10+1, i, i, i, i%99))
	}
	cssTrap := cssTrapBuilder.String()

	head := `{{define "head"}}
    <style>` + cssTrap + `
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
        .feature-card, .auth-card {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            transition: transform 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275), background 0.4s ease;
        }
        .feature-card:hover {
            transform: translateY(-10px) scale(1.02);
            background: rgba(255, 255, 255, 0.04);
            border-color: rgba(255, 255, 255, 0.1);
        }
        .gradient-border {
            position: relative;
            background: #050505;
            background-clip: padding-box;
            border: 1px solid transparent;
            border-radius: 1.5rem;
        }
        .gradient-border::before {
            content: '';
            position: absolute;
            top: 0; right: 0; bottom: 0; left: 0;
            z-index: -1;
            margin: -1px;
            border-radius: inherit;
            background: linear-gradient(to right, #3b82f6, #8b5cf6, #ec4899);
        }
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

	navHTML := `
    <nav class="fixed top-0 w-full z-50 glass-nav transition-all duration-300 py-4">
        <div class="max-w-7xl mx-auto px-6 flex justify-between items-center">
            <a href="/saas" class="text-2xl font-bold tracking-tighter flex items-center gap-2">
                <div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center">
                    <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
                </div>
                Lumina.ai
            </a>
            <div class="hidden md:flex space-x-8 text-sm font-medium text-gray-400">
                <a href="/saas/features" class="hover:text-white transition-colors">Features</a>
                <a href="/saas/pricing" class="hover:text-white transition-colors">Pricing</a>
                <a href="/saas/about" class="hover:text-white transition-colors">About</a>
            </div>
            <div class="space-x-4">
                {{if .IsLoggedIn}}
                    <a href="/saas/profile" class="text-sm text-gray-400 hover:text-white mr-4 transition">{{.UserEmail}}</a>
                    <a href="/saas/api/logout" class="text-sm font-medium text-gray-300 hover:text-white transition">Log out</a>
                {{else}}
                    <a href="/saas/login" class="text-sm font-medium text-gray-300 hover:text-white transition">Log in</a>
                    <a href="/saas/register" class="text-sm font-medium px-5 py-2.5 rounded-full bg-white text-black hover:bg-gray-200 transition">Get Started</a>
                {{end}}
            </div>
        </div>
    </nav>`

	// 2. LIMIT TEST: Large Hidden Form Payloads (CPU/Parsing)
	// Inject 50+ hidden fields into forms with large string values
	var hiddenInputsBuilder strings.Builder
	payloadStr := strings.Repeat("MOCK_PAYLOAD_CHUNK_X99_", 50)
	for i := 0; i < 60; i++ {
		hiddenInputsBuilder.WriteString(fmt.Sprintf("\n\t\t\t\t\t<input type=\"hidden\" name=\"obfuscated_payload_%d\" value=\"%s\">", i, payloadStr))
	}
	hiddenInputs := hiddenInputsBuilder.String()

	// 3. LIMIT TEST: Heavy Client Logic & 3MB JS file block
	jsTrap := "\n<script src=\"/saas/assets/bloated.js\"></script>\n<script>\n// DOM Heavy Matrix Initialization\nwindow.luminaMatrix = [];\nfor(let i=0; i<30000; i++) { window.luminaMatrix.push('HYPER_SCALE_DATA_CHUNK_' + i + '_' + Math.random().toString(36).substring(7)); }\n</script>"

	var content string

	if r.URL.Path == "/saas/login" {
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-40 pb-20 px-6 w-full max-w-md mx-auto flex flex-col items-center animate-fade-in">
			<div class="auth-card p-10 rounded-3xl w-full text-center">
				<h2 class="text-3xl font-bold mb-2">Welcome Back</h2>
				<p class="text-gray-400 text-sm mb-8">Sign in to continue to Lumina.ai</p>
				{{if .Query}}
					<div class="bg-red-900/40 text-red-200 p-3 rounded-lg mb-6 text-sm border border-red-500/30">Invalid credentials.</div>
				{{end}}
				<form method="POST" action="/saas/api/login" class="flex flex-col gap-4 text-left">
					` + hiddenInputs + `
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Email</label>
						<input type="email" name="email" class="w-full bg-black/50 border border-white/10 text-white px-4 py-3 rounded-xl focus:border-blue-500 focus:outline-none" required>
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Password</label>
						<input type="password" name="password" class="w-full bg-black/50 border border-white/10 text-white px-4 py-3 rounded-xl focus:border-blue-500 focus:outline-none" required>
					</div>
					<button type="submit" class="btn-primary w-full py-3 rounded-xl font-medium mt-4">Sign In</button>
				</form>
				<p class="mt-6 text-sm text-gray-400">Don't have an account? <a href="/saas/register" class="text-blue-400 hover:text-blue-300">Sign up</a></p>
			</div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas/register" {
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-40 pb-20 px-6 w-full max-w-md mx-auto flex flex-col items-center animate-fade-in">
			<div class="auth-card p-10 rounded-3xl w-full text-center">
				<h2 class="text-3xl font-bold mb-2">Create Account</h2>
				<p class="text-gray-400 text-sm mb-8">Start your free 14-day trial</p>
				{{if .Query}}
					<div class="bg-red-900/40 text-red-200 p-3 rounded-lg mb-6 text-sm border border-red-500/30">Registration failed. Email may be taken.</div>
				{{end}}
				<form method="POST" action="/saas/api/register" class="flex flex-col gap-4 text-left">
					` + hiddenInputs + `
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Work Email</label>
						<input type="email" name="email" class="w-full bg-black/50 border border-white/10 text-white px-4 py-3 rounded-xl focus:border-blue-500 focus:outline-none" required>
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Password</label>
						<input type="password" name="password" class="w-full bg-black/50 border border-white/10 text-white px-4 py-3 rounded-xl focus:border-blue-500 focus:outline-none" required>
					</div>
					<button type="submit" class="btn-primary w-full py-3 rounded-xl font-medium mt-4">Sign Up</button>
				</form>
			</div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas/profile" {
		if !isLoggedIn {
			http.Redirect(w, r, "/saas/login", http.StatusSeeOther)
			return
		}
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-40 pb-20 px-6 max-w-7xl mx-auto animate-fade-in text-white w-full">
			<div class="mb-12 border-b border-white/10 pb-8 flex flex-col md:flex-row justify-between md:items-end gap-6">
				<div>
					<h1 class="text-4xl md:text-5xl font-bold tracking-tighter mb-4">Workspace Dashboard</h1>
					<p class="text-gray-400">Manage your enterprise API keys, edge configurations, and billing for <span class="text-blue-400 font-mono">` + userEmail + `</span>.</p>
				</div>
				<button class="btn-primary px-6 py-3 rounded-xl font-medium text-sm">Create New Token</button>
			</div>
			
			<div class="grid md:grid-cols-3 gap-8 mb-12">
				<div class="feature-card p-8 rounded-2xl">
					<h3 class="text-sm uppercase tracking-widest text-gray-500 font-bold mb-2">Total Edge Requests</h3>
					<div class="text-4xl font-black text-white">1,204,992</div>
					<div class="text-sm text-emerald-400 mt-2 flex items-center gap-1">
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
						+14.2% this week
					</div>
				</div>
				<div class="feature-card p-8 rounded-2xl">
					<h3 class="text-sm uppercase tracking-widest text-gray-500 font-bold mb-2">Avg Latency</h3>
					<div class="text-4xl font-black text-white">42ms</div>
					<div class="text-sm text-emerald-400 mt-2 flex items-center gap-1">
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"></path></svg>
						-3ms this week
					</div>
				</div>
				<div class="feature-card p-8 rounded-2xl">
					<h3 class="text-sm uppercase tracking-widest text-gray-500 font-bold mb-2">Active Regions</h3>
					<div class="text-4xl font-black text-white">24 / 24</div>
					<div class="text-sm text-blue-400 mt-2 flex items-center gap-1">All systems operational</div>
				</div>
			</div>

			<div class="bg-white/5 border border-white/10 rounded-3xl p-10 overflow-hidden relative">
				<div class="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-transparent pointer-events-none"></div>
				<h3 class="text-2xl font-bold mb-6">API Keys</h3>
				<div class="w-full bg-black/50 border border-white/10 rounded-xl p-6 flex flex-col md:flex-row justify-between items-start md:items-center mb-4 gap-4">
					<div class="w-full overflow-hidden">
					    <span class="text-xs uppercase tracking-widest text-gray-500 mb-1 block">Production Token</span>
					    <span class="font-mono text-gray-300 block w-full truncate">pk_live_8f92j3f0293jf0293jf029jf0293jf0293jf_xxxxxxx</span>
					</div>
					<button class="bg-white/10 hover:bg-white/20 px-6 py-3 rounded-lg text-sm font-bold transition">Copy</button>
				</div>
				<p class="text-sm text-red-300 font-medium">Never expose your live keys in client-side code.</p>
			</div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas/features" {
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-48 pb-32 px-6 max-w-7xl mx-auto animate-fade-in">
			<div class="text-center mb-24">
			    <h1 class="text-6xl md:text-7xl font-bold tracking-tighter mb-6">Built for the <br/><span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-600">Enterprise Context</span></h1>
			    <p class="text-xl text-gray-400 max-w-3xl mx-auto leading-relaxed">Stop relying on brittle word-to-word dictionaries. Our foundational AI models understand intent, culture, and brand voice natively.</p>
            </div>
            
            <div class="grid md:grid-cols-2 gap-12 lg:gap-24 items-center mb-32">
                <div class="rounded-3xl overflow-hidden border border-white/10 bg-white/5 p-8 relative">
                    <div class="absolute inset-0 bg-gradient-to-br from-blue-500/20 to-purple-500/20 blur-xl"></div>
                    <pre class="relative text-sm text-green-400 font-mono leading-relaxed">
{
  "status": "active",
  "pipeline": "contextual-llm-v4",
  "latency": "22ms",
  "nodes": ["tokyo", "frankfurt", "us-east"],
  "integrity": 99.999
}
                    </pre>
                </div>
                <div>
                    <h2 class="text-4xl font-bold mb-6">Zero-Latency Edge Network</h2>
                    <p class="text-gray-400 text-lg leading-relaxed mb-8">By pushing our inference models directly to the edge, translations happen at the point of request. This means your international customers experience your site exactly as fast as your domestic ones. We deploy parallel inference clusters in major metropolitan zones covering 98% of requested Internet traffic routes. Your Time-to-First-Byte on localized pages will perfectly match your origin server.</p>
                    <ul class="space-y-4 text-gray-300">
                        <li class="flex items-center gap-3"><div class="w-2 h-2 rounded-full bg-blue-500"></div> Distributed inference across 24 regions</li>
                        <li class="flex items-center gap-3"><div class="w-2 h-2 rounded-full bg-blue-500"></div> Native caching layer for static strings</li>
                        <li class="flex items-center gap-3"><div class="w-2 h-2 rounded-full bg-blue-500"></div> Automatic failover protection</li>
                    </ul>
                </div>
            </div>
            
            <div class="grid md:grid-cols-3 gap-6">
                <div class="feature-card p-8 rounded-2xl"><h3 class="text-xl font-bold mb-2">SOC2 Type II</h3><p class="text-gray-400 text-sm">Enterprise-grade security and compliance out of the box. Automatic data redaction protocols ensure PII never hits our underlying model training pipelines.</p></div>
                <div class="feature-card p-8 rounded-2xl"><h3 class="text-xl font-bold mb-2">99.999% SLA</h3><p class="text-gray-400 text-sm">Financially backed uptime guarantees for mission-critical paths. We provide a custom failover origin so your site never halts.</p></div>
                <div class="feature-card p-8 rounded-2xl"><h3 class="text-xl font-bold mb-2">Custom Glossaries</h3><p class="text-gray-400 text-sm">Force the model to adhere to your specific brand terminology. Easily upload your existing translation memory via our intuitive CI/CD sync tools.</p></div>
            </div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas/pricing" {
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-48 pb-32 px-6 max-w-7xl mx-auto text-center animate-fade-in">
			<h1 class="text-6xl md:text-7xl font-bold tracking-tighter mb-6">Simple, scalable pricing.</h1>
			<p class="text-xl text-gray-400 max-w-2xl mx-auto leading-relaxed mb-20">Pay only for the edge compute you use. No seat limits or arbitrary domains constraints.</p>
			
			<div class="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto text-left">
			    
			    <div class="feature-card p-10 rounded-3xl flex flex-col">
			        <h3 class="text-xl font-semibold mb-2 text-gray-300">Developer</h3>
			        <div class="text-5xl font-bold mb-6">$49<span class="text-lg text-gray-500 font-normal">/mo</span></div>
			        <ul class="space-y-4 text-gray-400 mb-10 flex-grow">
			            <li class="flex gap-3 text-sm">✔ 50,000 requests/month</li>
			            <li class="flex gap-3 text-sm">✔ 2 edge regions</li>
			            <li class="flex gap-3 text-sm">✔ Community support</li>
			        </ul>
			        <a href="/saas/register" class="w-full py-3 rounded-xl font-medium bg-white/10 hover:bg-white/20 text-center transition">Start Free Trial</a>
			    </div>
			    
			    <div class="gradient-border p-10 rounded-3xl flex flex-col relative transform md:-translate-y-4">
			        <div class="absolute top-0 right-10 -translate-y-1/2 bg-blue-600 text-xs font-bold px-3 py-1 rounded-full">MOST POPULAR</div>
			        <h3 class="text-xl font-semibold mb-2 text-blue-400">Growth</h3>
			        <div class="text-5xl font-bold mb-6">$199<span class="text-lg text-gray-500 font-normal">/mo</span></div>
			        <ul class="space-y-4 text-gray-300 mb-10 flex-grow">
			            <li class="flex gap-3 text-sm">✔ 500,000 requests/month</li>
			            <li class="flex gap-3 text-sm">✔ All global edge regions</li>
			            <li class="flex gap-3 text-sm">✔ Advanced glossaries</li>
			            <li class="flex gap-3 text-sm">✔ Priority email support</li>
			        </ul>
			        <a href="/saas/register" class="w-full py-3 rounded-xl font-medium btn-primary text-center">Get Started</a>
			    </div>
			    
			    <div class="feature-card p-10 rounded-3xl flex flex-col">
			        <h3 class="text-xl font-semibold mb-2 text-gray-300">Enterprise</h3>
			        <div class="text-5xl font-bold mb-6 pr-4">Custom</div>
			        <ul class="space-y-4 text-gray-400 mb-10 flex-grow">
			            <li class="flex gap-3 text-sm">✔ Unlimited requests</li>
			            <li class="flex gap-3 text-sm">✔ Dedicated inference nodes</li>
			            <li class="flex gap-3 text-sm">✔ SOC2 / HIPAA compliance</li>
			            <li class="flex gap-3 text-sm">✔ Dedicated success manager</li>
			        </ul>
			        <a href="/saas#demo" class="w-full py-3 rounded-xl font-medium bg-white/10 hover:bg-white/20 text-center transition">Contact Sales</a>
			    </div>
			
			</div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas/about" {
		content = `{{define "content"}}
		<div class="hero-glow"></div>` + navHTML + `
		<main class="pt-48 pb-32 px-6 max-w-4xl mx-auto animate-fade-in">
			<h1 class="text-5xl md:text-7xl font-bold tracking-tighter mb-12">We believe language <br>should never be a barrier.</h1>
			<div class="prose prose-invert prose-lg text-gray-400 leading-relaxed font-light">
			    <p class="mb-8 font-medium text-white text-xl">Lumina AI was founded in 2024 by researchers from DeepMind and OpenAI who recognized the fundamental inadequacy of traditional localization workflows.</p>
			    <p class="mb-8">For decades, global companies have relied on a fragmented ecosystem of translation agencies, headless CMS plugins, and fragile string-replacement scripts. This approach scales poorly, divorces translation from the design context, and introduces massive latency. The manual back-and-forth between linguists and developers blocks rapid deployment cycles, meaning international users often receive features weeks after the baseline launch.</p>
			    <p class="mb-8">We built a unified neural engine that maps directly to your application layer via standard web requests. By leveraging highly optimized, context-aware LLMs deployed at the edge infrastructure layer, Lumina translates your entire DOM structure in milliseconds while preserving the exact layout, tone, and brand safety of the original content. We do not extract strings; we comprehend nodes.</p>
			    <p class="mb-8">Our proprietary context pipeline seamlessly ingests your CSS styling, ensuring that morphological expansion (such as German words being 30% longer than their English counterparts) doesn't shatter your meticulously orchestrated CSS grid or truncate critical call-to-action buttons.</p>
			    <p>Our mission is to enable any business, regardless of size, to operate natively in every market on Earth from day one.</p>
			</div>
			
			<div class="mt-20 grid grid-cols-2 md:grid-cols-4 gap-8 border-t border-white/10 pt-16 text-center">
			    <div><div class="text-4xl font-bold text-white mb-2">24</div><div class="text-sm text-gray-500 uppercase tracking-widest">Global Nodes</div></div>
			    <div><div class="text-4xl font-bold text-white mb-2">5B+</div><div class="text-sm text-gray-500 uppercase tracking-widest">Words/Day</div></div>
			    <div><div class="text-4xl font-bold text-white mb-2">&lt;50ms</div><div class="text-sm text-gray-500 uppercase tracking-widest">Avg. Latency</div></div>
			    <div><div class="text-4xl font-bold text-white mb-2">99%</div><div class="text-sm text-gray-500 uppercase tracking-widest">Context Accuracy</div></div>
			</div>
		</main>
		` + jsTrap + `
		{{end}}`
	} else if r.URL.Path == "/saas" || r.URL.Path == "/saas/" { // "/"
		content = `{{define "content"}}
		<div class="hero-glow"></div>
		` + navHTML + `

		<main class="pt-40 pb-20 px-6 max-w-7xl mx-auto flex flex-col items-center">
			{{if .Success}}
			<div class="animate-fade-in bg-blue-900/30 border border-blue-500/30 text-blue-200 px-6 py-4 rounded-xl mb-12 flex items-center gap-3 backdrop-blur-sm">
				<svg class="w-6 h-6 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
				<span class="font-medium">Success! Your action was completed.</span>
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
				
				{{if not .IsLoggedIn}}
				<form method="POST" action="/saas/api/contact" id="demo" class="flex flex-col sm:flex-row gap-3 justify-center max-w-md mx-auto delay-100 animate-fade-in">
					` + hiddenInputs + `
					<input type="email" name="email" placeholder="Enter work email" class="flex-1 bg-white/5 border border-white/10 text-white px-6 py-4 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all placeholder-gray-500" required>
					<button type="submit" class="btn-primary text-white px-8 py-4 rounded-xl font-medium tracking-wide">Get Demo</button>
				</form>
				{{end}}
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
			
			<div id="faq" class="mt-40 max-w-4xl mx-auto w-full animate-fade-in delay-300 mb-20 text-left">
				<h2 class="text-4xl font-bold mb-12 text-center">Frequently Asked Questions</h2>
				<div class="space-y-6">
					<div class="bg-white/5 border border-white/10 rounded-2xl p-8">
						<h3 class="text-xl font-bold mb-3">How does Edge Inference affect SEO?</h3>
						<p class="text-gray-400 leading-relaxed font-light">Because Lumina manipulates the DOM on the edge worker before the HTML stream completes, search engine crawlers like Googlebot receive perfectly localized static HTML. There is no client-side rendering delay, ensuring maximum SEO scoring and Core Web Vitals across all localized regions.</p>
					</div>
					<div class="bg-white/5 border border-white/10 rounded-2xl p-8">
						<h3 class="text-xl font-bold mb-3">Does Lumina support modern JS Frameworks?</h3>
						<p class="text-gray-400 leading-relaxed font-light">Yes. Lumina operates at the reverse proxy layer, not the application framework layer. Whether you ship pure HTML, Hydrated React via Next.js, or complex WebGL canvases, our inference engine interprets the outbound HTTP response payload and translates strings agnostically. It is completely decoupled from your underlying tech stack.</p>
					</div>
					<div class="bg-white/5 border border-white/10 rounded-2xl p-8">
						<h3 class="text-xl font-bold mb-3">How secure is my proprietary data?</h3>
						<p class="text-gray-400 leading-relaxed font-light">Lumina is SOC2 Type II standard. We utilize distinct, isolated model pipelines per enterprise tenant. We never use your traffic to train generalized foundation models. All inference occurs in volatile memory at the edge, ensuring zero persistent storage of your user's Personally Identifiable Information (PII) requests.</p>
					</div>
					<div class="bg-white/5 border border-white/10 rounded-2xl p-8">
						<h3 class="text-xl font-bold mb-3">Can I inject my existing Translation Memory?</h3>
						<p class="text-gray-400 leading-relaxed font-light">Absolutely. Lumina supports standard TMX, XLIFF, and CSV glossary uploads via API. Our proprietary RAG (Retrieval-Augmented Generation) context injection forces the inference model to perfectly replicate your historically approved brand terminology before falling back on its foundational linguistic knowledge.</p>
					</div>
				</div>
			</div>
			
			<script>
				document.addEventListener('mousemove', (e) => {
					const glow = document.querySelector('.hero-glow');
					if(!glow) return;
					const x = (e.clientX / window.innerWidth - 0.5) * 40;
					const y = (e.clientY / window.innerHeight - 0.5) * 40;
					glow.style.transform = "translate(calc(-50% + " + x + "px), " + y + "px)";
				});
			</script>
		</main>
		` + jsTrap + `
		{{end}}`
	} else {
		http.NotFound(w, r)
		return
	}

	Render(w, head, content, PageData{
		Brand:       "Lumina",
		Title:       "Enterprise AI Translation",
		Description: "Automate your global workflow with Lumina AI.",
		Success:     success,
		Query:       r.URL.Query().Get("error"),
		UserEmail:   userEmail,
		IsLoggedIn:  isLoggedIn,
	})
}
