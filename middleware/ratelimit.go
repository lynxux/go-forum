package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

// 每fillInterval(如果是数字表示纳秒)添加cap个令牌
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		//如果取不到令牌就返回
		if bucket.TakeAvailable(1) == 0 { //返回可用的令牌
			c.String(200, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
