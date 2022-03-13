package redis

import (
	"sync"

	radix "github.com/mediocregopher/radix/v3"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
)

type RadixRC3 struct {
	p       *radix.Pool
	l       sync.Mutex
	address string
	size    int
}

func NewRedisClient(conf *conf.Cache) (*RadixRC3, error) {
	rc := &RadixRC3{
		address: conf.Connect.Source,
		size:    int(conf.Connect.Size),
	}
	if err := rc.GetRedisClient(); err != nil {
		return rc, err
	}
	return rc, nil
}

func (r *RadixRC3) GetRedisClient() error {
	if r.p == nil {
		r.l.Lock()
		defer r.l.Unlock()
		if r.p == nil {
			var err error
			r.p, err = NewRedisClientPool(r.address, r.size)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RadixRC3) Get(key string) (value string, err error) {
	return key, nil
}

func (r *RadixRC3) Set(key string, value string, ttl int) error {
	return nil
}

func (r *RadixRC3) HashGet(key string, fields *[]string) (values *map[string]string, err error) {
	return &map[string]string{}, nil
}

func (r *RadixRC3) HashSet(key string, hash *map[string]string) error {
	return nil
}
