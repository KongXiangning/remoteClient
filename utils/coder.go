package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipEconder(in []byte)(out []byte,err error)  {
	var (
		buffer bytes.Buffer
	)

	writer := gzip.NewWriter(&buffer)
	if _,err = writer.Write(in);err != nil{
		writer.Close()
		return
	}

	if err = writer.Close();err != nil{
		return
	}

	return buffer.Bytes(),nil

}

func GzipDecode(in []byte)(out []byte,err error)  {
	reader ,err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
