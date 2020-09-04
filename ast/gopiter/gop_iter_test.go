package gopiter

import (
	"fmt"
	"testing"
)

func Test_Slice(t *testing.T) {
	c := []int{1, 2, 3, 4, 5}
	iter := NewIter(c)
	cnt := 0
	keySum := 0
	valSum := 0
	k := 0
	v := 0
	for iter.Next() {
		Key(iter, &k)
		Value(iter, &v)
		keySum += k
		valSum += v
		cnt++
	}
	fmt.Println(cnt, keySum, valSum)
	if cnt != 5 || keySum != 10 || valSum != 15 {
		t.Fail()
	}
}
