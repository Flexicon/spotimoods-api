package api

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// TokenOptions used for generating JWT tokens
type TokenOptions struct {
	DisplayName string
	Email       string
}

// generateToken prepares a new JWT based on input options
func generateToken(opts TokenOptions) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["display_name"] = opts.DisplayName
	claims["email"] = opts.Email
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	return token.SignedString([]byte(viper.GetString("app.secret")))
}

func useAuthMiddleware(g *echo.Group, opts Options) {
	g.Use(middleware.JWT([]byte(viper.GetString("app.secret"))), authUser(opts))
}

// authUser middleware to verify an existing user for a token
func authUser(opts Options) echo.MiddlewareFunc {
	type response struct {
		Msg string `json:"message"`
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			user, err := opts.Services.User().FindByEmail(claims["email"].(string))
			if err != nil {
				log.Println("Failed to retrieve user by email:", err)
				return c.JSON(http.StatusInternalServerError, response{Msg: "something went wrong"})
			}
			if user == nil {
				return c.JSON(http.StatusUnauthorized, response{Msg: "unauthorized"})
			}

			c.Set("user", user)
			if spotifyToken, err := opts.Services.User().FindTokenForUser(user.ID); err == nil {
				c.Set("user.spotify_token", spotifyToken)
			}

			return next(c)
		}
	}
}
