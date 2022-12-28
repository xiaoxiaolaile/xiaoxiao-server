package core

import "time"

type Time struct {
}

func (t *Time) Now() time.Time {
	return time.Now()
}
func (t *Time) Sleep(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}
func (t *Time) Unix(usec int64) time.Time {
	return time.UnixMicro(usec)
}
func (t *Time) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec)
}

func (t *Time) Parse(value, layout, name string) time.Time {
	if len(name) == 0 {
		name = "Asia/Shanghai"
	}
	loc, _ := time.LoadLocation(name)
	_t, _ := time.ParseInLocation(layout, value, loc)
	return _t
}
