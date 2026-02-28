package sites

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mock_server/internal/db"

	"github.com/google/uuid"
)

func ShopHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "artisan-session")

	cartSessionID, ok := session.Values["cart_id"].(string)
	if !ok {
		cartSessionID = uuid.New().String()
		session.Values["cart_id"] = cartSessionID
		session.Save(r, w)
	}

	if strings.HasPrefix(r.URL.Path, "/shop/cart/add") {
		productID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if productID > 0 {
			db.DB.Exec("INSERT INTO cart_items (session_id, product_id, quantity) VALUES (?, ?, 1) ON CONFLICT(session_id, product_id) DO UPDATE SET quantity=quantity+1", cartSessionID, productID)
		}
		http.Redirect(w, r, "/shop/cart", http.StatusSeeOther)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/shop/cart/remove") {
		productID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if productID > 0 {
			db.DB.Exec("DELETE FROM cart_items WHERE session_id = ? AND product_id = ?", cartSessionID, productID)
		}
		http.Redirect(w, r, "/shop/cart", http.StatusSeeOther)
		return
	}

	var cartCount int
	db.DB.QueryRow("SELECT COALESCE(SUM(quantity), 0) FROM cart_items WHERE session_id = ?", cartSessionID).Scan(&cartCount)

	head := `{{define "head"}}
    <style>
        body { font-family: 'Playfair Display', serif; background-color: #faf9f6; color: #1c1917; }
        .nav-link { position: relative; font-family: 'Inter', sans-serif; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.1em; }
        .nav-link::after { content: ''; position: absolute; width: 0; height: 1px; bottom: -4px; right: 0; background-color: #1c1917; transition: width 0.4s cubic-bezier(0.25, 0.8, 0.25, 1); }
        .nav-link:hover::after { width: 100%; left: 0; right: auto; }
        html { scroll-behavior: smooth; }
        .parallax-bg {
            background-attachment: fixed;
            background-position: center;
            background-repeat: no-repeat;
            background-size: cover;
        }
        .product-card { transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1); }
        .product-card img { transition: transform 0.8s cubic-bezier(0.165, 0.84, 0.44, 1); }
        .product-card:hover img { transform: scale(1.05); }
        .fade-up { opacity: 0; transform: translateY(40px); animation: fadeUpAnim 1.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
        @keyframes fadeUpAnim { to { opacity: 1; transform: translateY(0); } }
        .delay-100 { animation-delay: 100ms; }
        .delay-200 { animation-delay: 200ms; }
        .delay-300 { animation-delay: 300ms; }
        .delay-400 { animation-delay: 400ms; }
    </style>
    {{end}}`

	// LIMIT TEST: Large Base64 inline tracking icon inside the nav
	var largeB64Builder strings.Builder
	for i := 0; i < 5000; i++ {
		largeB64Builder.WriteString("M10 10 H 90 V 90 H 10 Z ")
	}
	hugeSVG := fmt.Sprintf(`<div style="display:none;"><svg viewBox="0 0 100 100"><path d="%s" fill="none"/></svg></div>`, largeB64Builder.String())

	navHTML := `
    <nav class="fixed top-0 w-full z-50 bg-[#faf9f6]/90 backdrop-blur-md border-b border-stone-200 py-6 transition-all">
        ` + hugeSVG + `
        <div class="max-w-7xl mx-auto px-8 flex justify-between items-center">
            <div class="flex gap-8">
                <a href="/shop" class="nav-link">Shop</a>
                <a href="/shop/journal" class="nav-link">Journal</a>
                <a href="/shop/about" class="nav-link">Our Story</a>
            </div>
            <a href="/shop" class="text-3xl font-bold tracking-tighter">Artisan Brews</a>
            <div class="flex gap-6 items-center">
                <a href="/shop/cart" class="nav-link flex items-center gap-2">
                    Cart 
                    <span class="w-5 h-5 rounded-full bg-stone-900 text-white text-xs flex items-center justify-center font-sans">{{.CartCount}}</span>
                </a>
            </div>
        </div>
    </nav>`

	var content string

	if strings.HasPrefix(r.URL.Path, "/shop/product/") {
		productID := strings.TrimPrefix(r.URL.Path, "/shop/product/")
		var p struct {
			ID          int
			Name        string
			Description string
			Price       float64
			ImageURL    string
		}

		err := db.DB.QueryRow("SELECT id, name, description, price, image_url FROM products WHERE id = ?", productID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL)
		if err != nil {
			log.Printf("Product lookup failed: %v", err)
			http.NotFound(w, r)
			return
		}

		content = `{{define "content"}}` + navHTML + `
		<main class="pt-32 pb-20 px-8 max-w-7xl mx-auto grid md:grid-cols-2 gap-16 items-center min-h-[80vh] fade-up">
			<div class="overflow-hidden bg-stone-100">
				<img src="` + p.ImageURL + `" alt="` + p.Name + `" class="w-full h-[600px] object-cover hover:scale-105 transition-transform duration-1000">
			</div>
			<div class="flex flex-col justify-center">
				<div class="font-sans text-xs tracking-widest text-stone-500 uppercase mb-4">Single Origin &bull; Whole Bean</div>
				<h1 class="text-5xl md:text-6xl font-bold mb-6 tracking-tight leading-tight">` + p.Name + `</h1>
				<p class="text-xl text-stone-600 mb-8 font-light leading-relaxed">` + p.Description + `</p>
				<div class="flex items-end gap-6 mb-12">
					<span class="text-3xl font-serif">$` + strconv.FormatFloat(p.Price, 'f', 2, 64) + ` <span class="text-base text-stone-400">USD</span></span>
				</div>
				<a href="/shop/cart/add?id=` + strconv.Itoa(p.ID) + `" class="w-full py-4 bg-stone-900 hover:bg-stone-800 text-white text-center font-sans tracking-widest uppercase text-sm transition-colors duration-300">
					Add to Cart
				</a>
                
                <div class="mt-16 pt-8 border-t border-stone-200">
                    <h3 class="font-bold text-lg mb-4">Brewing Details</h3>
                    <ul class="space-y-3 font-sans text-sm text-stone-500">
                        <li class="flex justify-between border-b border-stone-100 pb-2"><span>Process</span> <span>Washed & Sun-Dried</span></li>
                        <li class="flex justify-between border-b border-stone-100 pb-2"><span>Elevation</span> <span>1,800m - 2,100m</span></li>
                        <li class="flex justify-between border-b border-stone-100 pb-2"><span>Harvest</span> <span>Current Season</span></li>
                    </ul>
                </div>
			</div>
		</main>
		{{end}}`

	} else if r.URL.Path == "/shop/journal" {
		// Massive original content to ensure it scrolls past 3 viewports
		content = `{{define "content"}}` + navHTML + `
		<main class="pt-40 pb-32">
            <header class="text-center max-w-3xl mx-auto px-6 mb-24 fade-up">
                <h1 class="text-6xl font-bold tracking-tight mb-6">The Journal</h1>
                <p class="font-sans text-stone-500 uppercase tracking-widest text-sm">Notes on origin, roasting, and the pursuit of the perfect cup.</p>
            </header>

            <article class="max-w-5xl mx-auto px-6 grid md:grid-cols-2 gap-12 items-center mb-32 fade-up delay-100">
                <img src="https://images.unsplash.com/photo-1511920170033-f8396924c348?q=80&w=1000" alt="Roasting Process" class="w-full aspect-[4/5] object-cover">
                <div>
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-4">Origins &bull; Oct 12</div>
                    <h2 class="text-4xl font-bold mb-6 leading-tigher">The Science of the Maillard Reaction in Coffee Roasting</h2>
                    <p class="text-stone-600 mb-8 leading-relaxed">It's the precise moment when the beans transform. The sugars caramelize, the amino acids react, and the hundreds of aromatic compounds that define your morning cup are born. We spent three weeks adjusting our exhaust temperatures to perfect this reaction.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">The drum roaster hits 380 degrees Fahrenheit, marking the onset of first crack. Moisture inside the cellular structure of the bean violently expands, creating an audial signature that our master roasters listen for with almost obsessive dedication. This acoustic cue informs our dropping temperature, locking in the brightness without allowing bitter distillates to develop.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">Every batch behaves differently. Humidity, ambient temperature, and barometric pressure inside the roastery mean that profiling relies on sensory feedback as much as our digital logging softwares. The art lies within the tension of precise science and intuitive craft.</p>
                </div>
            </article>

            <article class="max-w-5xl mx-auto px-6 grid md:grid-cols-2 gap-12 items-center mb-32 fade-up delay-200">
                <div class="order-2 md:order-1">
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-4">Travel &bull; Sep 28</div>
                    <h2 class="text-4xl font-bold mb-6 leading-tigher">Sourcing at 2,000 Meters: Dispatches from Colombia</h2>
                    <p class="text-stone-600 mb-8 leading-relaxed">The air is thin, but the soil is rich. Our latest buying trip took us far past the paved roads of Antioquia to meet the farmers growing the gesha varietals that will make up our winter exclusive release. It is a labor of intense generational love.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">Navigating the dirt switchbacks of the Andes required four-wheel drive and a tremendous amount of patience. When we finally arrived at Finca El Paraiso, the clouds were literally rolling through the drying beds. The farmers here employ a specialized double anaerobic fermentation process, sealing the cherry in massive barrels for 72 hours before depulping.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">This meticulous process requires unparalleled sanitation. The slightest introduction of the wrong microbial flora can taint an entire seasonal harvest, rendering months of agricultural labor entirely void. We spent two weeks cupside, logging flavor notes, testing water activity, and establishing a relationship based on mutual reverence for the bean.</p>
                </div>
                <img src="https://images.unsplash.com/photo-1497935586351-b67a49e012bf?q=80&w=1000" alt="Coffee Farm" class="w-full aspect-[4/5] object-cover order-1 md:order-2">
            </article>

            <article class="max-w-5xl mx-auto px-6 grid md:grid-cols-2 gap-12 items-center mb-32 fade-up delay-300">
                <img src="https://images.unsplash.com/photo-1509042239860-f550ce710b93?q=80&w=1000" alt="Pour Over Coffee" class="w-full aspect-[4/5] object-cover">
                <div>
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-4">Technique &bull; Aug 14</div>
                    <h2 class="text-4xl font-bold mb-6 leading-tigher">Mastering the Bloom: A Treatise on Water Chemistry</h2>
                    <p class="text-stone-600 mb-8 leading-relaxed">Your extraction is only as good as your solvent. We dive deep into total dissolved solids, alkalinity, and magnesium ratios to understand why your tap water might be neutralizing the sparkling acidity of your single origin pour-overs.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">Consider the standard drip coffee maker. It sprays water unevenly, often barely clinging to 190 degrees Fahrenheit, channeling directly through the weakest point in the bed. True extraction requires thermal stability (precisely 205 degrees for light roasts) and a deliberate pouring technique that agitates the grounds evenly, allowing the trapped carbon dioxide to escape during the initial bloom phase.</p>
                    <p class="text-stone-600 mb-8 leading-relaxed">But even perfect technique fails if the buffer capacity of the water is too high. Bicarbonates will literally erase the delicate malic and citric acids inherent in lightly roasted African coffees. Our customized particulate filters ensure a brewing water profile optimized strictly for transparency and prolonged sweetness.</p>
                </div>
            </article>
            
            <section class="max-w-7xl mx-auto px-6 mt-32 grid md:grid-cols-3 gap-8 border-t border-stone-200 pt-16 fade-up delay-400">
                <div>
                    <img src="https://images.unsplash.com/photo-1541167760496-1628856ab772?q=80&w=600" class="w-full aspect-square object-cover mb-4">
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-2">Brew Guide</div>
                    <h3 class="text-xl font-bold">Perfecting the Chemex</h3>
                </div>
                <div>
                    <img src="https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?q=80&w=600" class="w-full aspect-square object-cover mb-4">
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-2">Culture</div>
                    <h3 class="text-xl font-bold">The Rise of Light Roasts</h3>
                </div>
                <div>
                    <img src="https://images.unsplash.com/photo-1620916566398-39f1143ab7be?q=80&w=600" class="w-full aspect-square object-cover mb-4">
                    <div class="font-sans text-stone-400 text-xs tracking-widest uppercase mb-2">Equipment</div>
                    <h3 class="text-xl font-bold">Burr vs Blade Grinders</h3>
                </div>
            </section>
		</main>
		{{end}}`

	} else if r.URL.Path == "/shop/about" {
		content = `{{define "content"}}` + navHTML + `
		<main>
		    <section class="h-[80vh] w-full parallax-bg flex items-center justify-center relative fade-up" style="background-image: url('https://images.unsplash.com/photo-1447933601403-0c6688de566e?q=80&w=1600');">
		        <div class="absolute inset-0 bg-stone-900/40"></div>
                <div class="relative z-10 text-center text-white px-6">
                    <h1 class="text-6xl md:text-8xl font-bold mb-6 tracking-tight">Our Story</h1>
                    <p class="font-sans tracking-widest text-sm uppercase">EST. 2018 &bull; PORTLAND, OREGON</p>
                </div>
            </section>

            <section class="max-w-4xl mx-auto px-6 py-32 text-center fade-up delay-100">
                <h2 class="text-4xl md:text-5xl font-bold mb-10 leading-tight">We believe that coffee is a <br/><span class="italic text-stone-500">culinary experience</span>, not just a caffeine delivery mechanism.</h2>
                <div class="w-px h-24 bg-stone-300 mx-auto mb-10"></div>
                <p class="text-xl text-stone-600 leading-relaxed font-light mb-8">
                    Artisan Brews began in a small garage with a 1kg San Franciscan roaster and a single origin Ethiopian bean. We were obsessed with tracing the lineage of the bean, understanding the soil terroir, and executing the roast profile to a microscopic degree of precision.
                </p>
                <p class="text-xl text-stone-600 leading-relaxed font-light mb-8">
                    Our initial foray into commercial roasting was fraught with failures. Sourcing green coffee inherently exposes you to the massive volatility of agricultural supply chains. We navigated container delays, fluctuating moisture contents, and the absolute destruction of a 60kg bag of rare Yemenite beans due to an over-zealous charge temperature.
                </p>
                <p class="text-xl text-stone-600 leading-relaxed font-light">
                    Today, we partner directly with farmers across four continents. We pay significantly above fair-trade minimums to ensure sustainable, ethical harvesting. Every bag we ship is roasted to order, ensuring it arrives at the absolute peak of its degassing phase. The journey from crop to cup is intensely collaborative.
                </p>
            </section>

            <section class="grid md:grid-cols-2 fade-up delay-200">
                <div class="bg-stone-900 text-stone-300 p-16 md:p-32 flex flex-col justify-center">
                    <h3 class="text-3xl font-bold text-white mb-6">The Roastery</h3>
                    <p class="font-sans font-light leading-relaxed mb-8">Located in the historic industrial waterfront, our flagship roastery is completely open to the public. We invite you to watch the Loring Smart Roaster in action, smell the chaff, and taste the immediate results at our cupping bar. We host weekly cupping sessions for the community to develop their palates.</p>
					<p class="font-sans font-light leading-relaxed mb-8">Within the reinforced brick walls, massive palettes of raw green beans age inside specialized hermetic silos. We monitor ambient temperature and humidity strictly, treating our green coffee inventory like a high-end wine cellar. This obsessive environmental control eliminates variables before the beans even enter the hopper.</p>
                    <a href="/shop" class="nav-link text-white self-start">Shop Now</a>
                </div>
                <!-- Fixed broken image below -->
                <img src="https://images.unsplash.com/photo-1509042239860-f550ce710b93?q=80&w=1000" alt="Our Roastery" class="w-full h-full object-cover">
            </section>
		</main>
		{{end}}`

	} else if r.URL.Path == "/shop/cart" {
		rows, err := db.DB.Query(`
			SELECT p.id, p.name, p.price, p.image_url, c.quantity
			FROM cart_items c
			JOIN products p ON c.product_id = p.id
			WHERE c.session_id = ?
		`, cartSessionID)

		var cartItemsHTML string
		var total float64
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, qty int
				var name, img string
				var price float64
				rows.Scan(&id, &name, &price, &img, &qty)
				itemTotal := price * float64(qty)
				total += itemTotal

				cartItemsHTML += `
				<div class="flex items-center justify-between border-b border-stone-200 py-6 fade-up">
					<div class="flex items-center gap-6">
						<img src="` + img + `" alt="` + name + `" class="w-24 h-24 object-cover bg-stone-100">
						<div>
							<h3 class="font-bold text-lg text-stone-900">` + name + `</h3>
							<p class="text-stone-500 font-sans text-sm mt-1">Quantity: ` + strconv.Itoa(qty) + `</p>
							<a href="/shop/cart/remove?id=` + strconv.Itoa(id) + `" class="text-red-800 hover:text-red-600 text-xs uppercase tracking-widest font-sans font-bold mt-2 inline-block transition-colors">Remove</a>
						</div>
					</div>
					<div class="text-xl font-serif font-medium">
						$` + strconv.FormatFloat(itemTotal, 'f', 2, 64) + `
					</div>
				</div>`
			}
		}

		if cartItemsHTML == "" {
			cartItemsHTML = `<p class="text-stone-500 italic font-serif py-12 text-center text-xl">Your cart is currently empty.</p>`
		}

		content = `{{define "content"}}` + navHTML + `
		<main class="pt-40 pb-32 px-6 max-w-4xl mx-auto min-h-[70vh]">
			<header class="mb-16 border-b border-stone-900 pb-8 fade-up">
				<h1 class="text-5xl font-bold tracking-tight">Your Cart</h1>
			</header>

			<div class="mb-12">
				` + cartItemsHTML + `
			</div>

			{{if gt .CartCount 0}}
			<div class="bg-stone-100 p-8 flex flex-col items-end fade-up delay-100">
				<div class="flex justify-between w-full md:w-1/2 mb-4 font-sans text-stone-500">
					<span>Subtotal</span>
					<span>$` + strconv.FormatFloat(total, 'f', 2, 64) + `</span>
				</div>
				<div class="flex justify-between w-full md:w-1/2 mb-8 font-sans text-stone-500">
					<span>Shipping</span>
					<span>Calculated at checkout</span>
				</div>
				<div class="flex justify-between w-full md:w-1/2 mb-8 text-2xl font-bold border-t border-stone-200 pt-4">
					<span>Total</span>
					<span>$` + strconv.FormatFloat(total, 'f', 2, 64) + `</span>
				</div>
				
				<a href="/shop" class="w-full md:w-1/2 py-4 bg-stone-900 hover:bg-stone-800 text-white text-center font-sans tracking-widest uppercase text-sm transition-colors duration-300">
					Proceed to Checkout
				</a>
			</div>
			{{end}}
		</main>
		{{end}}`

	} else if r.URL.Path == "/shop" || r.URL.Path == "/shop/" {

		rows, err := db.DB.Query("SELECT id, name, description, price, image_url FROM products LIMIT 16")
		var productsHTML string

		if err == nil {
			var i int
			for rows.Next() {
				var id int
				var name, desc, img string
				var price float64
				rows.Scan(&id, &name, &desc, &price, &img)

				marginTop := "mt-0"
				if i%2 != 0 {
					marginTop = "md:mt-32"
				}

				// LIMIT TEST: Div Soup Injection (Deep DOM Nesting ~20 levels)
				var soupOpen, soupClose strings.Builder
				for j := 0; j < 20; j++ {
					soupOpen.WriteString(fmt.Sprintf(`<div class="depth-layer layout-group-%d" data-render-cycle="%d">`, j, j))
					soupClose.WriteString(`</div>`)
				}

				// LIMIT TEST: High Attribute Density on Images
				longText := "Premium artisan roasted coffee beans sourced directly from sustainable farms. Our unique roasting profile brings out the natural sweetness and complex flavor notes of each distinct origin. Experience unparalleled transparency in every cup."
				denseAttrs := fmt.Sprintf(`alt="%s" title="%s" data-description="%s" aria-label="%s" data-category="single-origin" data-metadata="{'id':%d, 'stock':'available'}" data-blob-1="%s" data-blob-2="%s" data-blob-3="%s"`, name, name, desc, name, id, longText, longText, longText)

				productsHTML += `
				<a href="/shop/product/` + strconv.Itoa(id) + `" class="product-card group block ` + marginTop + ` mb-16 fade-up delay-` + strconv.Itoa((i%3)*100) + `">
					` + soupOpen.String() + `
					<div class="overflow-hidden bg-stone-100 mb-6 aspect-[4/5]">
						<img src="` + img + `" ` + denseAttrs + ` class="w-full h-full object-cover">
					</div>
					<div class="flex justify-between items-start font-sans">
						<div>
							<h3 class="text-lg font-bold text-stone-900 font-serif">` + name + `</h3>
							<p class="text-sm text-stone-500 mt-1">` + desc + `</p>
						</div>
						<span class="text-sm font-medium">$` + strconv.FormatFloat(price, 'f', 2, 64) + `</span>
					</div>
					` + soupClose.String() + `
				</a>`
				i++
			}
		}

		content = `{{define "content"}}` + navHTML + `
		<main class="pt-32 pb-20 px-6 max-w-7xl mx-auto">
			<header class="text-center mb-32 pt-20 fade-up">
				<h1 class="text-7xl md:text-9xl font-bold tracking-tighter mb-6 mx-auto">Roasted to <br><span class="italic font-light text-stone-500">Order.</span></h1>
				<p class="font-sans text-stone-500 uppercase tracking-widest text-sm mt-8">Single Origin &bull; Ethically Sourced</p>
			</header>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-x-12 lg:gap-x-24">
				` + productsHTML + `
			</div>
			
			<div class="text-center mt-20 fade-up delay-300">
			    <a href="#" class="inline-block border border-stone-900 text-stone-900 px-12 py-4 font-sans text-xs tracking-widest uppercase hover:bg-stone-900 hover:text-white transition-colors duration-300">View All Coffees</a>
			</div>
		</main>
		{{end}}`
	} else {
		http.NotFound(w, r)
		return
	}

	Render(w, head, content, PageData{
		Brand:       "Artisan",
		Title:       "Premium Roasters",
		Description: "Single origin specialty coffee.",
		CartCount:   cartCount,
	})
}
