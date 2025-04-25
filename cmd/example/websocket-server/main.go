package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"mymodule/contract/golang/notification"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go ReadMsg(conn)

	go SendMsg(conn)
	//for {
	//	messageType, p, err := conn.ReadMessage()
	//	fmt.Println("message:", string(p))
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	fmt.Println("enter message:")
	//	scanner := bufio.NewScanner(os.Stdin)
	//	scanner.Scan()
	//	msg := scanner.Text()
	//	if err := conn.WriteMessage(messageType, []byte(msg)); err != nil {
	//		log.Println(err)
	//		return
	//	}
	//}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ff", TestEncode())
	fmt.Println("websokcet runned")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ReadMsg(conn *websocket.Conn) {
	for {
		oType, msg, rErr := conn.ReadMessage()

		notifMsg := Protodecoder(string(msg))

		if rErr != nil {
			panic(rErr)
		}
		fmt.Println("msg:", notifMsg)
		fmt.Println("type:", oType)
	}
}

func SendMsg(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("enter message:")
		scanner.Scan()
		msg := scanner.Text()
		if wErr := conn.WriteMessage(websocket.TextMessage, []byte(msg)); wErr != nil {
			panic(wErr)
		}
	}
}

func Protodecoder(msg string) notification.Notification {
	decodedMsg, dErr := base64.StdEncoding.DecodeString(msg)
	if dErr != nil {
		log.Fatal(dErr)

		return notification.Notification{}
	}

	var notif notification.Notification
	uErr := proto.Unmarshal(decodedMsg, &notif)
	if uErr != nil {
		log.Fatal(uErr)

		return notification.Notification{}
	}

	return notif
}

func TestEncode() string {
	msg := notification.Notification{
		Event:   "matchingPlayer",
		Payload: "payload msg",
	}

	protoMsg, mErr := proto.Marshal(&msg)
	if mErr != nil {
		panic(mErr)
	}

	return base64.StdEncoding.EncodeToString(protoMsg)

}
