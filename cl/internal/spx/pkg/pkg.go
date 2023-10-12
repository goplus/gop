package pkg

var (
	GopPackage = true
)

type Vector struct {
	x int
	y int
}

func NewVector(x, y int) *Vector {
	return &Vector{x, y}
}

func (v *Vector) Add__0(x int, y int) {
	v.x += x
	v.y += y
}

func (v *Vector) Add__1(o *Vector) {
	v.Add__0(o.x, o.y)
}
