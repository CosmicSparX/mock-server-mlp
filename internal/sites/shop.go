package sites

import (
	"fmt"
	"net/http"
	"strconv"
)

func ShopHandler(w http.ResponseWriter, r *http.Request) {
	cartStr := r.URL.Query().Get("cart")
	cartCount, _ := strconv.Atoi(cartStr)

	action := r.URL.Query().Get("action")
	if action == "add" {
		http.Redirect(w, r, fmt.Sprintf("/shop?cart=%d", cartCount+1), http.StatusSeeOther)
		return
	}

	head := `{{define "head"}}
    <style>
        body { font-family: 'Inter', sans-serif; background-color: #f8f5f0; color: #1c1917; }
        .font-serif-heading { font-family: 'Playfair Display', serif; }
        
        .nav-link { position: relative; display: inline-block; padding-bottom: 2px; }
        .nav-link::after {
            content: ''; position: absolute; width: 0; height: 1px; bottom: 0; left: 0;
            background-color: #1c1917; transition: width 0.3s ease;
        }
        .nav-link:hover::after { width: 100%; }
        
        .product-card {
            transition: all 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
        }
        .image-container { overflow: hidden; }
        .image-container img {
            transition: transform 0.8s cubic-bezier(0.165, 0.84, 0.44, 1);
        }
        .product-card:hover .image-container img {
            transform: scale(1.05);
        }
        .add-btn {
            opacity: 0; transform: translateY(10px);
            transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
        }
        .product-card:hover .add-btn {
            opacity: 1; transform: translateY(0);
        }
        
        /* Reveal animations */
        @keyframes revealUp {
            from { opacity: 0; transform: translateY(40px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .reveal { animation: revealUp 1.2s cubic-bezier(0.19, 1, 0.22, 1) forwards; opacity: 0; }
        .delay-1 { animation-delay: 0.1s; }
        .delay-2 { animation-delay: 0.3s; }
        .delay-3 { animation-delay: 0.5s; }
        .delay-4 { animation-delay: 0.7s; }
    </style>
    {{end}}`

	content := `{{define "content"}}
    <header class="w-full px-8 py-6 flex justify-between items-center fixed top-0 bg-[#f8f5f0]/90 backdrop-blur-md z-50 border-b border-stone-200/50 reveal">
        <div class="flex gap-8 text-sm uppercase tracking-widest font-medium">
            <a href="#" class="nav-link">Shop</a>
            <a href="#" class="nav-link">Journal</a>
            <a href="#" class="nav-link">About</a>
        </div>
        <h1 class="text-2xl font-serif-heading font-semibold tracking-wide"><a href="/shop">ARTISAN</a></h1>
        <div class="flex gap-6 text-sm uppercase tracking-widest font-medium items-center">
            <a href="#" class="nav-link">Account</a>
            <button class="flex items-center gap-2 group">
                <span class="nav-link">Cart</span>
                <span class="w-6 h-6 rounded-full bg-stone-900 text-[#f8f5f0] flex items-center justify-center text-xs group-hover:scale-110 transition-transform">{{.CartCount}}</span>
            </button>
        </div>
    </header>

    <main class="w-full max-w-7xl mx-auto pt-40 px-8 pb-32">
        <div class="mb-24 text-center reveal delay-1 max-w-3xl mx-auto">
            <p class="text-sm uppercase tracking-widest text-stone-500 mb-6">Our Collection</p>
            <h2 class="text-5xl md:text-7xl font-serif-heading leading-tight text-stone-900">
                Ethically sourced.<br>
                <span class="italic text-stone-500">Masterfully roasted.</span>
            </h2>
        </div>

        <!-- Asymmetric Grid -->
        <div class="grid grid-cols-1 md:grid-cols-12 gap-y-20 md:gap-x-12 items-center">
            
            <!-- Large Focus Item -->
            <div class="md:col-span-7 product-card group cursor-pointer reveal delay-2">
                <div class="image-container rounded-sm aspect-[4/5] bg-stone-200 mb-6 relative">
                    <img src="https://images.unsplash.com/photo-1559525839-b184a4d698c7?q=80&w=1000&auto=format&fit=crop" class="w-full h-full object-cover">
                    <div class="absolute inset-0 bg-black/10 transition-opacity group-hover:opacity-0"></div>
                </div>
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="text-2xl font-serif-heading mb-1">Ethiopian Yirgacheffe</h3>
                        <p class="text-stone-500 font-light track-wide">Floral, Jasmine, Peach</p>
                    </div>
                    <div class="text-right flex flex-col items-end">
                        <span class="text-lg mb-2">$24.00</span>
                        <a href="/shop?action=add&cart={{.CartCount}}" class="add-btn inline-block border border-stone-900 text-stone-900 hover:bg-stone-900 hover:text-[#f8f5f0] px-4 py-2 text-sm uppercase tracking-widest transition-colors">Add to Cart</a>
                    </div>
                </div>
            </div>

            <!-- Offset Text / Quote block -->
            <div class="md:col-span-5 md:pl-10 reveal delay-3 flex justify-center">
                <div class="max-w-xs text-center md:text-left">
                     <p class="font-serif-heading italic text-3xl leading-snug text-stone-700 mb-6 border-l-4 border-stone-300 pl-6">
                        "The pursuit of the perfect cup begins with respecting the soil and the farmer."
                     </p>
                     <p class="text-sm uppercase tracking-widest text-stone-500">— Chief Roaster</p>
                </div>
            </div>

            <!-- Smaller Item 1 -->
            <div class="md:col-span-4 md:col-start-2 product-card group cursor-pointer reveal delay-3">
                <div class="image-container rounded-sm aspect-square bg-stone-200 mb-6 relative">
                    <img src="https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?q=80&w=600&auto=format&fit=crop" class="w-full h-full object-cover">
                     <div class="absolute inset-0 bg-black/10 transition-opacity group-hover:opacity-0"></div>
                </div>
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="text-xl font-serif-heading mb-1">Colombian Supremo</h3>
                        <p class="text-stone-500 font-light text-sm">Chocolate, Caramel, Nuts</p>
                    </div>
                    <div class="text-right flex flex-col items-end">
                        <span class="text-md mb-2">$21.00</span>
                        <a href="/shop?action=add&cart={{.CartCount}}" class="add-btn inline-block border border-stone-900 text-stone-900 hover:bg-stone-900 hover:text-[#f8f5f0] px-3 py-1.5 text-xs uppercase tracking-widest transition-colors">Add</a>
                    </div>
                </div>
            </div>

            <!-- Smaller Item 2 -->
            <div class="md:col-span-5 md:col-start-7 product-card group cursor-pointer reveal delay-4 mt-20 md:mt-40">
                <div class="image-container rounded-sm aspect-[3/4] bg-stone-200 mb-6 relative">
                    <img src="https://images.unsplash.com/photo-1587734195503-904fca47e0e9?q=80&w=600&auto=format&fit=crop" class="w-full h-full object-cover">
                     <div class="absolute inset-0 bg-black/10 transition-opacity group-hover:opacity-0"></div>
                </div>
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="text-xl font-serif-heading mb-1">House Espresso</h3>
                        <p class="text-stone-500 font-light text-sm">Dark Cocoa, Molasses</p>
                    </div>
                    <div class="text-right flex flex-col items-end">
                        <span class="text-md mb-2">$22.00</span>
                        <a href="/shop?action=add&cart={{.CartCount}}" class="add-btn inline-block border border-stone-900 text-stone-900 hover:bg-stone-900 hover:text-[#f8f5f0] px-3 py-1.5 text-xs uppercase tracking-widest transition-colors">Add</a>
                    </div>
                </div>
            </div>

        </div>
    </main>
    {{end}}`

	Render(w, head, content, PageData{
		Brand:       "Artisan Brews",
		Title:       "Shop Premium Coffee",
		Description: "Ethically sourced, masterfully roasted coffee beans.",
		CartCount:   cartCount,
	})
}
