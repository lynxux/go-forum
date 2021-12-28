package routes

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"web_app/controllers"
	"web_app/middleware"
)

func Setup(mode string) *gin.Engine {

	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //只有当mode=release时把gin设置为release模式
	}

	//初始化gin内置的校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Println("init gin validator translation failed! err: ", err)
		return nil
	}

	//r := gin.New()
	r := gin.Default()

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	v1 := r.Group("/api/v1")

	//logger, _ := zap.NewProduction()
	//设置中间件，这里使用zap库以及ginzap
	v1.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	v1.Use(ginzap.RecoveryWithZap(zap.L(), true))
	//v1.Use(middleware.RateLimitMiddleware(time.Second*2, 1)) //ratelimit

	//注册业务
	v1.POST("/signup", controllers.SignUpHandler)
	//登录业务
	v1.POST("/login", controllers.LoginHandler)

	//认证中间件
	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.PostDetailHandler)
		v1.GET("/posts", controllers.PostListHandler)

		//
		v1.GET("/post2", controllers.PostListHandler2)

		v1.POST("/vote", controllers.VotePostHandler)

	}

	//v1.GET("/ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
	//	// 如果是登录的用户，判断请求头中是否含有有效的JWT Token
	//	c.String(http.StatusOK, "pong")
	//})

	//版本号页面
	v1.GET("/version", func(c *gin.Context) {
		c.String(200, "to do")
	})

	pprof.Register(r) //注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
