package redis

import (
	"time"

	radix "github.com/mediocregopher/radix/v3"
)

func NewRedisClientPool(address string, size int) (*radix.Pool, error) {
	p, err := radix.NewPool("tcp", address, size, radix.PoolConnFunc(customConnFunc))

	if err != nil {
		return nil, err
	}

	return p, nil
}

func customConnFunc(network, addr string) (radix.Conn, error) {
	return radix.Dial(network, addr,
		radix.DialTimeout(5*time.Second),
	)
}

func Proccess(p *radix.Pool, f func(rc radix.Client) error) (e error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		e = util.PanicToErr(err)
	// 	}
	// }()

	return f(p)
}
