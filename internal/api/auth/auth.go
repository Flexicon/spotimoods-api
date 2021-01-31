package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// TokenOptions used for generating JWT tokens
type TokenOptions struct {
	DisplayName string
	Email       string
}

// GenerateToken prepares a new JWT based on input options
func GenerateToken(opts TokenOptions) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["display_name"] = opts.DisplayName
	claims["email"] = opts.Email
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	return token.SignedString([]byte(viper.GetString("app.secret")))
}

// Options for auth middleware
type Options struct {
	Services *internal.ServiceProvider
}

// UseMiddleware and attach it to the given echo.Group
func UseMiddleware(g *echo.Group, opts Options) {
	g.Use(middleware.JWT([]byte(viper.GetString("app.secret"))), setAuthUser(opts))
}

// setAuthUser middleware to verify an existing user for a token
func setAuthUser(opts Options) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			user, err := opts.Services.User().FindByEmail(claims["email"].(string))
			if err != nil {
				log.Println("failed to retrieve user by email:", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve user by email")
			}
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set("user", user)
			if spotifyToken, err := opts.Services.User().FindTokenForUser(user.ID); err == nil {
				c.Set("user.spotify_token", spotifyToken)
			}

			return next(c)
		}
	}
}
