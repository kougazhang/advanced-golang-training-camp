package main

import "testing"
import "time"

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

func BenchmarkFibPrepare(b *testing.B) {
	time.Sleep(time.Second * 3)
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

func BenchmarkFibResetPrepare(b *testing.B) {
	time.Sleep(time.Second * 3)
	b.ResetTimer() // 重置定时器
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}
