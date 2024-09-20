package controllers

import (
	"log"
	"net/http"
	"server/models"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type NewDB struct {
DB *gorm.DB
}


var upgrader = websocket.Upgrader{
  CheckOrigin : func(r *http.Request) bool {return true},
}

type Room struct {
  Client map[*websocket.Conn]bool
  Broadcast chan Message
}

var rooms = make(map[string]*Room)
var mutex sync.Mutex






type Message struct{
  gorm.Model
  Room string `json:"room"`
  Username string `json:"username"`
  Message string `json:"message"`
}




func (i *NewDB)CreateUser(c *gin.Context){
    var body models.Users
   if err := c.BindJSON(&body); err != nil {
    c.JSON(http.StatusInternalServerError , gin.H{"err": err.Error()})
    return
  }
  hashed , err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
  if err != nil {
   c.JSON(http.StatusInternalServerError , gin.H{"err":err.Error})
  return  
  }

  body.Password = string(hashed)
  if res := i.DB.Create(&body); res.Error != nil {
    c.JSON(http.StatusInternalServerError , gin.H{"err":res.Error})
    return

  }

  c.String(http.StatusOK , "created")



}
func (i *NewDB)Login(c *gin.Context){
  var LoginD struct {
    Username string `json:"username"`
		Password string `json:"password"`
  }

  if err := c.BindJSON(&LoginD);err != nil {
    c.JSON(http.StatusInternalServerError ,gin.H{"err3":err.Error()})
    return
  }
  var body models.Users
  if err := i.DB.Where("username = ?" ,LoginD.Username).First(&body).Error; err != nil {
    c.JSON(http.StatusInternalServerError ,gin.H{"err1":err})
    return
  }
  res := bcrypt.CompareHashAndPassword([]byte(body.Password) , []byte(LoginD.Password))
  if res != nil {
    c.JSON(http.StatusInternalServerError ,gin.H{"err2":res.Error})
    return
              
  }
  c.String(200  , "loggined")


}


func (i *NewDB)WebsocketHandler(c *gin.Context){
    user1 := c.Query("user1")
  user2 := c.Query("user2")
  users := []string{user1 , user2}
  sort.Strings(users)
  roomN := users[0] + "_" + users[1]

  conn , err := upgrader.Upgrade(c.Writer , c.Request , nil)
  if err != nil {
    log.Fatal("connection didn't happend" , err)
    return
  }

  mutex.Lock()
   room , ok := rooms[roomN]
    if !ok {
    room = &Room{
      Client : make(map[*websocket.Conn]bool),
      Broadcast : make(chan Message),

    }
    rooms[roomN] = room
    go handleWEB(room  )
  }
  room.Client[conn] = true
  mutex.Unlock()

  for {
    var msg Message
    err := conn.ReadJSON(&msg)
    if err != nil {
      log.Print("the read json " , err)
      mutex.Lock()
			delete(room.Client, conn)
			mutex.Unlock()
			conn.Close()
			break
    }
    msg.Room = roomN
    room.Broadcast <- msg

  if msgres := i.DB.Create(&msg); msgres.Error != nil {
      log.Fatal(msgres.Error)
    }
  }
}

func handleWEB(room *Room ){
defer close(room.Broadcast)
  for {
		// Grab the message from the broadcast channel
		msg , ok := <-room.Broadcast

    if !ok {
			log.Println("broadcast channel closed")
			return
		}


 
		// Send it out to every client currently connected
   
		mutex.Lock()
		for client := range room.Client {
			err := client.WriteJSON(msg)
    

			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(room.Client, client)
			}
		}
       
		mutex.Unlock()
	}

}


func (i *NewDB)GetAllmsg(c *gin.Context){
    var msg []Message
     room := c.Param("room")
     if res := i.DB.Where(" room = ?" , room).Order("created_at asc").Find(&msg); res.Error != nil {
     c.JSON(401 , gin.H{"err":res.Error})
    return
  }

  c.JSON(200 , gin.H{"msg":msg})
     
}


