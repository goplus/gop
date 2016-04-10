
grammar = `

term = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

doc = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/pushFloat |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call
`

fntable = nil

Stack = class {

	fn _init() {
		set(this, "stk", [])
	}

	fn clear() {
		set(this, "stk", this.stk[:0])
	}

	fn pop() {
		n = len(this.stk)
		if n > 0 {
			v = this.stk[n-1]
			set(this, "stk", this.stk[:n-1])
			return [v, true]
		}
		return [nil, false]
	}

	fn push(v) {
		set(this, "stk", append(this.stk, v))
	}

	fn popArgs(arity) {
		n = len(this.stk)
		if n < arity {
			panic("Stack.popArgs: unexpected")
		}
		args = slice("var", arity)
		copy(args, this.stk[n-arity:])
		set(this, "stk", this.stk[:n-arity])
		return args
	}
}

Calculator = class {

	fn _init() {
		set(this, "stk", new Stack)
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

main { // 使用main关键字将主程序括起来，是为了避免其中用的局部变量比如 err 对其他函数造成影响

	calc = new Calculator
	engine, err = interpreter(calc, nil)
	if err != nil {
		fprintln(os.stderr, err)
		return 1
	}

	scanner = bufio.scanner(os.stdin)
	for scanner.scan() {
		line = strings.trim(scanner.text(), " \t\r\n")
		if line != "" {
			err = engine.eval(line)
			if err != nil {
				fprintln(os.stderr, err)
			} else {
				printf("> %v\n\n", calc.ret())
			}
		}
	}
}

