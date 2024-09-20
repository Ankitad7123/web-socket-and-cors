package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func UrlPath(r *gin.Engine , db *gorm.DB){
  views := controllers.NewDB{DB:db}   
  
  r.POST("/" , views.CreateUser)
 r.POST("/login" , views.Login)
  r.GET("/ws" , views.WebsocketHandler)
  r.GET("/msgs/:room" , views.GetAllmsg)

 }
