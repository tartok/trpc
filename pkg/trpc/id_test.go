package trpc

import (
	"fmt"
	"testing"
)

func TestId(t *testing.T) {

	for i := 0; i < 10000000; i++ {
		s := GetId()
		_ = s
		fmt.Println(s)
		//fmt.Println(v, time.Now().UnixMilli())
		//if v == p {
		//	t.Fatal("time is not increasing", v-p)
		//}

		//bson.NewObjectID()
	}
}
