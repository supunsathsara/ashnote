package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/supunsathsara/ashnote/internal/db"
	"github.com/supunsathsara/ashnote/internal/handler"
)

func main() {
	// Initialize template engine with correct path
	workDir, _ := os.Getwd()
	// If running from the root directory, use "web/templates"
	// If running from cmd directory, use "../web/templates"
	templatesDir := "web/templates"
	if strings.HasSuffix(workDir, "cmd") {
		templatesDir = "../web/templates"
	}

	// Ensure the templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		log.Printf("Templates directory not found at: %s", templatesDir)
		log.Printf("Current working directory: %s", workDir)
		// Try the absolute path as a fallback
		templatesDir = ""
		if strings.HasSuffix(workDir, "cmd") {
			templatesDir = filepath.Join(workDir, "../web/templates")
		} else {
			templatesDir = filepath.Join(workDir, "web/templates")
		}
		log.Printf("Using fallback templates path: %s", templatesDir)
	}

	engine := html.New(templatesDir, ".html")

	// Initialize Fiber app with custom config
	app := fiber.New(fiber.Config{
		// Increase the header size limit to 10MB
		ReadBufferSize: 10 * 1024 * 1024,
		// Increase body limit to 10MB as well
		BodyLimit: 10 * 1024 * 1024,
		// Set up template engine
		Views: engine,
		// Custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Set up middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize database connection
	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize handler
	h := handler.New(database)

	// Register API routes
	h.RegisterRoutes(app)

	// Serve index page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// Serve static assets (if any)
	app.Static("/static", "web/static")

	// Print all routes for debugging
	log.Println("Available Routes:")
	for _, route := range app.GetRoutes() {
		log.Printf("%s %s", route.Method, route.Path)
	}

	// Start server
	log.Println("Starting AshNote on http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
