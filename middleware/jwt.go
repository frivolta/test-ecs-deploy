package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scopes      string   `json:"scopes"`
	Sub         string   `json:"sub"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

var C CustomClaims

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "https://apecalendar-dev.com"
		convAud, ok := token.Claims.(jwt.MapClaims)["aud"].([]interface{})
		if !ok {
			strAud, ok := token.Claims.(jwt.MapClaims)["aud"].(string)
			if !ok {
				return token, errors.New("Invalid audience.")
			}
			if strAud != aud {
				return token, errors.New("Invalid audience.")
			}
		} else {
			for _, v := range convAud {
				if v == aud {
					break
				} else {
					return token, errors.New("Invalid audience.")
				}
			}
		}
		iss := "https://dev-iv8n6772.us.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}
		// Setting custom claims
		C = CustomClaims{
			Scopes: token.Claims.(jwt.MapClaims)["scope"].(string),
			Sub:    token.Claims.(jwt.MapClaims)["sub"].(string),
			Email:  token.Claims.(jwt.MapClaims)["https://example.com/email"].(string),
		}

		// Converting permissions to []string
		permissions := token.Claims.(jwt.MapClaims)["permissions"].([]interface{})
		s := make([]string, len(permissions))
		for i, v := range permissions {
			s[i] = fmt.Sprint(v)
		}
		C.Permissions = s
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMid := *jwtMiddleware
		if err := jwtMid.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
			// Error handling here
		}
		c.Set("scopes", C.Scopes)
		c.Set("sub", C.Sub)
		c.Set("email", C.Email)
		c.Set("permissions", C.Permissions)
	}
}

func CheckScope(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Value("permissions").([]string)[:]
		for _, v := range permissions {
			ok := slices.Contains(p, v)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you don't have privileges to access this resource"})
			}
		}
		c.Next()
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + "dev-iv8n6772.us.auth0.com" + "/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}
	return cert, nil
}
