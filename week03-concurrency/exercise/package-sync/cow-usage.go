package main

func main() {
	// holds current server configuration 动态更新当前服务的配置
	var config atomic.Value

	go func() {
		// Reload config every 10 seconds 每 10 秒重载配置
		// and update config value with the new version 用新版本的值更新配置项
		for {
			time.Sleep(10 * time.Second)
			config.Store(loadConfig())
		}
	}()

	// Create worker goroutines that handle incoming requests
	// using the latest config value. 处理新请求时创建的 goroutine worker 使用最新的配置
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				c := config.Load()
				// Handle request r using config c.
				_, _ = r.c
			}
		}()
	}
}
