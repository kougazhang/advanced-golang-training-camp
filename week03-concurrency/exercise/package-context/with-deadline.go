package main

import (
	"context"
	"fmt"
	"time"
)

const shortDuration = time.Millisecond

func main() {
	d := time.Now().Add(shortDuration)
	ctx, _ := context.WithDeadline(context.Background(), d)

	// 即便 ctx 会过期，为了以防万一调用 cancellation 函数也是一个好习惯
	// 失败可能会导致 ctx 不过期，这样的话 parent 会比预期时存活更长
	// defer cancel()

	select {
	case <-time.After(time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
