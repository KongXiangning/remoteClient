package wsConn

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"remoteClient/utils"
	"log"
	"os"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	mutex    sync.Mutex
	isClosed bool
}

func Init(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	go conn.readLoop()
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}

	sKey := "dde4b1f8a9e6b814"
	ivParameter := "test121412412312"

	dest,_ := utils.Aes128Decrypt(data,[]byte(sKey),[]byte(ivParameter))
	out,err := coder.GzipDecode(dest)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(out))

	/*out,err := coder.GzipDecode(data)
	if err != nil {
		log.Println(err)
	}
	dest,_ := utils.Aes128Decrypt(out,[]byte(sKey),[]byte(ivParameter))
	fmt.Println(string(dest))*/

	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) Close() {
	conn.wsConn.Close()

	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()

}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			if websocket.IsCloseError(err,1005){
				log.Println(err)
				goto  EXIT
			}
			log.Println(err)
			goto EERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			fmt.Println("closeChan")
			goto EERR
		}
	}

EERR:
	conn.Close()
	return
EXIT:
	conn.Close()
	os.Exit(0)
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
