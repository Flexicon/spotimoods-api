package api

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

type authController struct {
	services *internal.ServiceProvider
}

func newAuth(services *internal.ServiceProvider) Controller {
	return &authController{services: services}
}

func (h *authController) Routes(g *echo.Group) {
	g = g.Group("/auth", authMiddlewareChain(Options{Services: h.services})...)

	g.GET("/refresh", h.Refresh())
}

func (h *authController) Refresh() echo.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}

	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		token, err := generateToken(TokenOptions{
			DisplayName: user.DisplayName,
			Email:       user.Email,
		})
		if err != nil {
			log.Println("Failed to generate token:", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, response{Token: token})
	}
}

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

func authMiddlewareChain(opts Options) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.JWT([]byte(viper.GetString("app.secret"))),
		authUser(opts),
	}
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

			return next(c)
		}
	}
}
