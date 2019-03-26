package wsConn

import (
	"fmt"
	"log"
	"remoteClient/utils"
)

func (conn *Connection) binaryReadHandler() {
	var (
		data []byte
	)
	data = <-conn.binaryChan
	go conn.binaryDecoder(data)

}

func (conn *Connection) binaryDecoder(data []byte) {
	var (
		msgData []byte
		method  string
		err     error
	)
	readIndex := 2 + data[1]
	isCC := data[0] | data[readIndex]
	switch isCC {
	case 0:
		fallthrough
	case 1:
		ivs := data[2:readIndex]
		msgData = data[60:]
		msgData, err = utils.Aes128Decrypt(msgData, utils.SKey, ivs)
		if err != nil {
			//conn.WriteMessage()
			log.Println(err)
		}
	case 2:
		fallthrough
	case 3:
		ivs := data[2:readIndex]
		msgData = data[60:]
		msgData, err = utils.Aes128Decrypt(msgData, utils.SKey, ivs)
		if err != nil {
			//conn.WriteMessage()
			log.Println(err)
		}
		msgData, err = utils.GzipDecode(msgData)
	}

	method = string(data[readIndex+2 : readIndex+2+data[readIndex+1]])
	fmt.Println(method)
	fmt.Println(string(msgData))
	fmt.Println(string(utils.SKey))
}
