package exec

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"

	"github.com/qiniu/qlang/goprj"
)

// -----------------------------------------------------------------------------

// TypeInferrer represents a TypeInferrer who can infer type from a ast.Expr.
type TypeInferrer struct {
	prj *goprj.Project
}

// NewTypeInferrer creates a new TypeInferrer.
func NewTypeInferrer(prj *goprj.Project) *TypeInferrer {
	return &TypeInferrer{prj}
}

// InferType infers type from a ast.Expr.
func (p *TypeInferrer) InferType(expr ast.Expr) goprj.Type {
	switch v := expr.(type) {
	case *ast.CallExpr:
		return p.inferTypeFromFun(v.Fun)
	case *ast.UnaryExpr:
		switch v.Op {
		case token.AND: // address
			t := p.InferType(v.X)
			return p.prj.UniqueType(&goprj.PointerType{Elem: t})
		default:
			log.Fatalln("InferType: unknown UnaryExpr -", v.Op)
		}
	case *ast.SelectorExpr:
		_ = v
	default:
		log.Fatalln("InferType:", reflect.TypeOf(expr))
	}
	return nil
}

func (p *TypeInferrer) inferTypeFromFun(fun ast.Expr) goprj.Type {
	switch v := fun.(type) {
	case *ast.SelectorExpr:
		_ = v
	default:
		log.Fatalln("inferTypeFromFun:", reflect.TypeOf(fun))
	}
	return nil
}

// -----------------------------------------------------------------------------
