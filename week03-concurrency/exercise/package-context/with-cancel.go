package main

import (
	"context"
	"fmt"
)

// WithCancel 返回携带新的 Done channel 的 parent context 的拷贝.
// 在返回函数被调用或当 parent context Done channel 被关闭时, 返回的 context Done channel 会被关闭

// 取消 context 会释放与之相关的资源, 所以在运行 context 完成时尽可能早地调用 `WithCancel`

// 这个 demo 展示了使用可取消的 context (a cancelable context) 预防 goroutine 泄露.
// 在函数结尾, gen 中已启动的 goroutine 会没有泄露的返回.

func main() {
	// gen 在独立的 goroutine 中产生整数并且把该整数发送到 channel 中.
	// 一旦产生的整数被消费完, gen 的调用者需要取消 context 以确保 gen 启动的 goroutine 不会出现泄露
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // 此处 return, 该 goroutine 不会导致内存泄露
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}

	ins := NewImplA(ImplB{})
	ins.A()
	ins.B()
	ins.C()
	ins.D()
	// 输出如下. 这表明方法调用时会调用最浅层的方法.
	//i am implA.A
	//i am implA.B
	//i am implA.C
	//i am ImplB.D

	var (
		a  ImplA
		iA InterfaceA
	)
	// 这一步能够通过编译,说明编译器认为 ImplA 实现了 InterfaceA 的方法
	iA = a
	iA.D()
	// 会抛出 runtime error, 因为 ImplA 没有实现
	// panic: runtime error: invalid memory address or nil pointer dereference
}

type InterfaceA interface {
	A()
	B()
	C()
	D()
}

type ImplA struct {
	InterfaceA
}

func (ImplA) A() {
	fmt.Println("i am implA.A")
}

func (ImplA) B() {
	fmt.Println("i am implA.B")
}

func (ImplA) C() {
	fmt.Println("i am implA.C")
}

type ImplB struct {
}

func (ImplB) A() {
	fmt.Println("i am ImplB.A")
}

func (ImplB) B() {
	fmt.Println("i am ImplB.B")
}

func (ImplB) C() {
	fmt.Println("i am ImplB.C")
}

func (ImplB) D() {
	fmt.Println("i am ImplB.D")
}

func NewImplA(a InterfaceA) ImplA {
	return ImplA{InterfaceA: a}
}
