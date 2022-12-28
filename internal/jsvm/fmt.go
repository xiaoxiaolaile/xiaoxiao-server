package jsvm

import "fmt"

type Fmt struct {
}

func (sender *Fmt) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func (sender *Fmt) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (sender *Fmt) Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func (sender *Fmt) Print(a ...interface{}) (int, error) {
	return fmt.Print(a...)
}
