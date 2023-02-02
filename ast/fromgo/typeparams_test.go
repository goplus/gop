//go:build go1.18
// +build go1.18

package fromgo

import "testing"

func TestIndexListExpr(t *testing.T) {
	test(t, `package main

type N1 = T1[int]

type N2 = T2[string, int]

var f = test[string, int](1, 2)
`)
}
