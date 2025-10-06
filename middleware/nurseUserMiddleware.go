// middleware/user_or_nurse.go
package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthUserOrNurse() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token não encontrado",
				"success": false,
			})
			return
		}

		if !strings.HasPrefix(header, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "formato do header Authorization inválido",
				"success": false,
			})
			return
		}

		tokenString := strings.TrimPrefix(header, BearerSchema)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signature method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token inválido",
				"success": false,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token inválido",
				"success": false,
			})
			return
		}

		fmt.Print("role: ", claims["role"])

		userId, _ := claims["sub"].(string)
		role, ok := claims["role"].(string)
		if !ok || (role != "PATIENT" && role != "NURSE"){
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "acesso restrito a usuários ou enfermeiros",
				"success": false,
			})
			return
		}

		hidden, ok := claims["hidden"].(bool)
		if !ok || hidden {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "acesso restrito para usuários ocultos",
				"success": false,
			})
			return
		}

		c.Set("claims", claims)
		c.Set("userId", userId)
		c.Set("role", role)
		c.Set("hidden", hidden)

		c.Next()
	}

}
