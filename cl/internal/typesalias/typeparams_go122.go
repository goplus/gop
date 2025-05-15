//go:build !go1.23
// +build !go1.23

package typesalias

import "go/types"

func TypeArgs(t *Alias) *types.TypeList {
	return nil
}

func TypeParams(t *Alias) *types.TypeParamList {
	return nil
}

func SetTypeParams(t *Alias, params []*types.TypeParam) {
	panic("unsupport Alias.SetTypeParams")
}
