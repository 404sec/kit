package typex

import (
	pinyin "github.com/mozillazg/go-pinyin"
	"regexp"
	"unicode"

	//"log"
	"strings"
)

type Pin string

func (p Pin) QuanPin() string {
	if len(string(p)) > 0 {
		a := pinyin.NewArgs()

		pins := pinyin.LazyPinyin(string(p), a)
		rePin := strings.Join(pins, "")

		return (rePin)
	} else {
		return ""
	}

}

func (p Pin) JianPin(han string) string {

	if len(han) > 0 {
		a := pinyin.NewArgs()
		//a.Heteronym = true
		a.Style = pinyin.FirstLetter
		pins := pinyin.LazyPinyin(han, a)
		rePin := strings.Join(pins, "")

		return (rePin)
	}
	return ""

}
func (p *Pin) PinAll() string {

	return p.QuanPin() + "," + p.JianPin(string(*p))
}

func (p *Pin) FristPin() string {
	han := string(*p)
	chiReg := regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z]+")
	str := chiReg.ReplaceAllString(han, "")

	if unicode.Is(unicode.Han, []rune(str)[0]) {
		return strings.ToUpper(string(p.JianPin(string(str))[:1]))

	}
	return strings.ToUpper(string(str[:1]))

}

/*
func (p *Pin) UnmarshalJSON(data []byte) (err error) {
	*p = Pin(string(data))

	return nil
}

func (p Pin) MarshalJSON() ([]byte, error) {
	if len(string(p)) > 0 {
		p = Pin(string(p)[0])
	}
	return []byte(p.JianPin()), nil
}
*/
