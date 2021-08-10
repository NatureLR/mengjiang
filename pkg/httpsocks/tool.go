package httpsocks

import (
	"fmt"
	"net"
	"reflect"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func connPrint(c net.Conn) {
	if c == nil {
		fmt.Println("C为空")
		return
	}
	buf := make([]byte, 1024)
	n, err := c.Read(buf[:])
	checkErr(err)
	fmt.Println("===========connPrint start===========")
	fmt.Println("connPrint", string(buf[:n]))
	fmt.Println("===========connPrint end===========")
}

func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
