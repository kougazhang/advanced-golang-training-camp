package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Pipeline 示范(demonstrate) 使用 Group 实现多阶段（multi-stage）pipeline:

func main() {
	m, err := MD5All(context.Background(), ".")
	if err != nil {
		log.Fatal(err)
	}

	for k, sum := range m {
		fmt.Printf("%s:\t%x\n", k, sum)
	}
}

type result struct {
	path string
	sum  [md5.Size]byte
}

// MD5All 读取在根目录下所有的文件，
// 返回一个map key 是文件路径 value 是文件内容的 md5。
// 如果目录遍历失败或任何读取操作失败，MD5All 返回 error。
func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	// 当 g.Wait() 返回时 ctx 会被取消。
	// 当这个版本的 MD5All 返回 error 时，我们知道所有的 goroutines 已经完成
	// 并且他们使用的内存可以垃圾回收了。

	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string)

	g.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result)
	const numDigesters = 20
	for i := 0; i < numDigesters; i++ {
		g.Go(func() error {
			for path := range paths {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				select {
				case c <- result{path, md5.Sum(data)}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}
	go func() {
		g.Wait()
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		m[r.path] = r.sum
	}
	// 检查是否有 goroutine 失败。由于 g 能对 error 进行累加，
	// 我们不需要对单个的 goroutine 进行发送或检查。
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return m, nil
}
