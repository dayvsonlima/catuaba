package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRoutes(t *testing.T) {
	dir := t.TempDir()

	routesCode := `package config

import (
	"myapp/app/controllers/auth"
	"myapp/app/controllers/products"
	"myapp/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/me", middleware.RequireAuth(), auth.Me)

	api := r.Group("/api")
	api.GET("/products", products.Index)
	api.POST("/products", products.Create)
	api.GET("/products/:id", products.Show)
	api.PUT("/products/:id", products.Update)
	api.DELETE("/products/:id", products.Delete)
}
`
	filePath := filepath.Join(dir, "routes.go")
	require.NoError(t, os.WriteFile(filePath, []byte(routesCode), 0644))

	routes, err := ParseRoutes(filePath)
	require.NoError(t, err)
	assert.Len(t, routes, 8)

	assert.Equal(t, "POST", routes[0].Method)
	assert.Equal(t, "/auth/register", routes[0].Path)
	assert.Equal(t, "auth.Register", routes[0].Handler)

	// Route with middleware
	assert.Equal(t, "GET", routes[2].Method)
	assert.Equal(t, "/auth/me", routes[2].Path)
	assert.Equal(t, "auth.Me", routes[2].Handler)
	assert.Equal(t, []string{"middleware.RequireAuth()"}, routes[2].Middleware)

	assert.Equal(t, "GET", routes[3].Method)
	assert.Equal(t, "/products", routes[3].Path)
	assert.Equal(t, "products.Index", routes[3].Handler)

	assert.Equal(t, "DELETE", routes[7].Method)
	assert.Equal(t, "/products/:id", routes[7].Path)
}

func TestParseRoutes_NoFile(t *testing.T) {
	_, err := ParseRoutes("/nonexistent/routes.go")
	assert.Error(t, err)
}
