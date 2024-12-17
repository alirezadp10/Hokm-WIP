package middleware

import (
    "github.com/alirezadp10/hokm/internal/utils/crypto"
    "github.com/labstack/echo/v4"
    "net/http"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        username := c.QueryParam("username") // Example: Get from query params
        if username == "" {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "username is required"})
        }

        // Decrypt and authenticate
        decryptedUsername, err := crypto.Decrypt([]byte(username))
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid username"})
        }

        // Store the decrypted username in context
        c.Set("username", string(decryptedUsername))

        // Proceed to the next handler
        return next(c)
    }
}
