grammar = `

term1 = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

term2 = term1 *('+' term1/add | '-' term1/sub)

term3 = term2 *('<' term2/lt | '>' term2/gt | "==" term2/eq | "<=" term2/le | ">=" term2/ge | "!=" term2/ne)

term4 = term3 *("&&" term3/and)

expr = term4 *("||" term4/or)

s =
	(IDENT '='! expr)/assign |
	(IDENT ','!)/name IDENT/name % ','/ARITY '=' expr /massign |
	(IDENT "++")/inc |
	(IDENT "--")/dec |
	(IDENT "+="! expr)/adda |
	(IDENT "-="! expr)/suba |
	(IDENT "*="! expr)/mula |
	(IDENT "/="! expr)/quoa |
	(IDENT "%="! expr)/moda |
	"return" ?expr/ARITY /retEngine |
	"defer"/_mute! expr/_codeSrc/_unmute/defer |
	expr

doc = s *(';'/clear s | ';'/pushn)

ifbody = '{' ?doc/_codeSrc '}'

swbody = *("case"! expr/_codeSrc ':' ?doc/_codeSrc)/_ARITY ?("default"! ':' ?doc/_codeSrc)/_ARITY

fnbody = '(' IDENT/name %= ','/ARITY ?"..."/ARITY ')' '{'/_mute ?doc/_codeSrc '}'/_unmute

clsname = '(' IDENT/ref ')' | IDENT/ref

newargs = ?('(' expr %= ','/ARITY ')')/ARITY

classb = "fn"! IDENT/name fnbody ?';'/mfn

atom =
	'(' expr %= ','/ARITY ?"..."/ARITY ?',' ')'/call |
	'.' IDENT/mref |
	'[' ?expr/ARITY ?':'/ARITY ?expr/ARITY ']'/index

factor =
	INT/pushInt |
	FLOAT/pushFloat |
	STRING/pushString |
	CHAR/pushChar |
	(IDENT/ref | '('! expr ')' | "fn"! fnbody/fnEngine | '[' expr %= ','/ARITY ?',' ']'/slice) *atom |
	"if"/_mute! expr/_codeSrc ifbody *("elif" expr/_codeSrc ifbody)/_ARITY ?("else" ifbody)/_ARITY/_unmute/ifEngine |
	"switch"/_mute! ?(~'{' expr)/_codeSrc '{' swbody '}'/_unmute/switchEngine |
	"for"/_mute! (~'{' s)/_codeSrc %= ';'/_ARITY '{' ?doc/_codeSrc '}'/_unmute/forEngine |
	"new"! clsname newargs /new |
	"class"! '{' *classb/ARITY '}'/classEngine |
	'{'! (expr ':' expr) %= ','/ARITY ?',' '}'/map |
	'!' factor/not |
	'-' factor/neg |
	'+' factor
`

errReturn = newRuntimeError("return")

