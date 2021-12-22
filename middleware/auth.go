package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"web_app/controllers"
	"web_app/pkg/jwt"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//客户端携带token有三种方式：1.放在请求头 2.放在请求体中 3. 放在URI
		//这里假设Token放在Header的Authorization中，并使用Bearer开头  Authorization: Bearer xxx.xxx.xx
		//具体的实现方式依据业务而定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			c.Abort()
			return
		}

		//空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}

		//parts[1]是获取的tokenstring
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前的userID保存到请求的上下文中
		c.Set(controllers.CtxUserIDKey, mc.UserID)
		c.Next() //后续的处理中，我们可以使用c.get("CtxUserIDKey")来获取当前的用户信息
	}
}
