// +build go1.6

package runtime

import "runtime"

func init() {
	Exports["readTrace"] = runtime.ReadTrace
	Exports["startTrace"] = runtime.StartTrace
	Exports["stopTrace"] = runtime.StopTrace
}
