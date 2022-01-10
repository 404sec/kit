package typex

import (
	"time"
)

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

//序列化
func (t Time) MarshalJSON() ([]byte, error) {

	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	if t.Year() > 1000 {
		b = time.Time(t).AppendFormat(b, timeFormart)
		b = append(b, '"')
		return b, nil
	} else {
		b = append(b, '"')
		return b, nil
	}
}

//反序列化
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	newstring := string(data)
	var err error
	tt, err2 := time.ParseInLocation(timeFormart, newstring, time.Local)
	if err != nil {
		return err2
	} else {
		*t = Time(tt)
	}
	return err
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}

func (t Time) Year() int {
	return time.Time(t).Year()
}

func (t Time) Now() Time {
	tt := Time(time.Now())
	return tt
}
