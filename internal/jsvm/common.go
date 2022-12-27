package jsvm

import (
	"fmt"
	"strconv"
)

var Int = func(s interface{}) int {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return i
}
