package middleware

import (
    "net/http"

    "github.com/alirezadp10/hokm/internal/util/crypto"
    "github.com/alirezadp10/hokm/internal/util/trans"
    sqliteRepo "github.com/alirezadp10/hokm/pkg/repository/sqlite"
    "github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
    playerRepo *sqliteRepo.PlayersRepository
}

func NewAuthMiddleware(playerRepo *sqliteRepo.PlayersRepository) *AuthMiddleware {
    return &AuthMiddleware{
        playerRepo: playerRepo,
    }
}

func (m *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        username := c.Request().Header.Get("user-reference-key")

        if username == "" {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": trans.Get("user-reference-key is required.")})
        }

        // Decrypt and authenticate
        decryptedUsername, err := crypto.Decrypt(username)

        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": trans.Get("invalid user-reference-key.")})
        }

        if !m.playerRepo.CheckPlayerExistence(decryptedUsername) {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": trans.Get("invalid user-reference-key.")})
        }

        // Store the decrypted username in context
        c.Set("username", decryptedUsername)

        // Proceed to the next handler
        return next(c)
    }
}