Stack = class {

	fn _init() {
		this.stk = []
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

	fn pushByte(lit) {
		v, multibyte, tail, err = strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
		if err != nil {
			panicf("invalid char `%s`: %v", lit, err)
		}
		if tail != "" || multibyte {
			panic("invalid char: " + lit)
		}
		this.stk = append(this.stk, v)
	}

	fn pushString(lit) {
		v, err = strconv.Unquote(lit)
		if err != nil {
			panicf("invalid string `%s`: %v", lit, err)
		}
		this.stk = append(this.stk, v)
	}

	fn pushNil() {
		this.stk = append(this.stk, nil)
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

	fn baseFrame() {
		return len(this.stk)
	}

	fn setFrame(n) {
		this.stk = this.stk[:n]
	}

	fn index() {
		arity2, _ = this.pop()
		if arity2 != 0 {
			arg2, _ = this.pop()
		}
		arityMid, _ = this.pop()
		arity1, _ = this.pop()
		if arity1 != 0 {
			arg1, _ = this.pop()
		}
		arr, _ = this.pop()

		if arityMid == 0 {
			if arity1 == 0 || arity2 != 0 {
				panic("call operator[] error: illegal index")
			}
			this.push(arr[arg1])
		} else {
			this.push(arr[arg1:arg2])
		}
	}
}

evalCode = fn(e, ip, name, code) {

	old = ip.base
	ip.base = ip.stk.baseFrame()
	err = e.EvalCode(ip, name, code)
	ip.base = old
	return err
}

externVar = class {

	fn _init(vars) {
		this.vars = vars
	}
}

functionInfo = class {

	fn _init(args, fnb, variadic) {
		this.args, this.fnb, this.variadic = args, fnb, variadic
	}
}

Interpreter = class {

	fn _init() {
		this.stk = new Stack
		this.vars = {}
		this.parent = nil
		this.retv = nil
		this.defers = []
		this.base = 0
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

	fn clear() {
		this.stk.setFrame(this.base)
	}

	fn ret() {
		v, _ = this.stk.pop()
		this.clear()
		return v
	}

	fn map() {
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity*2)
		this.stk.push(mapFrom(args...))
	}

	fn slice() {
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity)
		this.stk.push(sliceFrom(args...))
	}

	fn call() {
		variadic, _ = this.stk.pop()
		arity, _ = this.stk.pop()
		if variadic != 0 {
			if arity == 0 {
				panic("what do you mean of `...`?")
			}
			v, _ = this.stk.pop()
			n = len(v)
			for i = 0; i < n; i++ {
				this.stk.push(v[i])
			}
			arity += (n - 1)
		}
		args = this.stk.popArgs(arity)
		f, _ = this.stk.pop()
		this.stk.push(f(args...))
	}

	fn getVar(name) {
		vars = this.vars
		val = vars[name]
		if val == undefined {
			panic("variable not found: " + name)
		}
		if classof(val) == externVar {
			vars = val.vars
			val = vars[name]
		}
		return vars, val
	}

	fn inc(name) {
		vars, val = this.getVar(name)
		val++
		vars[name] = val
		this.stk.push(val)
	}

	fn dec(name) {
		vars, val = this.getVar(name)
		val--
		vars[name] = val
		this.stk.push(val)
	}

	fn addAssign(name) {
		vars, val = this.getVar(name)
		v, ok = this.stk.pop()
		val += v
		vars[name] = val
		this.stk.push(val)
	}

	fn subAssign(name) {
		vars, val = this.getVar(name)
		v, ok = this.stk.pop()
		val -= v
		vars[name] = val
		this.stk.push(val)
	}

	fn mulAssign(name) {
		vars, val = this.getVar(name)
		v, ok = this.stk.pop()
		val *= v
		vars[name] = val
		this.stk.push(val)
	}

	fn quoAssign(name) {
		vars, val = this.getVar(name)
		v, ok = this.stk.pop()
		val /= v
		vars[name] = val
		this.stk.push(val)
	}

	fn modAssign(name) {
		vars, val = this.getVar(name)
		v, ok = this.stk.pop()
		val %= v
		vars[name] = val
		this.stk.push(val)
	}

	fn getVars(name) {
		vars = this.vars
		val = vars[name]
		if val != undefined { // 变量已经存在
			if classof(val) == externVar {
				vars = val.vars
			}
			return vars
		}
		if name[0] != '_' { // NOTE: '_' 开头的变量是私有的，不可继承
			for t = this.parent; t != nil; t = t.parent {
				if t.vars[name] != undefined {
					panicf("variable `%s` exists in extern function", name)
				}
			}
		}
		return vars
	}

	fn getRef(name) {
		val = this.vars[name]
		if val != undefined {
			if classof(val) == externVar {
				val = val.vars[name]
			}
		} else {
			if name[0] != '_' { // NOTE: '_' 开头的变量是私有的，不可继承
				for t = this.parent; t != nil; t = t.parent {
					val = t.vars[name]
					if val != undefined {
						if classof(val) == externVar {
							val = val.vars[name]
						} else {
							e = new externVar(t.vars)
						}
						this.vars[name] = e // 缓存访问过的变量
						return val
					}
				}
			}
			val = fntable[name]
			if val == undefined {
				panic("symbol not found: " + name)
			}
		}
		return val
	}

	fn multiAssign() {
		val, _ = this.stk.pop()
		arity, _ = this.stk.pop()
		arity++
		if arity != len(val) {
			panicf("multi assignment error: require %d variables, but we got %d", len(val), arity)
		}
		names = this.stk.popArgs(arity)
		for i = 0; i < arity; i++ {
			name = names[i]
			vars = this.getVars(name)
			vars[name] = val[i]
		}
		this.stk.push(val)
	}

	fn assign(name) {
		vars = this.getVars(name)
		v, _ = this.stk.pop()
		vars[name] = v
		this.stk.push(v)
	}

	fn ref(name) {
		this.stk.push(this.getRef(name))
	}

	fn function(engine) {
		fnb, _ = this.stk.pop()
		variadic, _ = this.stk.pop()
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity)
		f = new Function
		f.args = args
		f.fnb = fnb
		f.engine = engine
		f.parent = this
		f.cls = nil
		f.variadic = (variadic != 0)
		this.stk.push(f.call)
	}

	fn fnReturn(engine) {
		arity, _ = this.stk.pop()
		if arity == 0 {
			this.retv = nil
		} else {
			v, _ = this.stk.pop()
			this.retv = v
		}
		if this.parent != nil {
			panic(errReturn) // 利用 panic 来实现 return (正常退出)
		}
		this.execDefers(engine)
		if this.retv == nil {
			os.Exit(0)
		}
		os.Exit(this.retv)
	}

	fn execDefers(engine) {
		n = len(this.defers)
		for i = n-1; i >= 0; i-- {
			err = engine.evalCode(this, "expr", this.defers[i])
			if err != nil {
				panic(err)
			}
		}
		this.defers = this.defers[:0]
	}

	fn fnDefer() {
		src, _ = this.stk.pop()
		this.defers = append(this.defers, src)
		this.stk.push(nil)
	}

	fn memberFuncDecl() {
		fnb, _ = this.stk.pop()
		variadic, _ = this.stk.pop()
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity + 1)
		f = new functionInfo(args, fnb, variadic != 0)
		this.stk.push(f)
	}

	fn fnClass(engine) {
		arity, _ = this.stk.pop()
		args = this.stk.popArgs(arity)
		fns = {}
		cls = new Class(fns, engine, this)
		for i = 0; i < len(args); i++ {
			v = args[i]
			name = v.args[0]
			v.args[0] = "this"
			f = new Function
			f.args = v.args
			f.fnb = v.fnb
			f.engine = engine
			f.parent = this
			f.cls = cls
			f.variadic = v.variadic
			fns[name] = f.call
		}
		this.stk.push(cls)
	}

	fn fnNew() {

		args = []
		hasArgs, _ = this.stk.pop()
		if hasArgs != 0 {
			nArgs, _ = this.stk.pop()
			args = this.stk.popArgs(nArgs)
		}

		cls, _ = this.stk.pop()
		if classof(cls) != Class {
			panic("new argument isn't a class")
		}
		obj = new Object(cls)
		init = cls.fns["_init"]
		if init != undefined { // 构造函数
			init(obj, args...)
		}
		this.stk.push(obj)
	}

	fn memberRef(name) {

		o, _ = this.stk.pop()
		switch classof(o) {
		case Object:
			val = o.vars[name]
			if val == undefined {
				f = o.cls.fns[name]
				if f != undefined {
					val = new thisDeref(o, f)
				} else {
					panicf("object doesn't has member `%s`", name)
				}
			}
			this.stk.push(val)
			return
		case Class:
			val = o.fns[name]
			if val == undefined {
				panicf("class doesn't has method `%s`", name)
			}
			this.stk.push(val)
			return
		}

		obj = reflect.ValueOf(o)
		switch {
		case obj.Kind() == reflect.Map:
			m = obj.MapIndex(reflect.ValueOf(name))
			if m.IsValid() {
				this.stk.push(m.Interface())
			} else {
				panicf("member `%s` not found", name)
			}
		default:
			name = strings.Title(name)
			m = obj.MethodByName(name)
			if !m.IsValid() {
				m = reflect.Indirect(obj).FieldByName(name)
				if !m.IsValid() {
					panicf("type `%v` doesn't has member `%s`", obj.Type(), name)
				}
			}
			this.stk.push(m.Interface())
		}
	}

	fn fnIf(engine) {

		elseCode = nil
		elseArity, _ = this.stk.pop()
		if elseArity != 0 {
			elseCode, _ = this.stk.pop()
		}

		condArity, _ = this.stk.pop()
		n = 2 * (condArity + 1)
		ifbr = this.stk.popArgs(n)

		for i = 0; i <= condArity; i++ {
			condCode = ifbr[i*2]
			err = evalCode(engine, this, "expr", condCode)
			if err != nil {
				return err
			}
			cond, _ = this.stk.pop()
			if cond {
				bodyCode = ifbr[i*2+1]
				if bodyCode == nil {
					this.stk.push(nil)
					return
				}
				return evalCode(engine, this, "doc", bodyCode)
			}
		}
		if elseCode != nil {
			return evalCode(engine, this, "doc", elseCode)
		}
		this.stk.push(nil)
	}

	fn fnSwitch(engine) {

		defaultCode = nil
		defaultArity, _ = this.stk.pop()
		if defaultArity != 0 {
			defaultCode, _ = this.stk.pop()
		}

		caseArity, _ = this.stk.pop()
		n = 2 * caseArity
		casebr = this.stk.popArgs(n)

		switchCode, _ = this.stk.pop()
		if switchCode == nil {
			switchCond = true
		} else {
			err = evalCode(engine, this, "expr", switchCode)
			if err != nil {
				return err
			}
			switchCond, _ = this.stk.pop()
		}

		for i = 0; i < caseArity; i++ {
			caseCode = casebr[i*2]
			err = evalCode(engine, this, "expr", caseCode)
			if err != nil {
				return err
			}
			caseCond, _ = this.stk.pop()
			if switchCond == caseCond {
				bodyCode = casebr[i*2+1]
				if bodyCode == nil {
					this.stk.push(nil)
					return
				}
				return evalCode(engine, this, "doc", bodyCode)
			}
		}
		if defaultCode != nil {
			return evalCode(engine, this, "doc", defaultCode)
		}
		this.stk.push(nil)
	}

	fn fnFor(engine) {

		bodyCode, _ = this.stk.pop()

		condCode = nil
		stepCode = nil
		arity, _ = this.stk.pop()
		switch {
		case arity == 0:
		case arity == 1 || arity == 3:
			forCode = this.stk.popArgs(arity)
			if arity == 3 {
				initCode = forCode[0]
				if initCode != nil {
					evalCode(engine, this, "s", initCode)
					this.stk.pop()
				}
				stepCode = forCode[2]
			}
			condCode = forCode[arity/2]
		default:
			panic("invalid for loop form")
		}

		n = 0
		for condCode == nil || this.evalCondS(engine, condCode) {
			n++
			if bodyCode != nil {
				err = evalCode(engine, this, "doc", bodyCode)
				if err != nil {
					return err
				}
				this.stk.pop()
			}
			if stepCode != nil {
				err = evalCode(engine, this, "s", stepCode)
				if err != nil {
					return err
				}
				this.stk.pop()
			}
		}
		this.stk.push(n)
		return nil
	}

	fn evalCondS(engine, condCode) {
		err = evalCode(engine, this, "s", condCode)
		if err != nil {
			panic(err)
		}
		cond, _ = this.stk.pop()
		return cond
	}
}

