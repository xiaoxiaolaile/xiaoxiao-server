package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logs "github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func cmdWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)

			break
		}
		log.Printf("recv: %s", message)
		go func() {
			f := &WebsocketIm{
				Faker: Faker{
					Type:     "carry",
					Platform: "xiao_websocket",
					Message:  fmt.Sprintf("%s", message),
					Admin:    true,
				},
				ws:          ws,
				messageType: mt,
			}
			Senders <- f

		}()
		//err = ws.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}

type WebsocketIm struct {
	Faker
	ws          *websocket.Conn
	messageType int
}

func (i *WebsocketIm) Reply(msgs ...interface{}) (arr []string, err error) {
	logs.Info("reply message2", msgs)
	arr, err = i.Faker.Reply(msgs)
	if len(msgs) > 0 {
		//if i.s.f != nil {
		//	i.data["content"] = msgs[0]
		//	i.s.f(i.data)
		//}
		message := fmt.Sprintf("%v", msgs[0])
		if "undefined" == message {
			return
		}
		message = getMessage(msgs)
		err = i.ws.WriteMessage(i.messageType, []byte(message))
		if err != nil {
			log.Println("write:", err)
		}

	}
	return
}
