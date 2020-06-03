package foo

func ReverseMap(_arg_0 map[string]int) (_ret_1 map[int]string) { 
//line ./foo.ql:4
	return func() (_qlang_ret map[int]string) {
		_qlang_ret = make(map[int]string)
		for k, v := range _arg_0 {
			_qlang_ret[v] = k
		}
		return
	}()
}
