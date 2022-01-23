package main

import "fmt"
import "strings"
import "testing"

func join() string {
	return strings.Join([]string{"a", "b", "c", "d", "e", "f"}, "|")
}

func format() string {
	return fmt.Sprintf("%s%s%s%s%s%s", "a", "b", "c", "d", "e", "f")
}

func BenchmarkJoinString(b *testing.B) {
	join()
}

func BenchmarkFormatString(b *testing.B) {
	format()
}
