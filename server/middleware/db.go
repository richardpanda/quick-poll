package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
