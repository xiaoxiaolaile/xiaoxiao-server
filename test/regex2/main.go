package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	str := "  ? 红烧肉怎么做"
	str = strings.TrimSpace(str)
	str = regexp.MustCompile(`^\?`).ReplaceAllString(str, `(自定义开始) ?`)
	fmt.Println(str)
}
