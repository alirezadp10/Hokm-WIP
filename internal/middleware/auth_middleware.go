package middleware

import (
    "github.com/alirezadp10/hokm/internal/database/sqlite"
    "github.com/alirezadp10/hokm/internal/util/crypto"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "net/http"
)

func AuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            username := c.Request().Header.Get("user-reference-key")

            if username == "" {
                return c.JSON(http.StatusBadRequest, map[string]string{"error": "user-reference-key is required"})
            }

            // Decrypt and authenticate
            decryptedUsername, err := crypto.Decrypt(username)

            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user-reference-key"})
            }

            if !sqlite.CheckPlayerExistence(db, decryptedUsername) {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user-reference-key"})
            }

            // Store the decrypted username in context
            c.Set("username", decryptedUsername)

            // Proceed to the next handler
            return next(c)
        }
    }
}
