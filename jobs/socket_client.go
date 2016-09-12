package jobs

import (
	"bufio"
	"net"
	"strconv"
)

type Socket interface {
	Send([]byte) error
	Read() ([]byte, error)
	Close() error
}

type Data struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time string      `json:"time"`
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SocketClient struct {
	Conn net.Conn
}

func NewSocketClient(address string, port int) (Socket, error) {
	service := address + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", service)

	if err != nil {
		return nil, err
	}
	return &SocketClient{conn}, nil
}

func (sc *SocketClient) Send(data []byte) error {
	_, e := sc.Conn.Write(data)
	return e
}

func (sc *SocketClient) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	r := bufio.NewReader(sc.Conn)

	n, err := r.Read(buf)
	return buf[:n], err
}

func (sc *SocketClient) Close() error {
	return sc.Conn.Close()
}
