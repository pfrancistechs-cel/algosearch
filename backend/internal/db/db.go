package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectAndMigrate initializes the database pool and executes auto-migrations and seeding.
func ConnectAndMigrate(connStr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid database connection string: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}
	log.Println("[Database] Connected to PostgreSQL (Neon) successfully")

	if err := runMigrations(ctx, pool); err != nil {
		return nil, fmt.Errorf("migration failure: %w", err)
	}

	if err := autoSeedIfEmpty(ctx, pool); err != nil {
		log.Printf("[Database] Seeding warning: %v\n", err)
	}

	return pool, nil
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_trgm;`,
		`CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			category TEXT NOT NULL,
			tags TEXT[] DEFAULT '{}',
			price NUMERIC(10,2) DEFAULT 0,
			popularity_score INT DEFAULT 0,
			search_vector TSVECTOR GENERATED ALWAYS AS (
				setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
				setweight(to_tsvector('english', coalesce(category, '')), 'B') ||
				setweight(to_tsvector('english', coalesce(description, '')), 'C')
			) STORED,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_documents_search_vector ON documents USING gin(search_vector);`,
		`CREATE INDEX IF NOT EXISTS idx_documents_title_trgm ON documents USING gin(title gin_trgm_ops);`,
		`CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(category);`,
		`CREATE INDEX IF NOT EXISTS idx_documents_popularity ON documents(popularity_score DESC);`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed executing query: %s: %w", q, err)
		}
	}

	log.Println("[Database] Schema migrations verified and applied")
	return nil
}

func autoSeedIfEmpty(ctx context.Context, pool *pgxpool.Pool) error {
	var count int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM documents").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		log.Printf("[Database] Found %d existing documents; skipping seed\n", count)
		return nil
	}

	log.Println("[Database] Empty collection detected. Seeding 30 initial sample products...")

	sampleDocs := []struct {
		title       string
		description string
		category    string
		tags        []string
		price       float64
		popularity  int
	}{
		{"Apple MacBook Pro 16-inch M3 Max", "Liquid Retina XDR display with 128GB unified memory and blazing fast Apple Silicon performance.", "Laptops", []string{"apple", "laptop", "m3", "pro"}, 3499.00, 98},
		{"Apple MacBook Air 13-inch M2", "Incredibly thin and light laptop with all-day battery life and stunning Liquid Retina display.", "Laptops", []string{"apple", "lightweight", "m2", "portable"}, 1099.00, 95},
		{"Dell XPS 15 OLED", "High performance creative laptop with InfinityEdge 3.5K OLED touchscreen and NVIDIA RTX 4070 graphics.", "Laptops", []string{"dell", "windows", "oled", "creative"}, 2199.00, 89},
		{"Lenovo ThinkPad X1 Carbon Gen 11", "Ultra-rugged enterprise ultrabook with carbon-fiber chassis and class-leading ergonomic keyboard.", "Laptops", []string{"lenovo", "thinkpad", "business", "rugged"}, 1749.00, 87},
		{"ASUS ROG Zephyrus G14", "Compact 14-inch gaming laptop powered by AMD Ryzen 9 and GeForce RTX 4080 with AniMe Matrix display.", "Laptops", []string{"gaming", "asus", "rog", "portable"}, 1899.00, 91},
		{"Apple iPhone 15 Pro Max", "Aerospace-grade titanium design, A17 Pro chip, customizable Action button, and 5x optical zoom camera.", "Smartphones", []string{"apple", "iphone", "titanium", "5g"}, 1199.00, 99},
		{"Samsung Galaxy S24 Ultra", "Galaxy AI powered flagship phone with built-in S Pen stylus, titanium frame, and 200MP camera system.", "Smartphones", []string{"samsung", "galaxy", "android", "ai"}, 1299.00, 96},
		{"Google Pixel 8 Pro", "Google Tensor G3 chip, advanced generative photography, Super Actua display, and 7 years of OS updates.", "Smartphones", []string{"google", "pixel", "pure-android", "camera"}, 999.00, 90},
		{"Sony WH-1000XM5 Wireless Headphones", "Industry-leading active noise canceling over-ear headphones with two processors and 8 microphones.", "Audio", []string{"sony", "headphones", "anc", "bluetooth"}, 399.00, 94},
		{"Apple AirPods Pro 2nd Gen", "Active noise cancellation, Adaptive Audio, Transparency mode, and USB-C MagSafe charging case.", "Audio", []string{"apple", "earbuds", "wireless", "spatial-audio"}, 249.00, 97},
		{"Bose QuietComfort Ultra", "World-class noise cancellation with revolutionary spatialized audio and ultra-plush leather cushions.", "Audio", []string{"bose", "audio", "anc", "travel"}, 429.00, 88},
		{"Sennheiser Momentum 4 Wireless", "Audiophile-inspired 42mm transducer system delivering pristine acoustic clarity and 60-hour battery life.", "Audio", []string{"sennheiser", "hifi", "battery-beast"}, 349.00, 86},
		{"Logitech MX Master 3S Mouse", "Quiet clicks, 8K DPI any-surface tracking sensor, and ultra-fast electromagnetic MagSpeed scroll wheel.", "Peripherals", []string{"logitech", "mouse", "ergonomic", "productivity"}, 99.00, 96},
		{"Keychron Q1 Pro Custom Keyboard", "Full aluminum wireless custom mechanical keyboard with hot-swappable switches and QMK/VIA support.", "Peripherals", []string{"mechanical", "keychron", "custom", "rgb"}, 199.00, 89},
		{"Apple Magic Keyboard with Touch ID", "Sleek wireless rechargeable keyboard with integrated biometric Touch ID for Mac computers.", "Peripherals", []string{"apple", "keyboard", "touch-id", "minimalist"}, 179.00, 85},
		{"Elgato Stream Deck MK.2", "15 customizable LCD keys for controlling apps, audio, lighting, and streaming macros with visual feedback.", "Peripherals", []string{"elgato", "streaming", "macros", "studio"}, 149.00, 88},
		{"Dell UltraSharp 32-inch 4K USB-C Monitor", "IPS Black technology monitor with 2000:1 contrast ratio, 90W power delivery, and built-in KVM switch.", "Monitors", []string{"dell", "4k", "ips", "usb-c"}, 849.00, 92},
		{"LG UltraGear 27-inch OLED Gaming Monitor", "Lightning fast 240Hz refresh rate and 0.03ms response time with infinite contrast OLED panel.", "Monitors", []string{"lg", "gaming", "oled", "240hz"}, 799.00, 90},
		{"Apple Studio Display 27-inch 5K", "Breathtaking 5K Retina display with 12MP Ultra Wide camera, three-mic array, and six-speaker sound system.", "Monitors", []string{"apple", "5k", "retina", "studio"}, 1599.00, 91},
		{"Apple iPad Pro 12.9-inch M2", "Ultra-fast M2 chip, Liquid Retina XDR display with ProMotion 120Hz, and Apple Pencil hover experience.", "Tablets", []string{"apple", "ipad", "tablet", "creative"}, 1099.00, 93},
		{"Samsung Galaxy Tab S9 Ultra", "Massive 14.6-inch Dynamic AMOLED 2X screen with included S Pen and water-resistant armor aluminum casing.", "Tablets", []string{"samsung", "tablet", "android", "amoled"}, 1199.00, 87},
		{"Apple Watch Ultra 2", "Rugged 49mm titanium case, precision dual-frequency GPS, 3000-nit brightest display, and 36-hour battery.", "Wearables", []string{"apple", "smartwatch", "fitness", "outdoor"}, 799.00, 94},
		{"Garmin Fenix 7 Pro Sapphire Solar", "Multisport GPS smartwatch with solar charging lens, built-in LED flashlight, and topographic maps.", "Wearables", []string{"garmin", "gps", "solar", "rugged"}, 899.00, 88},
		{"Sony PlayStation 5 Console", "Lightning-fast loading with an ultra-high speed SSD, deeper immersion with haptic feedback and 3D Audio.", "Gaming", []string{"sony", "ps5", "console", "4k-gaming"}, 499.00, 97},
		{"Nintendo Switch OLED Model", "Vibrant 7-inch OLED screen, wide adjustable stand, dock with wired LAN port, and 64GB internal storage.", "Gaming", []string{"nintendo", "switch", "handheld", "mario"}, 349.00, 93},
		{"Sonos Move 2 Portable Speaker", "Spacious stereo sound indoors or outdoors with up to 24 hours of playback and IP56 weather resistance.", "Smart Home", []string{"sonos", "smart-speaker", "wifi", "bluetooth"}, 449.00, 89},
		{"Philips Hue White and Color Starter Kit", "Smart lighting ecosystem with millions of rich colors, automated scheduling, and voice assistant sync.", "Smart Home", []string{"smart-home", "lighting", "iot", "ambient"}, 199.00, 86},
		{"Sony Alpha 7 IV Full-frame Mirrorless Camera", "33MP Exmor R CMOS sensor, 4K 60p 10-bit recording, and real-time eye autofocus for creators.", "Cameras", []string{"sony", "mirrorless", "4k", "photography"}, 2498.00, 91},
		{"DJI Mini 4 Pro Drone", "Ultralight sub-249g folding drone with omnidirectional obstacle sensing, 4K/60fps HDR, and 20km transmission.", "Cameras", []string{"dji", "drone", "4k", "aerial"}, 759.00, 92},
		{"Anker Prime 27,650mAh Power Bank (250W)", "Ultra-high output 3-port power bank capable of fast-charging multiple MacBooks and iPhones simultaneously.", "Peripherals", []string{"anker", "charger", "battery", "usb-c"}, 179.00, 91},
	}

	for _, d := range sampleDocs {
		_, err := pool.Exec(ctx, `
			INSERT INTO documents (title, description, category, tags, price, popularity_score)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, d.title, d.description, d.category, d.tags, d.price, d.popularity)
		if err != nil {
			return fmt.Errorf("failed inserting seed record %s: %w", d.title, err)
		}
	}

	log.Printf("[Database] Successfully seeded %d tech products\n", len(sampleDocs))
	return nil
}

