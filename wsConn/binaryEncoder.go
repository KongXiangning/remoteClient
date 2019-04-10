package wsConn

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
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
		err               error
		compData          []byte
		encryptResultData []byte
		encryptJsonData   []byte
		methodData        []byte
		resultGzip        byte
		resultLen         int32
		jsonLen           int32
		mlen              int
	)

	data := new(bytes.Buffer)
	ivs := make([]byte, 16)
	if _, err = io.ReadFull(rand.Reader, ivs); err != nil {
		log.Print(err)
		goto ERR
	}

	if len(transfer.Result) > 500 {
		resultGzip = 1
		if compData, err = utils.GzipEconder([]byte(transfer.Result)); err != nil {
			log.Println(err)
			goto ERR
		}
	} else {
		resultGzip = 0
		compData = []byte(transfer.Result)
	}
	if encryptResultData, err = utils.Aes128Encrypt(compData, utils.SKey, ivs); err != nil {
		log.Println(err)
		goto ERR
	}
	resultLen = int32(len(encryptResultData))

	compData = nil
	if compData, err = utils.GzipEconder(transfer.HandleJson); err != nil {
		log.Println(err)
		goto ERR
	}
	if encryptJsonData, err = utils.Aes128Encrypt(compData, utils.SKey, ivs); err != nil {
		log.Println(err)
		goto ERR
	}
	jsonLen = int32(len(encryptJsonData))

	methodData = []byte(transfer.Method)
	mlen = len(methodData)
	if mlen > 255 {
		log.Println("method is to long")
		goto ERR
	}

	data.WriteByte(1)
	data.WriteByte(16)
	data.Write(ivs)
	data.WriteByte(transfer.Types)
	data.WriteByte(byte(mlen))
	data.Write(methodData)
	data.WriteByte(resultGzip)
	data.Write(Int32ToBytes(resultLen))
	data.Write(encryptResultData)
	data.Write(Int32ToBytes(jsonLen))
	data.Write(encryptJsonData)

	conn.outBinaryChan <- data.Bytes()
	return
ERR:
	data.WriteByte(0)
	data.WriteString(err.Error())
	conn.outBinaryChan <- data.Bytes()
}

func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}
