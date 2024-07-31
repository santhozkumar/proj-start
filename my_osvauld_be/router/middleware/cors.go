package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CorsMiddleware() func(*gin.Context) {

	return func(ctx *gin.Context) {

		allowedOrigins := viper.GetString("AC_ALLOW_ORIGINS")
        log.Println(allowedOrigins)

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
		ctx.Writer.Header().Set("Access-Control-Max-Age", viper.GetString("AC_MAX_AGE"))
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", viper.GetString("AC_ALLOW_METHODS"))
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", viper.GetString("AC_ALLOW_HEADERS"))
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", viper.GetString("AC_ALLOW_CREDENTIALS"))
		ctx.Writer.Header().Set("Cache-Control", "no-cache")

        if ctx.Request.Method == http.MethodOptions {
            // Must return 2xx status for options method, 200 or 204 is fine
            // 204 seems didn't work in the old firefox browsers, so use 200
            ctx.AbortWithStatus(http.StatusOK)
            return
        }
        ctx.Next()
	}

}
