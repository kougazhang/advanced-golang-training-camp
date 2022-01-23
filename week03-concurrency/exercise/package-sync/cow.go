package main

import (
	"sync"
	"sync/atomic"
)

func main() {
	type Map map[string]string
	var m atomic.Value
	m.Store(make(Map))
	var mu sync.Mutex // used only by writes, 只用来写
	// read function can be used to read the data without further synchronization，读函数可以直接读取数据不需要加锁
	read := func(key string) (val string) {
		m1 := m.Load().(Map)
		return m1[key]
	}
	// insert function can be used to update the data without further synchronization
	// 插入函数不需要同步吗？感觉这个注释有问题, 不是 without 而是 with
	insert := func(key, val string) {
		mu.Lock() // synchronize with other potential writers 与其他潜在的 writers 要同步更新
		defer mu.Unlock()
		m1 := m.Load().(Map) // load current value of the data strcture 加载当前的数据结构
		m2 := make(Map)
		// 把所有的老数据拷贝到 m2 中
		for k, v := range m1 {
			m2[k] = v // copy all data from the current object to the new one
		}
		m2[key] = val // do the update that we need, 按需更新
		m.Store(m2)   // atomically replace the current object with the new one
		// At this point all new readers start working with the new version.
		// The old version will be garbage collected once the existing readers
		// (if any) are done with it.
		// 这样就完成了原子化的更新操作。
		// 更新操作完成后，所有新的 readers 会读到新版本的数据
		// 老版本的数据会在（如果存在）当前的 readers 读完后，会被 gc 掉（垃圾回收）
	}
	_, _ = read, insert
}
