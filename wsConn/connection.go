package wsConn

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Connection struct {
	wsConn        *websocket.Conn
	binaryChan    chan []byte
	textChan      chan []byte
	outBinaryChan chan []byte
	outTextChan   chan []byte
	closeChan     chan byte
	IsExist       chan int

	mutex    sync.Mutex
	isClosed bool
}

func Init(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:        wsConn,
		binaryChan:    make(chan []byte, 1000),
		textChan:      make(chan []byte, 1000),
		outBinaryChan: make(chan []byte, 1000),
		outTextChan:   make(chan []byte, 1000),
		closeChan:     make(chan byte, 1),
		IsExist:       make(chan int),
	}

	go conn.readLoop()
	go conn.writeLoop()
	go func() {
		for {
			conn.binaryReadHandler()
		}
	}()

	return
}

func (conn *Connection) WriteBinaryMessage(data []byte) (err error) {
	select {
	case conn.outBinaryChan <- data:
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
		index int
		data  []byte
		err   error
	)
	for {
		if index, data, err = conn.wsConn.ReadMessage(); err != nil {
			if websocket.IsCloseError(err, 1005) {
				log.Println(err)
				goto EXIT
			}
			goto EERR
		}
		switch index {
		case 1:
			conn.textChan <- data
		case 2:
			conn.binaryChan <- data
		}
	}

EERR:
	conn.Close()
	conn.IsExist <- 1
	return
EXIT:
	conn.Close()
	conn.IsExist <- 0
	//os.Exit(0)
}

func (conn *Connection) writeLoop() {
	var (
		data  []byte
		wtype int
		err   error
	)

	for {
		select {
		case data = <-conn.outBinaryChan:
			wtype = websocket.BinaryMessage
		case data = <-conn.outTextChan:
			wtype = websocket.TextMessage
		}
		if err = conn.wsConn.WriteMessage(wtype, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
	conn.IsExist <- 2
}
