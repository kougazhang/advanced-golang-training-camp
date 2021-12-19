package homework

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// 测试: 一个退出，全部注销退出
func Test1(t *testing.T) {
	// Create workgroup
	var wg Group
	// Add function to cancel execution using os signal
	wg.Add(Signal())
	// Create http server
	srv := http.Server{Addr: "127.0.0.1:8080"}
	// Add function to start and stop http server
	wg.Add(
		Server(
			func() error {
				fmt.Printf("Server listen at %v\n", srv.Addr)
				err := srv.ListenAndServe()
				fmt.Printf("Server stopped listening with error: %v\n", err)
				if err != http.ErrServerClosed {
					return err
				}
				return nil
			},
			func() error {
				fmt.Println("Server is about to shutdown")
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
				defer cancel()

				err := srv.Shutdown(ctx)
				fmt.Printf("Server shutdown with error: %v\n", err)
				return err
			},
		))
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context canceled")
		cancel()
	}()
	// Add function to cancel execution using context
	wg.Add(Context(ctx))
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup run stopped with error: %v\n", err)
}

// 测试: 支持接收 signal 退出
func Test2(t *testing.T) {
	// Create workgroup
	var wg Group
	// Add function to cancel execution using os signal
	wg.Add(Signal())
	// Create http server
	srv := http.Server{Addr: "127.0.0.1:8080"}
	// Add function to start and stop http server
	wg.Add(
		Server(
			func() error {
				fmt.Printf("Server listen at %v\n", srv.Addr)
				err := srv.ListenAndServe()
				fmt.Printf("Server stopped listening with error: %v\n", err)
				if err != http.ErrServerClosed {
					return err
				}
				return nil
			},
			func() error {
				fmt.Println("Server is about to shutdown")
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
				defer cancel()

				err := srv.Shutdown(ctx)
				fmt.Printf("Server shutdown with error: %v\n", err)
				return err
			},
		))
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context canceled")
		cancel()
	}()
	// Add function to cancel execution using context
	wg.Add(Context(ctx))
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup run stopped with error: %v\n", err)
}
