include "qlang.ql"

main { // 使用main关键字将主程序括起来，是为了避免其中用的局部变量比如 err 对其他函数造成影响

	ipt = new Interpreter

	if len(os.Args) > 1 {
		engine, err = interpreter(ipt, insertSemis)
		if err != nil {
			fprintln(os.Stderr, err)
			return 1
		}
		fname = os.Args[1]
		b, err = ioutil.ReadFile(fname)
		if err != nil {
			fprintln(os.Stderr, err)
			return 2
		}
		err = engine.Exec(b, fname)
		if err != nil {
			fprintln(os.Stderr, err)
			return 3
		}
		return
	}

	engine, err = interpreter(ipt, nil)
	if err != nil {
		fprintln(os.Stderr, err)
		return 1
	}

	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = strings.Trim(scanner.Text(), " \t\r\n")
		if line != "" {
			err = engine.Eval(line)
			if err != nil {
				fprintln(os.Stderr, err)
			} else {
				printf("> %v\n\n", ipt.ret())
			}
		}
	}
}
