
grammar = `

term = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

doc = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/pushFloat |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call
`

Stack = class {

	fn _init() {
		this.stk = []
	}

	fn clear() {
		this.stk = this.stk[:0]
	}

	fn pop() {
		n = len(this.stk)
		if n > 0 {
			v = this.stk[n-1]
			this.stk = this.stk[:n-1]
			return v, true
		}
		return nil, false
	}

	fn push(v) {
		this.stk = append(this.stk, v)
	}

	fn popArgs(arity) {
		n = len(this.stk)
		if n < arity {
			panic("Stack.popArgs: unexpected")
		}
		args = make([]var, arity)
		copy(args, this.stk[n-arity:])
		this.stk = this.stk[:n-arity]
		return args
	}
}

Calculator = class {

	fn _init() {
		this.stk = new Stack
	}

	fn grammar() {
		return grammar
	}

	fn stack() {
		return this.stk
	}

	fn fntable() {
		return fntable
	}

	fn ret() {
		v, _ = this.stk.pop()
		this.stk.clear()
		return v
	}

	fn call(name) {
		f = fntable[name]
		if f == undefined {
			panic("function not found: " + name)
		}
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity)
		ret = f(args...)
		this.stk.push(ret)
	}
}

fntable = {
	"sin": sin,
	"cos": cos,
	"pow": pow,
	"max": max,
	"min": min,

	"$mul": fn(a, b) { return a*b },
	"$quo": fn(a, b) { return a/b },
	"$mod": fn(a, b) { return a%b },
	"$add": fn(a, b) { return a+b },
	"$sub": fn(a, b) { return a-b },
	"$neg": fn(a) { return -a },

	"$call": Calculator.call,
	"$pushFloat": Stack.push,
	"$ARITY": Stack.push,
}
