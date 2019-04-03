package wsConn

import "fmt"

func (conn *Connection) textReadHandler() {
	var (
		data []byte
	)
	data = <-conn.inTextChan
	fmt.Println(string(data))
}
