package auth

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Middleware handles authentication for routes requiring authentication
type Middleware interface {
	APIKeyMiddlewareFunc() gin.HandlerFunc
	GoogleJWTMiddlewareFunc() gin.HandlerFunc
	GinJWTMiddleware(authenticator func(c *gin.Context) (interface{}, error)) (middleware *jwt.GinJWTMiddleware, err error)
	GinJWTMiddlewareForClientLogin(authenticator func(c *gin.Context) (interface{}, error)) (middleware *jwt.GinJWTMiddleware, err error)
}

// NewAuthMiddleware returns a new auth.AuthMiddleware
func NewAuthMiddleware(config *config.APIConfig) (authMiddleware Middleware) {
	authMiddleware = &authMiddlewareImpl{
		config: config,
	}

	return
}

type authMiddlewareImpl struct {
	config *config.APIConfig
}

func (m *authMiddlewareImpl) APIKeyMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := c.GetHeader("Authorization")
		bearerToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		// check if bearer token equals api key, otherwise check if it's a valid jwt
		if bearerToken != m.config.Auth.APIKey {
			_, err := ValidateJWT(m.config, bearerToken)
			if err != nil {
				log.Error().
					Str("authorizationHeader", authorizationHeader).
					Msg("Authorization header bearer token is incorrect")
				c.Status(http.StatusUnauthorized)
				return
			}
		}

		// set 'user' to enforce a handler method to require api key auth with `user := c.MustGet(gin.AuthUserKey).(string)` and ensuring the user equals 'apiKey'
		c.Set(gin.AuthUserKey, "apiKey")
	}
}

func (m *authMiddlewareImpl) GoogleJWTMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			return
		}

		bearerTokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		valid, err := isValidGoogleJWT(bearerTokenString)

		if err != nil {
			log.Warn().Err(err).Str("bearer", bearerTokenString).Msgf("Error when validating Google JWT")
			return
		}

		if !valid {
			log.Warn().Err(err).Str("bearer", bearerTokenString).Msgf("Google JWT is not valid")
			return
		}

		// set 'user' to enforce a handler method to require api key auth with `user := c.MustGet(gin.AuthUserKey).(string)` and ensuring the user equals 'apiKey'
		c.Set(gin.AuthUserKey, "google-jwt")
	}
}

func (m *authMiddlewareImpl) coreGinJWTMiddleware(authenticator func(c *gin.Context) (interface{}, error)) (middleware *jwt.GinJWTMiddleware, err error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:         m.config.Auth.JWT.Domain,
		Key:           []byte(m.config.Auth.JWT.Key),
		TokenLookup:   "header:Authorization, cookie:jwt",
		Authenticator: authenticator,
		TimeFunc:      time.Now,
	})
}

func (m *authMiddlewareImpl) GinJWTMiddleware(authenticator func(c *gin.Context) (interface{}, error)) (middleware *jwt.GinJWTMiddleware, err error) {
	middleware, err = m.coreGinJWTMiddleware(authenticator)
	if err != nil {
		return nil, err
	}

	// send cookie
	middleware.SendCookie = true
	middleware.SecureCookie = true
	middleware.CookieHTTPOnly = true
	middleware.CookieDomain = m.config.Auth.JWT.Domain

	// redirect after login
	middleware.LoginResponse = func(c *gin.Context, code int, token string, expire time.Time) {
		// see if gin context has a return url
		returnURL, exists := c.Get("returnURL")
		if exists {
			c.Redirect(http.StatusFound, returnURL.(string))
			return
		}

		// cookie is used, so token does not need to be returned via response
		c.Redirect(http.StatusFound, "/")
	}

	// redirect after logout
	middleware.LogoutResponse = func(c *gin.Context, code int) {
		c.Redirect(http.StatusFound, "/login")
	}

	// set some user properties as claims
	middleware.PayloadFunc = func(data interface{}) jwt.MapClaims {
		// add user properties as claims
		if user, ok := data.(*contracts.User); ok {

			organizations := []string{}
			for _, o := range user.Organizations {
				organizations = append(organizations, o.Name)
			}

			groups := []string{}
			for _, g := range user.Groups {
				groups = append(groups, g.Name)
			}

			return jwt.MapClaims{
				jwt.IdentityKey: user.ID,
				"email":         user.GetEmail(),
				"roles":         user.Roles,
				"groups":        groups,
				"organizations": organizations,
			}
		}
		return jwt.MapClaims{}
	}

	return middleware, nil
}

func (m *authMiddlewareImpl) GinJWTMiddlewareForClientLogin(authenticator func(c *gin.Context) (interface{}, error)) (middleware *jwt.GinJWTMiddleware, err error) {
	middleware, err = m.coreGinJWTMiddleware(authenticator)
	if err != nil {
		return nil, err
	}

	// set some client properties as claims
	middleware.PayloadFunc = func(data interface{}) jwt.MapClaims {
		// add client properties as claims
		if client, ok := data.(*contracts.Client); ok {
			return jwt.MapClaims{
				jwt.IdentityKey: client.ID,
				"clientID":      client.ClientID,
				"roles":         client.Roles,
			}
		}
		return jwt.MapClaims{}
	}

	return middleware, nil

}
