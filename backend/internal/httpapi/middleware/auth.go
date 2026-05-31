package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/auth"
)

const claimsKey = "claims"

// RequireAuth enforces bearer token authentication and stores parsed claims in the request context.
func RequireAuth(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "缺少登录凭证，请重新登录"})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := authService.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录状态无效，请重新登录"})
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

// RequirePortalRole restricts a route group to a specific distributor portal role.
func RequirePortalRole(role auth.PortalRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok || claims.PortalRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权访问当前页面"})
			return
		}
		c.Next()
	}
}

// GetClaims returns the parsed JWT claims previously stored by RequireAuth.
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get(claimsKey)
	if !ok {
		return nil, false
	}
	claims, ok := value.(*auth.Claims)
	return claims, ok
}
