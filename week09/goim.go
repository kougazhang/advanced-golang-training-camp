package week_9

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	MaxBodySize = int32(1 << 12)
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

type TcpObj4GoIm struct {
	ProtocolVersion int32
	Operation       int32
	Seq             int32
	Content         string
}

func (obj *TcpObj4GoIm) Body() (string, error) {
	return obj.Content, nil
}

func (obj *TcpObj4GoIm) Pretty() (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

var _ TcpObj = (*TcpObj4GoIm)(nil)

type GoImDecoder struct {
	buf []byte
	rcv *bufio.Reader
	err error
	dst *TcpObj4GoIm
}

func (decoder *GoImDecoder) reset() {
	decoder.err = nil
	decoder.dst = nil
	decoder.buf = nil
	decoder.rcv = nil
}

func (decoder *GoImDecoder) WithStream(stream *bufio.Reader) Decoder {
	decoder.reset()
	decoder.rcv = stream
	return decoder
}

func (decoder *GoImDecoder) Decode() Decoder {
	decoder.buf, decoder.err = decoder.rcv.Peek(_rawHeaderSize)
	if decoder.err != nil {
		return decoder
	}
	pkgLen := decoder.readInt32(_packOffset, _packSize)
	headerLen := decoder.readInt32(_headerOffset, _headerSize)

	if headerLen != _rawHeaderSize {
		decoder.err = errors.New("header length error")
		return decoder
	}

	decoder.dst = &TcpObj4GoIm{
		ProtocolVersion: decoder.readInt32(_verOffset, _verSize),
		Operation:       decoder.readInt32(_opOffset, _opSize),
		Seq:             decoder.readInt32(_seqOffset, _seqSize),
	}

	if bodyLen := pkgLen - headerLen; bodyLen > 0 {
		decoder.buf, decoder.err = decoder.rcv.Peek(int(pkgLen))
		if decoder.err != nil {
			return decoder
		}
		decoder.dst.Content = decoder.readString(_rawHeaderSize, int(bodyLen))
	}
	return decoder
}

func (decoder *GoImDecoder) readInt32(start int, step int) int32 {
	if step <= 0 {
		step = 1
	}
	bytes := decoder.buf[start : start+step]
	var tmp = int32(bytes[step-1])
	if step == 1 {
		return tmp
	}
	for i := 1; i <= step-1; i++ {
		tmp = tmp | int32(bytes[step-1-i])<<i*8
	}
	return tmp
}

func (decoder *GoImDecoder) readString(start int, step int) string {
	return string(decoder.buf[start : start+step])
}

func (decoder *GoImDecoder) Result() (TcpObj, error) {
	return decoder.dst, decoder.err
}

var _ Decoder = (*GoImDecoder)(nil)
var GoImD = &GoImDecoder{}
