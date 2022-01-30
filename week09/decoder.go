package week_9

import (
	"bufio"
)

type Decoder interface {
	WithStream(stream *bufio.Reader) Decoder
	Decode() Decoder
	Result() (TcpObj, error)
}

type TcpObj interface {
	Pretty() (string, error)
	Body() (string, error)
}