Class = class {

	fn _init(fns, engine, parent) {
		this.fns, this.engine, this.parent = fns, engine, parent
	}
}

Object = class {

	fn _init(cls) {
		this.vars, this.cls = {}, cls
	}

	fn setVar(name, val) {
		f = this.cls.fns[name]
		if f != undefined {
			panic("set failed: class already have a method named " + name)
		}
		this.vars[name] = val
	}
}

Function = class {

	fn call(args...) {

		n = len(this.args)
		if this.variadic {
			if len(args) < n-1 {
				panicf("function requires >= %d arguments, but we got %d", n-1, len(args))
			}
		} else {
			if len(args) != n {
				panicf("function requires %d arguments, but we got %d", n, len(args))
			}
		}

		if this.fnb == nil {
			return nil
		}

		parent = this.parent
		stk = parent.stk
		vars = {}
		base = stk.baseFrame()
		ip = new Interpreter
		ip.vars = vars
		ip.parent = parent
		ip.stk = stk
		ip.base = base

		if this.variadic {
			for i = 0; i < n-1; i++ {
				vars[this.args[i]] = args[i]
			}
			vars[this.args[n-1]] = args[n-1:]
		} else {
			for i = 0; i < len(args); i++ {
				vars[this.args[i]] = args[i]
			}
		}

		err = fn() {
			defer fn() {
				ip.execDefers(this.engine)
				stk.setFrame(base)
				e = recover()
				if e != nil {
					if e == errReturn { // 正常 return 导致，见 fnReturn 函数
						return
					}
					panic(e)
				}
			}()
			return this.engine.evalCode(ip, "doc", this.fnb)
		}()
		if err != nil {
			panic(err)
		}
		return ip.retv
	}
}

