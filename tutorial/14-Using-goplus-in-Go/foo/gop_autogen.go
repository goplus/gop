package foo

func ReverseMap(_arg_0 map[string]int) (_ret_1 map[int]string) {
//line ./foo.gop:4
	return func() (_gop_ret map[int]string) {
		_gop_ret = make(map[int]string)
		for k, v := range _arg_0 {
//line ./foo.gop:4
			_gop_ret[v] = k
		}
		return
	}()
}
