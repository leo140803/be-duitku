package config

import (
	"log"
	"os"
	"strings"

	supabase "github.com/lengzuo/supa"
)

var SupaClient *supabase.Client

func InitDB() {
	supabaseProjectID := os.Getenv("SUPABASE_PROJECT_ID")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseProjectID == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_PROJECT_ID and SUPABASE_ANON_KEY environment variables are required")
	}

	// Clean the project ID (remove any URL parts if accidentally included)
	supabaseProjectID = strings.TrimSpace(supabaseProjectID)
	if strings.Contains(supabaseProjectID, "https://") {
		// Extract project ID from URL if full URL was provided
		parts := strings.Split(supabaseProjectID, ".")
		if len(parts) > 0 {
			supabaseProjectID = parts[0]
		}
	}

	log.Printf("Connecting to Supabase project: %s", supabaseProjectID)

	client, err := supabase.New(supabase.Config{
		ApiKey:     supabaseKey,
		ProjectRef: supabaseProjectID,
	})
	if err != nil {
		log.Fatal("Failed to create Supabase client:", err)
	}

	SupaClient = client
	log.Println("Supabase connected 🚀")
}
