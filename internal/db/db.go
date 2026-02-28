package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./mock_server.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()
}

func createTables() {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			image_url TEXT,
			stock INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS cart_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			UNIQUE(session_id, product_id)
		);`,
		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			excerpt TEXT,
			content TEXT NOT NULL,
			category TEXT,
			author TEXT,
			image_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS newsletter_subs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL
		);`,
	}

	for _, schema := range schemas {
		_, err := DB.Exec(schema)
		if err != nil {
			log.Fatalf("Failed to create table: %v. Schema: %s", err, schema)
		}
	}

	seedData()
}

func seedData() {
	// Seed some products if empty
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		DB.Exec(`INSERT INTO products (name, description, price, stock, image_url) VALUES 
			('Ethiopian Yirgacheffe', 'Floral, Jasmine, Peach - Light Roast', 24.00, 50, 'https://images.unsplash.com/photo-1675306408031-a9aad9f23308?q=80&w=1000'),
			('Colombian Supremo', 'Chocolate, Caramel, Nuts - Medium Roast', 21.00, 30, 'https://images.unsplash.com/photo-1606486544554-164d98da4889?q=80&w=600'),
			('House Espresso', 'Dark Cocoa, Molasses - Dark Roast', 22.00, 100, 'https://images.unsplash.com/photo-1580933073521-dc49ac0d4e6a?q=80&w=600'),
			('Guatemalan Antigua', 'Spicy, Smoky, Rich - Medium Roast', 23.50, 45, 'https://images.unsplash.com/photo-1513530176992-0cf39c4cbed4?q=80&w=600'),
			('Kenyan AA', 'Grapefruit, Blackcurrant, Bright - Light Roast', 26.00, 20, 'https://images.unsplash.com/photo-1607681034540-2c46cc71896d?q=80&w=600'),
			('Sumatra Mandheling', 'Earthy, Herbal, Heavy Body - Dark Roast', 22.50, 60, 'https://images.unsplash.com/photo-1690983326555-8b8e27843a32?q=80&w=600'),
			('Costa Rican Tarrazu', 'Honey, Citrus, Crisp - Medium/Light Roast', 24.50, 40, 'https://images.unsplash.com/photo-1587734195503-904fca47e0e9?q=80&w=600'),
			('Decaf Swiss Water', 'Mild, Sweet, Balanced - Medium Roast', 20.00, 80, 'https://images.unsplash.com/photo-1624976921221-4478954a4bc9?q=80&w=600'),
			('Panama Gesha', 'Bergamot, Melon, Silky - Light Roast', 45.00, 15, 'https://images.unsplash.com/photo-1587049016823-69ef9d68bd44?q=80&w=600'),
			('Honduras Marcala', 'Vanilla, Orange, Smooth - Medium Roast', 21.50, 35, 'https://images.unsplash.com/photo-1692296113053-76f240e5ce33?q=80&w=600'),
			('Burundi Ngozi', 'Cranberry, Black Tea, Syrupy - Medium/Light', 25.00, 25, 'https://images.unsplash.com/photo-1511920170033-f8396924c348?q=80&w=600'),
			('Mexican Chiapas', 'Milk Chocolate, Almond, Mild - Medium Roast', 19.50, 55, 'https://images.unsplash.com/photo-1692299116305-762729f2e9d5?q=80&w=600'),
			('Tanzania Peaberry', 'Lemon, Blackberry, Winey - Light Roast', 27.00, 20, 'https://images.unsplash.com/photo-1606791405792-1004f1718d0c?q=80&w=600'),
			('Rwanda Kivu', 'Cherry, Brown Sugar, Juicy - Medium Roast', 23.00, 40, 'https://images.unsplash.com/photo-1610632380989-680fe40816c6?q=80&w=600'),
			('Indian Monsooned Malabar', 'Spicy, Musty, Heavy - Dark Roast', 22.00, 50, 'https://images.unsplash.com/photo-1692299108834-038511803008?q=80&w=600'),
			('Papua New Guinea Sigri', 'Mango, Papaya, Clean - Medium Roast', 24.00, 30, 'https://images.unsplash.com/photo-1770941550709-a555ac4a69b3?q=80&w=600')
		`)
	}

	// Seed some articles if empty
	var countArts int
	DB.QueryRow("SELECT COUNT(*) FROM articles").Scan(&countArts)
	if countArts == 0 {
		DB.Exec(`INSERT INTO articles (title, excerpt, content, category, author, image_url) VALUES 
			('The AI Infrastructure Squeeze', 'Why raw compute is becoming the new international currency.', 'Detailed analysis of the ongoing supply constraints across global data centers affecting hyperscalers.', 'Tech', 'Sarah Jenkins', 'https://images.unsplash.com/photo-1485827404703-89b55fcc595e?q=80&w=1200'),
			('Central Banks Signal Pause', 'Interest rates may remain steady through Q3.', 'Global markets reacted swiftly to the latest minutes released by the Federal Reserve indicating a prolonged pause.', 'Markets', 'David Chen', 'https://images.unsplash.com/photo-1726566289392-011dc554e604?q=80&w=1200'),
			('A Return to Brutalism in Web Design', 'Why the web is getting uglier, on purpose.', 'An exploration of the anti-design movement pushing back against corporatized, homogenized web experiences.', 'Culture', 'Elena Rostova', 'https://images.unsplash.com/photo-1586125674857-4eb86880905d?q=80&w=1200'),
			('Global Supply Chains Reroute', 'The subtle shifts in shipping lanes are reshaping the global economy.', 'Deep dive into the logistical nightmares and triumphs of the post-2023 shipping industry.', 'World', 'James Thorne', 'https://images.unsplash.com/photo-1583932334951-9a74f88ea6aa?q=80&w=1200'),
			('The End of the Endless Scroll', 'Social media platforms are quietly killing algorithmic feeds.', 'Why chronological feeds are making a comeback among Gen Z users.', 'Culture', 'Maya Patel', 'https://images.unsplash.com/photo-1579532537902-1e50099867b4?q=80&w=1200'),
			('Quantum Advantage Achieved in Simulation', 'A major breakthrough in molecular modeling.', 'Researchers have successfully utilized a 256-qubit system to simulate complex protein folding.', 'Tech', 'Dr. Arthur Vance', 'https://images.unsplash.com/photo-1726569058494-a8e6ddcf1799?q=80&w=1200'),
			('OPEC Plunges Markets into Chaos', 'Unexpected production cuts send crude prices soaring.', 'Analysis of the geopolitical fallout from the latest summit in Vienna.', 'Markets', 'Sarah Jenkins', 'https://images.unsplash.com/photo-1697382608848-c5fea6dd9d60?q=80&w=1200'),
			('Silicon Valleys Silicon Problem', 'The chip fabrication crunch hitting startups.', 'How early-stage hardware companies are finding it impossible to secure manufacturing time.', 'Tech', 'David Chen', 'https://images.unsplash.com/photo-1581092787765-e3feb951d987?q=80&w=1200'),
			('The Silent Extinction of Local News', 'Why the death of the town paper threatens democracy.', 'A sobering look at the expansion of news deserts across rural America and Europe.', 'Opinion', 'Elena Rostova', 'https://images.unsplash.com/photo-1647166545674-ce28ce93bdca?q=80&w=1200'),
			('Europe Votes on Sweeping AI Act', 'The first comprehensive regulatory framework for artificial intelligence.', 'What the new legislation means for foundation models and open-source developers.', 'World', 'Marco Rossi', 'https://images.unsplash.com/photo-1726568313407-c7d9c8a8ce88?q=80&w=1200')
		`)
	}
}
