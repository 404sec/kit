package typex

import (
	"fmt"
	//"strconv"
)

type Mask string

//接受数据进行解密 数据转为加密字符串
func (n Mask) MarshalJSON() ([]byte, error) {

	//加密
	return []byte(fmt.Sprintf("%.2f", n)), nil
}

//接受数据进行加密 数据放到结构体
func (n *Mask) UnmarshalJSON(b []byte) error {

	//解密
	//d, err := strconv.ParseFloat(string(b), 64)
	//*n = Mask(d)
	return nil
	//n=fmt.Sprintf("%.2f", b))
}
