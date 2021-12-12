package main

import (
	"errors"
	"fmt"
	"log"
)

// 在 Dao 层遇到 sql.ErrNoRows 的时候，应该 Wrap 这个 error 抛给上层。
// 但是 sql.ErrNoRows 是一个底层 error ，不能直接暴露给业务层，Dao 层应该把它封装为 Dao 层自定义的 error 类型。

// main 方法内为业务层代码
func main() {
	_, err := QueryByName("foo")
	// 只关心是否有 error, 处理如下：
	if err != nil {
		log.Printf("handleErr: query %s err %v\n", "foo", err)
	}
	// 关心 error 的 root error, 处理如下：
	if errors.Is(err, &ErrNoRows{}) {
		log.Printf("errors.Is: %v\n", err)
	}
	// 对 error 进行断言
	var userErr *UserNotFoundErr
	if errors.As(err, &userErr) {
		log.Printf("errors.As: name %s, root error %v \n", userErr.Name, userErr.Err)
	}
	// 如果当前的业务处理不了, 想继续上抛, 处理如下：
	err = fmt.Errorf("query: err %w", err)
	log.Printf("WrapErr: %v", err)
}

// Dao 层 User 用户模型
type User struct {
	ID   string
	Name string
	Age  int
}

func QueryByName(name string) (*User, error) {
	return nil, &UserNotFoundErr{
		Name: name,
		Err:  &ErrNoRows{},
	}
}

// Dao 层对外暴露的 error
type UserNotFoundErr struct {
	Name string
	Err  error
}

func (u *UserNotFoundErr) Error() string {
	return u.Name + ": not found"
}

func (u *UserNotFoundErr) Unwrap() error {
	return u.Err
}

// 定义 ErrNoRows 指代 sql.ErrNoRows
type ErrNoRows struct{}

func (e *ErrNoRows) Error() string {
	return "no rows"
}
