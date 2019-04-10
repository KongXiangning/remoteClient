package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"kube-client/execute"
	"log"
	"net/http"
	"net/url"
	"remoteClient/utils"
	"remoteClient/wsConn"
	"time"
)

//var origin = "http://localhost:7777"
//var url = "ws://localhost:7777/ws"
var addr = flag.String("addr", "192.168.5.77:7777", "http service address")

func main() {
	var (
		wsClient      *websocket.Conn
		conn          *wsConn.Connection
		err           error
		u             url.URL
		requestHeader http.Header
	)
	execute.GetClient()

	flag.Parse()
	log.SetFlags(0)
	u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	requestHeader = http.Header{}

	password := "suneee"
	skeyByte, err := ioutil.ReadFile("e:\\key.txt")
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, 16)
	if len(skeyByte) != 0 {
		utils.SKey = skeyByte
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}
		if b, err := utils.Aes128Encrypt([]byte(password), utils.SKey, iv); err != nil {
			panic(err)
		} else {
			password = base64.RawStdEncoding.EncodeToString(b)
		}

	}
	requestHeader.Add("kid", "1010")
	requestHeader.Add("login", password)
	requestHeader.Add("key", base64.RawStdEncoding.EncodeToString(iv))

ERR:
	if wsClient, _, err = websocket.DefaultDialer.Dial(u.String(), requestHeader); err != nil {
		log.Println(err)
		time.Sleep(10 * time.Second)
		goto ERR
	}

	if conn, err = wsConn.Init(wsClient); err != nil {
		log.Println(err)
		wsClient.Close()
		goto ERR
	}

	wsClient.SetPongHandler(func(message string) error {
		var (
			data      []byte
			ivs       []byte
			validText []byte
			echo      string
			err       error
		)
		data = []byte(message)
		ivs = data[2:18]
		switch data[0] {
		case 0:
			index := 19 + data[18]
			origKey := data[19:index]
			utils.SKey, err = utils.Aes128Decrypt(origKey, utils.SKey, ivs)
			if err != nil {
				panic(err)
			}
			validText = data[index+1 : index+1+data[index]]
		case 1:
			validText = data[19 : 19+data[18]]
		}
		if becho, err := utils.Aes128Decrypt(validText, utils.SKey, ivs); err != nil {
			panic(err)
		} else {
			echo = string(becho)
		}

		if echo == "suneee" {

		} else {
			conn.IsExist <- 0
		}
		return nil
	})

	log.Println("connection.go success!")

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
			time.Sleep(1 * time.Second)
			b := []byte("heartbeat")
			fmt.Println(b)
			if err = conn.WriteTextMessage([]byte("heartbeat")); err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()

	/*for {
		conn.BinaryDecoder()
	}
	for {
		if _, err = conn.ReadMessage(); err != nil {
			log.Println(err)
		}
	}*/

	select {
	case flag := <-conn.IsExist:
		if flag == 1 {
			fmt.Println("read is error")
			goto ERR
		} else if flag == 2 {
			fmt.Println("write is error")
			goto ERR
		}
	}

}
