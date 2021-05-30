// Package user provide Go+ "os/user" package, as "os/user" package in Go.
package user

import (
	user "os/user"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execCurrent(_ int, p *gop.Context) {
	ret0, ret1 := user.Current()
	p.Ret(0, ret0, ret1)
}

func execLookup(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := user.Lookup(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupGroup(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := user.LookupGroup(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupGroupId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := user.LookupGroupId(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := user.LookupId(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmUnknownGroupErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(user.UnknownGroupError).Error()
	p.Ret(1, ret0)
}

func execmUnknownGroupIdErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(user.UnknownGroupIdError).Error()
	p.Ret(1, ret0)
}

func execmUnknownUserErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(user.UnknownUserError).Error()
	p.Ret(1, ret0)
}

func execmUnknownUserIdErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(user.UnknownUserIdError).Error()
	p.Ret(1, ret0)
}

func execmUserGroupIds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*user.User).GroupIds()
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("os/user")

func init() {
	I.RegisterFuncs(
		I.Func("Current", user.Current, execCurrent),
		I.Func("Lookup", user.Lookup, execLookup),
		I.Func("LookupGroup", user.LookupGroup, execLookupGroup),
		I.Func("LookupGroupId", user.LookupGroupId, execLookupGroupId),
		I.Func("LookupId", user.LookupId, execLookupId),
		I.Func("(UnknownGroupError).Error", (user.UnknownGroupError).Error, execmUnknownGroupErrorError),
		I.Func("(UnknownGroupIdError).Error", (user.UnknownGroupIdError).Error, execmUnknownGroupIdErrorError),
		I.Func("(UnknownUserError).Error", (user.UnknownUserError).Error, execmUnknownUserErrorError),
		I.Func("(UnknownUserIdError).Error", (user.UnknownUserIdError).Error, execmUnknownUserIdErrorError),
		I.Func("(*User).GroupIds", (*user.User).GroupIds, execmUserGroupIds),
	)
	I.RegisterTypes(
		I.Type("Group", reflect.TypeOf((*user.Group)(nil)).Elem()),
		I.Type("UnknownGroupError", reflect.TypeOf((*user.UnknownGroupError)(nil)).Elem()),
		I.Type("UnknownGroupIdError", reflect.TypeOf((*user.UnknownGroupIdError)(nil)).Elem()),
		I.Type("UnknownUserError", reflect.TypeOf((*user.UnknownUserError)(nil)).Elem()),
		I.Type("UnknownUserIdError", reflect.TypeOf((*user.UnknownUserIdError)(nil)).Elem()),
		I.Type("User", reflect.TypeOf((*user.User)(nil)).Elem()),
	)
}
