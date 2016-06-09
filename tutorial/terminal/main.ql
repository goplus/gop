include "../calc/calc.ql"

main { // 使用main关键字将主程序括起来，是为了避免其中用的局部变量比如 err 对其他函数造成影响

	calc = new Calculator
	engine, err = interpreter(calc, nil)
	if err != nil {
		fprintln(os.stderr, err)
		return 1
	}

	historyFile = os.getenv("HOME") + "/.qcalc.history"
	term = terminal.new(">>> ", "... ", nil)
	term.loadHistroy(historyFile)
	defer term.saveHistroy(historyFile)

	println(`Q-Calculator - http://qlang.io, version 1.0.00
Copyright (C) 2015 Qiniu.com - Shanghai Qiniu Information Technologies Co., Ltd.
Use Ctrl-D (i.e. EOF) to exit.
`)

	for {
		expr, err = term.scan()
		if err != nil {
			if err == terminal.ErrPromptAborted {
				continue
			} elif err == io.EOF {
				println("^D")
				break
			}
			fprintln(os.stderr, err)
			continue
		}
		expr = strings.trimSpace(expr)
		if expr == "" {
			continue
		}
		err = engine.eval(expr)
		if err != nil {
			fprintln(os.stderr, err)
		} else {
			println(calc.ret())
		}
	}
}
