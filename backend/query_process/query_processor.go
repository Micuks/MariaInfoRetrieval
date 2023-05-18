package query_process

import (
	"fmt"
	"runtime/debug"

	"github.com/yanyiwu/gojieba"
)

func ProcessQuery(query string) []string {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
		}
	}()

	seg := gojieba.NewJieba()
	defer seg.Free()

	words := seg.Cut(query, true) // Use Jeba for Chinese text segment

	return words
}
