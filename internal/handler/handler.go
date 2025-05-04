package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supunsathsara/ashnote/internal/crypto"
	"github.com/supunsathsara/ashnote/internal/db"
)

// Handler contains all the dependencies needed for the handlers
type Handler struct {
	DB *db.DB
}

// New creates a new handler
func New(db *db.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// RegisterRoutes registers all the routes for the application
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Create a new message
	app.Post("/api/messages", h.CreateMessage)

	// Get a message
	app.Get("/api/messages/:id", h.GetMessage)

	// Serve static files and handle frontend routes
	app.Get("/message/:id", h.ViewMessagePage)
}

// CreateMessage handles the creation of a new message
func (h *Handler) CreateMessage(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Message  string `json:"message"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Message == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message and password are required",
		})
	}

	// Encrypt the message
	encrypted, err := crypto.Encrypt(req.Message, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to encrypt message",
		})
	}

	// Store the encrypted message in the database
	id, err := h.DB.StoreMessage(encrypted)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store message",
		})
	}

	// Return the message ID
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      id,
		"message": "Message created successfully",
		"url":     c.BaseURL() + "/message/" + id,
	})
}

// GetMessage handles retrieving a message
func (h *Handler) GetMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	password := c.Query("password")

	// Validate input
	if id == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message ID and password are required",
		})
	}

	// Get the encrypted message from the database
	encryptedMessage, err := h.DB.GetMessage(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Message not found or already accessed",
		})
	}

	// Decrypt the message
	decrypted, err := crypto.Decrypt(encryptedMessage, password)
	if err != nil {
		// If decryption fails, reset the access count
		// This allows users to retry with the correct password
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	// Delete the message after successful decryption
	if err := h.DB.DeleteMessage(id); err != nil {
		// Log the error but don't fail the request
		// The message has already been marked as accessed
	}

	// Return the decrypted message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": decrypted,
	})
}

// ViewMessagePage serves the HTML page for viewing a message
func (h *Handler) ViewMessagePage(c *fiber.Ctx) error {
	// This is just passing the request to the frontend
	// The actual HTML template will be served
	return c.Render("message", fiber.Map{
		"ID": c.Params("id"),
	})
}
