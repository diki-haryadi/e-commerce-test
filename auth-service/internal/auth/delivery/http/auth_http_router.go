package articleHttpController

import (
	"github.com/diki-haryadi/go-micro-template/config"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/labstack/echo/v4"
	"net/http"

	authDomain "github.com/diki-haryadi/go-micro-template/internal/auth/domain"
)

type Router struct {
	controller authDomain.HttpController
}

func NewRouter(controller authDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {
	e.POST("/auth/signup", r.controller.SignUp)
	e.POST("/auth/signin", r.controller.SignIn)
	e.POST("/auth/refresh", r.controller.RefreshToken)

	protected := e.Group("")
	protected.Use(r.JWTMiddleware())
	protected.GET("/me", r.controller.GetProfile)
}

// Custom middleware to extract and validate JWT token
func (r *Router) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			// Remove 'Bearer ' prefix if present
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			// Validate token and extract claims
			claims, err := pkg.ValidateToken(tokenString, config.BaseConfig.App.JWTSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Set claims to context
			c.Set("claims", claims)
			c.Set("userId", claims.Sub)
			c.Set("scope", claims.Scope)

			return next(c)
		}
	}
}
