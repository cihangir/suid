package suid

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestSuid(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	suid := NewSUID(12)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for i := 0; i < 1025; i++ {
			suid.Generate()
			fmt.Println(suid.Generate())
		}
	}()
	// go func() {
	// 	wg.Add(1)
	// 	for i := 0; i < 2024; i++ {
	// 		suid.Generate()

	// 		// fmt.Println(suid.Generate())
	// 	}
	// }()
	wg.Wait()
}
