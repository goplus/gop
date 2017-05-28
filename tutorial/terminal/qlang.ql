include "../qlang/qlang.ql"

main { // 使用main关键字将主程序括起来，是为了避免其中用的局部变量比如 err 对其他函数造成影响

	ipt = new Interpreter
	engine, err = interpreter(ipt, nil)
	if err != nil {
		fprintln(os.Stderr, err)
		return 1
	}

	historyFile = os.Getenv("HOME") + "/.qlang.history"
	term = terminal.New(">>> ", "... ", nil)
	term.LoadHistroy(historyFile)
	defer term.SaveHistroy(historyFile)

	println(`Q-language - http://qlang.io, version 1.0.00
Copyright (C) 2015 Qiniu.com - Shanghai Qiniu Information Technologies Co., Ltd.
Use Ctrl-D (i.e. EOF) to exit.
`)

	for {
		expr, err = term.Scan()
		if err != nil {
			if err == terminal.ErrPromptAborted {
				continue
			} elif err == io.EOF {
				println("^D")
				break
			}
			fprintln(os.Stderr, err)
			continue
		}
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}
		err = engine.Eval(expr)
		if err != nil {
			fprintln(os.Stderr, err)
		} else {
			println(ipt.ret())
		}
	}
}