thisDeref = class {

	fn _init(p, f) {
		this.p, this.f = p, f
	}

	fn call(args...) {
		return this.f.call(this.p, args...)
	}
}

fnSet = fn(o, args...) {
	if classof(o) == Object {
		for i = 0; i < len(args); i += 2 {
			o.setVar(args[i], args[i+1])
		}
	} else {
		set(o, args...)
	}
}

osMod = {
	"Args":   os.Args[1:],
	"Open":   os.Open,
	"Stderr": os.Stderr,
	"Stdout": os.Stdout,
	"Stdin":  os.Stdin,
}

fntable = {
	"undefined": undefined,
	"nil":       nil,
	"true":      true,
	"false":     false,

	"sin":      sin,
	"cos":      cos,
	"pow":      pow,
	"max":      max,
	"min":      min,
	"len":      len,
	"append":   append,
	"type":     reflect.TypeOf,
	"printf":   printf,
	"println":  println,
	"fprintln": fprintln,
	"set":      fnSet,
	"os":       osMod,
	"strconv":  strconv,
	"strings":  strings,

	"$add": fn(a, b) { return a+b },
	"$sub": fn(a, b) { return a-b },
	"$mul": fn(a, b) { return a*b },
	"$quo": fn(a, b) { return a/b },
	"$mod": fn(a, b) { return a%b },
	"$neg": fn(a)    { return -a  },

	"$lt":  fn(a, b) { return a<b },
	"$gt":  fn(a, b) { return a>b },
	"$le":  fn(a, b) { return a<=b },
	"$ge":  fn(a, b) { return a>=b },
	"$eq":  fn(a, b) { return a==b },
	"$ne":  fn(a, b) { return a!=b },
	"$and": fn(a, b) { return a&&b },
	"$or":  fn(a, b) { return a||b },
	"$not": fn(a)    { return !a   },

	"$name":       Stack.push,
	"$ARITY":      Stack.push,
	"$_ARITY":     Stack.push,
	"$_codeSrc":   Stack.push,
	"$pushInt":    Stack.push,
	"$pushFloat":  Stack.push,
	"$pushString": Stack.pushString,
	"$pushChar":   Stack.pushByte,
	"$pushn":      Stack.pushNil,
	"$index":      Stack.index,

	"$clear":   Interpreter.clear,
	"$ref":     Interpreter.ref,
	"$slice":   Interpreter.slice,
	"$map":     Interpreter.map,
	"$call":    Interpreter.call,
	"$assign":  Interpreter.assign,
	"$massign": Interpreter.multiAssign,
	"$inc":     Interpreter.inc,
	"$dec":     Interpreter.dec,
	"$adda":    Interpreter.addAssign,
	"$suba":    Interpreter.subAssign,
	"$mula":    Interpreter.mulAssign,
	"$quoa":    Interpreter.quoAssign,
	"$moda":    Interpreter.modAssign,
	"$defer":   Interpreter.fnDefer,
	"$mfn":     Interpreter.memberFuncDecl,
	"$mref":    Interpreter.memberRef,
	"$new":     Interpreter.fnNew,

	"$classEngine":  Interpreter.fnClass,
	"$ifEngine":     Interpreter.fnIf,
	"$switchEngine": Interpreter.fnSwitch,
	"$forEngine":    Interpreter.fnFor,

	"$retEngine": Interpreter.fnReturn,
	"$fnEngine":  Interpreter.function,
}
