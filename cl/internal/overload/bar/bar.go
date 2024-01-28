package bar

const GopPackage = true

type M = map[string]any

type basetype interface {
	string | int | bool | float64
}

type Var__0[T basetype] struct {
	val T
}

func (p *Var__0[T]) Value() T {
	return p.val
}

type Var__1[T map[string]any] struct {
	val T
}

func (p *Var__1[T]) Value() T {
	return p.val
}

func Gopx_Var_Cast__0[T basetype]() *Var__0[T] {
	return new(Var__0[T])
}

func Gopx_Var_Cast__1[T map[string]any]() *Var__1[T] {
	return new(Var__1[T])
}
