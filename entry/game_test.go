package entry

import (
	"fmt"
	"testing"

	_ "golang.org/x/image/bmp"
)

func TestGame_CreateBullet(t *testing.T) {
	num := 2
	x := 10
	num1 := 0
	var xList = make([]int, num)
	for i := 0; i < num; i++ {
		if i%2 == 0 {
			xList[i] = x + num1
			num1 = x + num1
			continue
		}
		xList[i] = x - num1
	}
	fmt.Println(xList)
}
