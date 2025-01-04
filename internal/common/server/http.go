package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yuyayang02/BlogLite-api/internal/common/auth"
	"github.com/yuyayang02/BlogLite-api/internal/common/config"
	"log"
	"net/http"
)

func RunHTTPServer(port int, createHandler func(router gin.IRouter)) {

	if config.Cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	root := gin.New()
	setMiddlewares(root)

	createHandler(root.Group("/api/v1"))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server running %s", addr)

	err := http.ListenAndServe(addr, root)
	if err != nil {
		log.Fatal("Unable to start HTTP server", "err", err)
	}
}

func setMiddlewares(r *gin.Engine) {
	r.Use(gin.Recovery())

	addCorsMiddleware(r)
	addAuthMiddleware(r)

	r.Use(func(c *gin.Context) {
		//这个头部指示浏览器在处理响应时不要尝试猜测内容类型。它要求浏览器遵循服务器指定的 MIME 类型，从而防止某些类型的攻击（例如，恶意代码通过混淆的内容类型执行）
		c.Header("X-Content-Type-Options", "nosniff")
		//这个头部防止页面被嵌入到 <frame> 或 <iframe> 中，从而抵御点击劫持（clickjacking）攻击。设置为 deny 意味着页面无法被任何网站嵌入
		c.Header("X-Frame-Options", "deny")
		c.Next()
	})

	//r.Use(noCache())
}

func addAuthMiddleware(r *gin.Engine) {
	r.Use(auth.AuthMiddleware())
}

func addCorsMiddleware(r *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.New(corsConfig))
}

func noCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
