package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/pkg/helpers"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/sirupsen/logrus"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", constants.AllowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", constants.AllowCredential)
		c.Writer.Header().Set("Access-Control-Allow-Headers", constants.AllowHeader)
		c.Writer.Header().Set("Access-Control-Allow-Methods", constants.AllowMethods)
		c.Writer.Header().Set("Access-Control-Max-Age", constants.MaxAge)

		if !helpers.IsArrayContains(strings.Split(constants.AllowMethods, ", "), c.Request.Method) {
			logger.InfoF("method %s is not allowed\n", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryCORS}, c.Request.Method)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error method": "forbidden with CORS policy"})
			return
		}

		for key, value := range c.Request.Header {
			if !helpers.IsArrayContains(strings.Split(constants.AllowHeader, ", "), key) {
				logger.InfoF("ini header %s: %s\n", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryCORS}, key, value)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error header": "forbidden with CORS policy"})
				return
			}
		}

		if constants.AllowOrigin != "*" {
			if !helpers.IsArrayContains(strings.Split(constants.AllowOrigin, ", "), c.Request.Host) {
				logger.InfoF("host '%s' is not part of '%v'\n", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryCORS}, c.Request.Host, constants.AllowOrigin)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error origin": "forbidden with CORS policy"})
				return
			}
		}

		c.Next()
	}
}
