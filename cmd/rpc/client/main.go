package main

import (
	"fmt"
	"time"
)

func main() {
	//	r,err:=tcp()
	r, err := http()
	if err != nil {
		panic(err)
	}
	for {
		//<-w
		var s string
		req, resp, err := r.Request("test", nil, &s)
		if err != nil {
			fmt.Println(err)

		} else {
			fmt.Println(req, resp, err, s)
		}
		time.Sleep(time.Second)
	}

}
