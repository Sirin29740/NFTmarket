package api

import (
	"NFTmarket/api/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Router() {
	r := gin.Default()

	// -----------------------------------------------------------------
	// 1. CORS 配置 (解决跨域问题)
	// -----------------------------------------------------------------
	r.Use(cors.New(cors.Config{
		// 假设前端 Vue 开发服务器在 5173 端口运行
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 必须允许 Authorization 头部，用于传输 JWT Token
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	const frontendDistDir = "./frontend/nftmarket/dist"

	// -----------------------------------------------------------------
	// 1. API 路由 (必须优先注册)
	// -----------------------------------------------------------------
	// 公开路由
	r.POST("/register", Register)
	r.POST("/login", Login)

	// 认证路由组
	authGroup := r.Group("/api")
	authGroup.Use(AuthMiddleware())
	{
		// 确保这个具体的 API 路径在 StaticFS 之前
		authGroup.GET("/profile", v1.GetProfile)
		authGroup.POST("/upload", v1.Upload)
	}

	r.GET("/", func(c *gin.Context) {
		c.File(frontendDistDir + "/index.html")
	})

	r.Static("/assets", frontendDistDir+"/assets")

	// -----------------------------------------------------------------
	// 3. SPA 历史模式回退路由 (NoRoute)
	// -----------------------------------------------------------------

	// NoRoute 捕获所有未被 StaticFS 或 API 匹配到的路径。
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只有非 API 路径，才返回 index.html
		if !strings.HasPrefix(path, "/api/") {
			// 返回 index.html
			c.File(frontendDistDir + "/index.html")
			return
		}

		// 未知 API 路由
		c.JSON(http.StatusNotFound, gin.H{"message": "API endpoint not found"})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
