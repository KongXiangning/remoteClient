package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"remoteClient/wsConn"
	"log"
	"net/http"
	"net/url"
	"time"
)

//var origin = "http://localhost:7777"
//var url = "ws://localhost:7777/ws"

var addr = flag.String("addr", "192.168.50.122:7777", "http service address")

func main()  {
	var (
		wsClient *websocket.Conn
		conn *wsConn.Connection
		err error
		u url.URL
		requestHeader http.Header
	)


	flag.Parse()
	log.SetFlags(0)
	u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	requestHeader = http.Header{}
	requestHeader.Add("test","1asdafsdf")

ERR:
	if wsClient,_,err = websocket.DefaultDialer.Dial(u.String(), requestHeader);err != nil{
		log.Println(err)
		time.Sleep(10*time.Second)
		goto  ERR
	}

	if conn,err = wsConn.Init(wsClient);err != nil{
		log.Println(err)
		wsClient.Close()
		goto ERR
	}

	log.Println("connection success!")

	go func() {
		var (
			err error
		)
		for {
			//b  := []byte("heartbeatheartbeatdfasdheartbeatsdsasdheartbeatheartbeatheartbeatsdfasdadsaheartbeat")
			//out,_ := coder.GzipEconder(b)
			/*fmt.Print("b:")
			fmt.Println(len(b))
			fmt.Println("--------------------")
			fmt.Print("out:")
			fmt.Println(len(out))*/
			b := []byte("heartbeat")
			fmt.Println(b)
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(1000 * time.Second)
		}
	}()


	for {
		if _, err = conn.ReadMessage(); err != nil {
			log.Println(err)
		}
	}

}
