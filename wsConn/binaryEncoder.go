package wsConn

import (
	"bytes"
	"crypto/rand"
	"io"
	"kube-client/execute"
	"log"
	"remoteClient/utils"
)

func (conn *Connection) binaryWriteHandler() {
	var (
		transfer execute.KubeTransfer
	)
	transfer = <-conn.transferChan
	go conn.binaryEncoder(transfer)
}

func (conn *Connection) binaryEncoder(transfer execute.KubeTransfer) {
	var (
		err         error
		compData    []byte
		encryptData []byte
		methodData  []byte
	)
	ivs := make([]byte, 16)
	if _, err = io.ReadFull(rand.Reader, ivs); err != nil {
		log.Print(err)
	}

	if compData, err = utils.GzipEconder([]byte(transfer.Result)); err != nil {
		log.Println(err)
	}

	if encryptData, err = utils.Aes128Encrypt(compData, utils.SKey, ivs); err != nil {
		log.Println(err)
	}

	data := new(bytes.Buffer)
	methodData = []byte(transfer.Method)
	mlen := len(methodData)
	if mlen > 255 {
		log.Println("err")
	}

	data.WriteByte(1)
	data.WriteByte(16)
	data.Write(ivs)
	data.WriteByte(transfer.Types)
	data.WriteByte(byte(mlen))
	data.Write(methodData)
	data.Write(encryptData)

	conn.outBinaryChan <- data.Bytes()
}
