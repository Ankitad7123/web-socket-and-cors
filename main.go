package main

import (
	"fmt"
	"log"
	"server/controllers"
	"server/models"
	"server/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB 


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Proceed to the next middleware
		c.Next()
	}
}

func main(){
  r := gin.Default()
  dsn := ""
  db , err := gorm.Open(postgres.Open(dsn) , &gorm.Config{})
  if err != nil {
    log.Fatal(err)
  }
  fmt.Print("data base connected")
  if err = db.AutoMigrate(&models.Users{} , &controllers.Message{}); err != nil {
    log.Fatal(err)
  }

 r.Use(CORSMiddleware())

  fmt.Println("migrated")
  routes.UrlPath(r , db)
  r.Run()
}
