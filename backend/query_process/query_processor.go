package query_process

import (
	"github.com/yanyiwu/gojieba"
)

func ProcessQuery(query string) []string {
	seg := gojieba.NewJieba()
	defer seg.Free()

	words := seg.Cut(query, true) // Use Jeba for Chinese text segment

	return words
}
