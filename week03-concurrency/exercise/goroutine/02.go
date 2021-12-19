import (
	"context"
	"fmt"
	"net/http"
)

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello, QCon!")
	})
	return serve("0.0.0.0:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8001", http.DefaultServeMux, stop)
}

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- serveApp(stop)
		// 下面代码可以模拟一个 serve error
		//time.Sleep(3 * time.Second)
		//done <- errors.New("timeoutErr")
	}()
	go func() {
		done <- serveDebug(stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		// 当没有 server 发生错误时这里会被阻塞, 不会执行
		if err := <-done; err != nil {
			fmt.Println("server err", err)
		}
		// 表示某个 server 已经发生异常, 通知其余 server 关闭
		// stopped 这个 flag 表示 channel stop 只关闭一次
		// 重复关闭 channel 会造成 panic
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}

