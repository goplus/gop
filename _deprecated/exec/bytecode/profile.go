package bytecode

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------

var (
	mutexProf    sync.Mutex
	doProfile    bool
	gOpCalls     [64]uint64
	gOpDurations [64]time.Duration
)

var (
	gInstrCalls     = make(map[Instr]uint64, 1024)
	gInstrDurations = make(map[Instr]time.Duration, 1024)
	gOpDetails      = [64]uint32{
		opBuiltinOp: ^uint32(0),
		opPushInt:   ^uint32(bitsOpIntOperand >> 1),
	}
)

func instrProfile(i Instr, dur time.Duration) {
	mutexProf.Lock()
	defer mutexProf.Unlock()
	op := i >> bitsOpShift
	gOpCalls[op]++
	gOpDurations[op] += dur
	if mask := gOpDetails[op]; mask != 0 {
		gInstrCalls[i&mask]++
		gInstrDurations[i&mask] += dur
	}
}

// -----------------------------------------------------------------------------

func sortOpCalls(opCalls []uint64) []uint32 {
	ops := make([]uint32, 64)
	for i := uint32(0); i < 64; i++ {
		ops[i] = i
	}
	sort.Slice(ops, func(i, j int) bool {
		return gOpDurations[ops[i]] > gOpDurations[ops[j]]
	})
	return ops
}

// -----------------------------------------------------------------------------

type instrCalls struct {
	i     Instr
	dur   time.Duration
	calls uint64
}

func sortInstrCalls() []instrCalls {
	data := make([]instrCalls, 0, len(gInstrCalls))
	for i, calls := range gInstrCalls {
		data = append(data, instrCalls{i, gInstrDurations[i], calls})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].dur > data[j].dur
	})
	return data
}

// -----------------------------------------------------------------------------

// SetProfile sets profile flag.
func SetProfile(profile bool) {
	doProfile = profile
}

// ProfileReport reports profile information.
func ProfileReport() {
	opCalls := gOpCalls[:]
	sortedOps := sortOpCalls(opCalls)
	mostCall := opCalls[sortedOps[0]]
	if mostCall == 0 {
		fmt.Println("No profile information.")
		return
	}
	fmt.Println("Generating op profile information:")
	for i := 0; i < 16; i++ {
		op := sortedOps[i]
		time := gOpDurations[op]
		if time == 0 {
			break
		}
		fmt.Println("op:", op, "time(ms):", time.Milliseconds(), "calls:", opCalls[op])
	}
	fmt.Println("Generating instr profile information:")
	sortedInstrs := sortInstrCalls()
	n := len(sortedInstrs)
	if n > 32 {
		n = 32
	}
	for i := 0; i < n; i++ {
		item := sortedInstrs[i]
		info, p1, p2 := DecodeInstr(item.i)
		time := item.dur.Milliseconds()
		fmt.Println("instr:", item.i, "op:", info.Name, "params:", p1, p2, "time(ms):", time, "calls:", item.calls)
	}
}

// -----------------------------------------------------------------------------
