package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"alc/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("=== JR del Perú - Database Seeder ===")

	// Get database URL
	dbURL := os.Getenv("POSTGRESQL_URL")
	if dbURL == "" {
		log.Fatal("❌ POSTGRESQL_URL environment variable is required")
	}

	// Get admin credentials from environment
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminEmail == "" || adminPassword == "" {
		log.Fatal("❌ ADMIN_USERNAME, ADMIN_EMAIL, and ADMIN_PASSWORD environment variables are required")
	}

	// Connect to database
	log.Println("📦 Connecting to database...")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Verify database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Unable to ping database: %v\n", err)
	}
	log.Println("✓ Database connection successful")

	// Initialize repository
	queries := repository.New(dbpool)

	// Check if admin already exists
	log.Printf("🔍 Checking if admin user '%s' already exists...\n", adminUsername)
	existingAdmin, err := queries.GetAdminByUsername(context.Background(), adminUsername)
	if err == nil {
		log.Printf("⚠️  Admin user '%s' already exists (ID: %d)\n", existingAdmin.Username, existingAdmin.AdminID)
		log.Println("ℹ️  Skipping admin creation. Delete the existing user first if you want to recreate it.")
		return
	}

	// Hash password
	log.Println("🔐 Hashing password...")
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Error hashing password: %v\n", err)
	}

	// Create admin user
	log.Printf("👤 Creating admin user '%s'...\n", adminUsername)
	admin, err := queries.CreateAdmin(context.Background(), repository.CreateAdminParams{
		Username:     adminUsername,
		Email:        adminEmail,
		PasswordHash: string(hash),
	})
	if err != nil {
		log.Fatalf("❌ Error creating admin: %v\n", err)
	}

	// Success message
	fmt.Println()
	log.Println("✓ Admin user created successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  ID:       %d\n", admin.AdminID)
	fmt.Printf("  Username: %s\n", admin.Username)
	fmt.Printf("  Email:    %s\n", admin.Email)
	fmt.Printf("  Active:   %v\n", admin.IsActive.Bool)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	log.Println("🎉 Seeder completed successfully!")
	log.Printf("📍 You can now login at http://localhost:8080/admin/login\n")
}
