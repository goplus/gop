package main

import (
	"fmt"
	"reflect"
)

type basetype interface {
	int | string
}

type Table struct {
}

func Gopt_Table_Gopx_Col__0[T basetype](p *Table, name string) {
	fmt.Printf("Gopt_Table_Gopx_Col__0 %v: %s\n", reflect.TypeOf((*T)(nil)).Elem(), name)
}

func Gopt_Table_Gopx_Col__1[Array any](p *Table, name string) {
	fmt.Printf("Gopt_Table_Gopx_Col__1 %v: %s\n", reflect.TypeOf((*Array)(nil)).Elem(), name)
}
