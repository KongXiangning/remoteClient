package wsConn

import (
	"fmt"
	"kube-client/execute"
	"log"
	"remoteClient/utils"
)

func (conn *Connection) binaryReadHandler() {
	var (
		data []byte
	)
	data = <-conn.inBinaryChan
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
		msgData = data[82:]
		msgData, err = utils.Aes128Decrypt(msgData, utils.SKey, ivs)
		if err != nil {
			//conn.WriteMessage()
			log.Println(err)
		}
	case 2:
		fallthrough
	case 3:
		ivs := data[2:readIndex]
		msgData = data[82:]
		msgData, err = utils.Aes128Decrypt(msgData, utils.SKey, ivs)
		if err != nil {
			//conn.WriteMessage()
			log.Println(err)
		}
		msgData, err = utils.GzipDecode(msgData)
	}

	method = string(data[28 : 28+data[27]])
	fmt.Println(method)
	fmt.Println(string(msgData))
	fmt.Println(string(utils.SKey))

	kubeTransfer := execute.KubeTransfer{Types: data[81], Method: method, Hid: data[19:27], HandleJson: msgData}
	client := execute.GetClient()
	client.Execute(kubeTransfer, conn.transferChan)
}
