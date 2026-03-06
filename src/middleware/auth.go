package middleware

import (
	"automation-developer-guide/src/utils"

	"github.com/gofiber/fiber/v2"
)

// IsAuthenticated checks if the user is logged in by parsing the JWT from the session cookie.
// This replaces the old AuthProxyMiddleware that made HTTP calls to the auth service.
func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("session_id")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	claims, err := utils.ParseJWT(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// Store username in locals for easy access in handlers
	if username, ok := claims["username"].(string); ok {
		c.Locals("username", username)
	}

	// Store user_id in locals
	if userID, ok := claims["user_id"].(string); ok {
		c.Locals("user_id", userID)
	}

	if email, ok := claims["email"].(string); ok {
		c.Locals("email", email)
	}

	if avatarURL, ok := claims["avatar_url"].(string); ok {
		c.Locals("avatar_url", avatarURL)
	}

	return c.Next()
}
