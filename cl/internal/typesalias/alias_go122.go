//go:build go1.22
// +build go1.22

package typesalias

import "go/types"

const Support = true

type Alias = types.Alias

func NewAlias(obj *types.TypeName, rhs types.Type) *Alias {
	return types.NewAlias(obj, rhs)
}

func Unalias(t types.Type) types.Type {
	return types.Unalias(t)
}
