package homework

import (
	"fmt"
	"net/http"
	"testing"
)

// 测试收到退出信号时退出
func TestGroupServer_Run(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", fn)
	srv1 := &Server{
		Name: "srv1",
		Server: &http.Server{
			Addr:    ":8000",
			Handler: mux,
		},
		CertFile: "",
		KeyFile:  "",
	}
	srv2 := &Server{
		Name: "srv2",
		Server: &http.Server{
			Addr:    ":8001",
			Handler: mux,
		},
		CertFile: "",
		KeyFile:  "",
	}
	var group GroupServer
	group.Register(srv1, srv2)
	err := group.Run()
	if err != nil {
		fmt.Println("err", err)
	}
}
