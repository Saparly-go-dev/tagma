package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, role, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	c.Set(userCtx, userId)
	c.Set(roleCtx, role)
	c.Next()
}

func (h *Handler) requireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, ok := c.Get(roleCtx)
		if !ok {
			newErrorResponse(c, http.StatusUnauthorized, "user role not found")
			return
		}
		roleName, ok := role.(string)
		if !ok {
			newErrorResponse(c, http.StatusUnauthorized, "invalid user role")
			return
		}
		if _, ok := allowed[roleName]; !ok {
			newErrorResponse(c, http.StatusForbidden, "insufficient permissions")
			return
		}
		if c.Request.Method != http.MethodGet && roleName != "admin" && !canMutate(roleName, c.Request.URL.Path) {
			newErrorResponse(c, http.StatusForbidden, "insufficient permissions")
			return
		}
		c.Next()
	}
}

func canMutate(role, path string) bool {
	allowedPrefixes := map[string][]string{
		"agent":        {"/api/orders", "/api/payments", "/api/extra"},
		"ekspeditor":   {"/api/orders", "/api/payments", "/api/ekspeditors/save", "/api/ekspeditors/product/delivered"},
		"merchandiser": {"/api/furniture"},
	}
	for _, prefix := range allowedPrefixes[role] {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func getUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get(roleCtx)
	if !ok {
		return "", errors.New("user role not found")
	}
	roleName, ok := role.(string)
	if !ok {
		return "", errors.New("invalid user role")
	}
	return roleName, nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not invalid type")
		return 0, errors.New("user id is not invalid type")
	}

	return idInt, nil
}
