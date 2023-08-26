package ng

import (
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"strconv"
)

const (
	GopPackage = true // to indicate this is a Go+ package
)

const (
	Uint128_Max       = 1<<128 - 1
	Uint128_Min       = 0
	Uint128_IsUntyped = true
)

const (
	maxUint64 = (1 << 64) - 1
	intSize   = 32 << (^uint(0) >> 63)
)

//
// Gop_: Go+ object prefix
// Gop_xxx_Cast: type Gop_xxx typecast
// xxxx__N: the Nth overload function
//

type Gop_ninteger = uint

// -----------------------------------------------------------------------------

type Uint128 struct {
	hi, lo uint64
}

// Uint128_Init: func uint128.init(v int) uint128
func Uint128_Init__0(v int) (out Uint128) {
	if v < 0 {
		panic("TODO: can't init uint128 from a negative integer")
	}
	return Uint128{lo: uint64(v)}
}

// Uint128_Init: func bigint.init(v untyped_bigint) bigint
func Uint128_Init__1(v UntypedBigint) (out Uint128) {
	out, inRange := Uint128_Cast__9(v)
	if !inRange {
		log.Panicf("value %v was not in valid uint128 range\n", v)
	}
	return
}

// Uint128_Cast: func uint128(v untyped_int) uint128
func Uint128_Cast__0(v int) Uint128 {
	return Uint128_Cast__6(int64(v))
}

// Uint128_Cast: func uint128(v untyped_bigint) uint128
func Uint128_Cast__1(v UntypedBigint) Uint128 {
	return Uint128_Init__1(v)
}

// Uint128_Cast: func uint128(v uint64) uint128
func Uint128_Cast__2(v uint64) Uint128 {
	return Uint128{lo: v}
}

// Uint128_Cast: func uint128(v uint32) uint128
func Uint128_Cast__3(v uint32) Uint128 {
	return Uint128{lo: uint64(v)}
}

// Uint128_Cast: func uint128(v uint16) uint128
func Uint128_Cast__4(v uint16) Uint128 {
	return Uint128{lo: uint64(v)}
}

// Uint128_Cast: func uint128(v uint8) uint128
func Uint128_Cast__5(v uint8) Uint128 {
	return Uint128{lo: uint64(v)}
}

// Uint128_Cast: func uint128(v int64) uint128
func Uint128_Cast__6(v int64) (out Uint128) {
	if v < 0 {
		return Uint128{hi: maxUint64, lo: uint64(v)}
	}
	return Uint128{lo: uint64(v)}
}

// Uint128_Cast: func uint128(v int64) (uint128, bool)
func Uint128_Cast__7(v int64) (out Uint128, inRange bool) {
	if v < 0 {
		return
	}
	return Uint128{lo: uint64(v)}, true
}

// Uint128_Cast: func uint128(v *big.Int) uint128
func Uint128_Cast__8(v *big.Int) Uint128 {
	out, _ := Uint128_Cast__9(v)
	return out
}

// Uint128_Cast: func uint128(v *big.Int) (uint128, bool)
func Uint128_Cast__9(v *big.Int) (out Uint128, inRange bool) {
	if v.Sign() < 0 {
		return out, false
	}

	words := v.Bits()
	switch intSize {
	case 64:
		lw := len(words)
		switch lw {
		case 0:
			return Uint128{}, true
		case 1:
			return Uint128{lo: uint64(words[0])}, true
		case 2:
			return Uint128{hi: uint64(words[1]), lo: uint64(words[0])}, true
		default:
			return Uint128{hi: maxUint64, lo: maxUint64}, false
		}

	case 32:
		lw := len(words)
		switch lw {
		case 0:
			return Uint128{}, true
		case 1:
			return Uint128{lo: uint64(words[0])}, true
		case 2:
			return Uint128{lo: (uint64(words[1]) << 32) | (uint64(words[0]))}, true
		case 3:
			return Uint128{hi: uint64(words[2]), lo: (uint64(words[1]) << 32) | (uint64(words[0]))}, true
		case 4:
			return Uint128{
				hi: (uint64(words[3]) << 32) | (uint64(words[2])),
				lo: (uint64(words[1]) << 32) | (uint64(words[0])),
			}, true
		default:
			return Uint128{hi: maxUint64, lo: maxUint64}, false
		}

	default:
		panic("unsupported bit size")
	}
}

// Uint128_Cast: func uint128() uint128
func Uint128_Cast__a() Uint128 {
	return Uint128{}
}

// Uint128_Cast: func uint128(v uint) uint128
func Uint128_Cast__b(v uint) Uint128 {
	return Uint128{lo: uint64(v)}
}

// Uint128_Cast: func uint128(hi, lo uint64) uint128
func Uint128_Cast__c(hi, lo uint64) Uint128 {
	return Uint128{hi: hi, lo: lo}
}

// Gop_Rcast: func float64(v uint128) float64
func (u Uint128) Gop_Rcast__0() float64 {
	if u.hi == 0 {
		return float64(u.lo)
	}
	return (float64(u.hi) * (1 << 64)) + float64(u.lo)
}

// Gop_Rcast: func uint64(v uint128) uint64
func (u Uint128) Gop_Rcast__1() uint64 {
	return u.lo
}

// Gop_Rcast: func uint64(v uint128) (uint64, bool)
func (u Uint128) Gop_Rcast__2() (out uint64, inRange bool) {
	return u.lo, u.hi == 0
}

// Gop_Rcast: func int64(v uint128) int64
func (u Uint128) Gop_Rcast__3() int64 {
	return int64(u.lo)
}

// Gop_Rcast: func int64(v uint128) (int64, bool)
func (u Uint128) Gop_Rcast__4() (out int64, inRange bool) {
	return int64(u.lo), u.hi == 0 && u.lo <= maxInt64
}

// -----------------------------------------------------------------------------

func (u Uint128) IsZero() bool {
	return u.lo == 0 && u.hi == 0
}

func (u *Uint128) Scan(state fmt.ScanState, verb rune) (err error) {
	t, err := state.Token(true, nil)
	if err != nil {
		return
	}
	v, err := ParseUint128(string(t), 10)
	if err == nil {
		*u = v
	}
	return
}

func (u Uint128) Format(s fmt.State, c rune) {
	// TODO: not so good
	u.BigInt().Format(s, c)
}

func (u Uint128) String() string {
	return u.Text(10)
}

func (u Uint128) Text(base int) string {
	if u.hi == 0 {
		return strconv.FormatUint(u.lo, base)
	}
	// TODO: not so good
	return u.BigInt().Text(base)
}

func (u Uint128) BigInt() *big.Int {
	var v big.Int
	u.ToBigInt(&v)
	return &v
}

func (u Uint128) ToBigInt(b *big.Int) {
	switch intSize {
	case 64:
		bits := b.Bits()
		ln := len(bits)
		if len(bits) < 2 {
			bits = append(bits, make([]big.Word, 2-ln)...)
		}
		bits = bits[:2]
		bits[0] = big.Word(u.lo)
		bits[1] = big.Word(u.hi)
		b.SetBits(bits)

	case 32:
		bits := b.Bits()
		ln := len(bits)
		if len(bits) < 4 {
			bits = append(bits, make([]big.Word, 4-ln)...)
		}
		bits = bits[:4]
		bits[0] = big.Word(u.lo & 0xFFFFFFFF)
		bits[1] = big.Word(u.lo >> 32)
		bits[2] = big.Word(u.hi & 0xFFFFFFFF)
		bits[3] = big.Word(u.hi >> 32)
		b.SetBits(bits)

	default:
		if u.hi > 0 {
			b.SetUint64(u.hi)
			b.Lsh(b, 64)
		}
		var lo big.Int
		lo.SetUint64(u.lo)
		b.Add(b, &lo)
	}
}

// Bit returns the value of the i'th bit of x. That is, it returns (x>>i)&1.
// The bit index i must be 0 <= i < 128
func (u Uint128) Bit(i int) uint {
	if i < 0 || i >= 128 {
		panic("bit out of range")
	}
	if i >= 64 {
		return uint((u.hi >> uint(i-64)) & 1)
	} else {
		return uint((u.lo >> uint(i)) & 1)
	}
}

// SetBit returns a Uint128 with u's i'th bit set to b (0 or 1).
// If b is not 0 or 1, SetBit will panic. If i < 0, SetBit will panic.
func (u Uint128) SetBit(i int, b uint) (out Uint128) {
	if i < 0 || i >= 128 {
		panic("bit out of range")
	}
	if b == 0 {
		if i >= 64 {
			u.hi = u.hi &^ (1 << uint(i-64))
		} else {
			u.lo = u.lo &^ (1 << uint(i))
		}
	} else if b == 1 {
		if i >= 64 {
			u.hi = u.hi | (1 << uint(i-64))
		} else {
			u.lo = u.lo | (1 << uint(i))
		}
	} else {
		panic("bit value not 0 or 1")
	}
	return u
}

func (u Uint128) LeadingZeros() int {
	if u.hi == 0 {
		return bits.LeadingZeros64(u.lo) + 64
	} else {
		return bits.LeadingZeros64(u.hi)
	}
}

func (u Uint128) TrailingZeros() int {
	if u.lo == 0 {
		return bits.TrailingZeros64(u.hi) + 64
	} else {
		return bits.TrailingZeros64(u.lo)
	}
}

// BitLen returns the length of the absolute value of u in bits. The bit length of 0 is 0.
func (u Uint128) BitLen() int {
	if u.hi > 0 {
		return bits.Len64(u.hi) + 64
	}
	return bits.Len64(u.lo)
}

// OnesCount returns the number of one bits ("population count") in u.
func (u Uint128) OnesCount() int {
	return bits.OnesCount64(u.hi) + bits.OnesCount64(u.lo)
}

func (u Uint128) Reverse() Uint128 {
	return Uint128{hi: bits.Reverse64(u.lo), lo: bits.Reverse64(u.hi)}
}

func (u Uint128) ReverseBytes() Uint128 {
	return Uint128{hi: bits.ReverseBytes64(u.lo), lo: bits.ReverseBytes64(u.hi)}
}

// Cmp compares 'u' to 'n' and returns:
//
//	< 0 if u <  n
//	  0 if u == n
//	> 0 if u >  n
//
// The specific value returned by Cmp is undefined, but it is guaranteed to
// satisfy the above constraints.
func (u Uint128) Cmp__1(n Uint128) int {
	if u.hi == n.hi {
		if u.lo > n.lo {
			return 1
		} else if u.lo < n.lo {
			return -1
		}
	} else {
		if u.hi > n.hi {
			return 1
		} else if u.hi < n.hi {
			return -1
		}
	}
	return 0
}

func (u Uint128) Cmp__0(n uint64) int {
	if u.hi > 0 || u.lo > n {
		return 1
	} else if u.lo < n {
		return -1
	}
	return 0
}

func (u Uint128) Gop_Dup() (v Uint128) {
	return u
}

func (u *Uint128) Gop_Inc() {
	u.lo++
	if u.lo == 0 {
		u.hi++
	}
}

func (u *Uint128) Gop_Dec() {
	if u.lo == 0 {
		u.hi--
	}
	u.lo--
}

// Gop_AddAssign: func (a *uint128) += (b uint128)
func (u *Uint128) Gop_AddAssign(b Uint128) {
	*u = u.Gop_Add__1(b)
}

// Gop_SubAssign: func (a *uint128) -= (b uint128)
func (u *Uint128) Gop_SubAssign(b Uint128) {
	*u = u.Gop_Sub__1(b)
}

// Gop_MulAssign: func (a *uint128) *= (b uint128)
func (u *Uint128) Gop_MulAssign(b Uint128) {
	*u = u.Gop_Mul__1(b)
}

// Gop_QuoAssign: func (a *uint128) /= (b uint128) {
func (u *Uint128) Gop_QuoAssign(b Uint128) {
	*u = u.Gop_Quo__1(b)
}

// Gop_RemAssign: func (a *uint128) %= (b uint128)
func (u *Uint128) Gop_RemAssign(b Uint128) {
	*u = u.Gop_Rem__1(b)
}

// Gop_OrAssign: func (a *uint128) |= (b uint128)
func (u *Uint128) Gop_OrAssign(b Uint128) {
	*u = u.Gop_Or__1(b)
}

// Gop_XorAssign: func (a *uint128) ^= (b uint128)
func (u *Uint128) Gop_XorAssign(b Uint128) {
	*u = u.Gop_Xor__1(b)
}

// Gop_AndAssign: func (a *uint128) &= (b uint128)
func (u *Uint128) Gop_AndAssign(b Uint128) {
	*u = u.Gop_And__1(b)
}

// Gop_AndNotAssign: func (a *uint128) &^= (b uint128)
func (u *Uint128) Gop_AndNotAssign(b Uint128) {
	*u = u.Gop_AndNot(b)
}

// Gop_LshAssign: func (a *uint128) <<= (n untyped_uint)
func (u *Uint128) Gop_LshAssign(n Gop_ninteger) {
	*u = u.Gop_Lsh(n)
}

// Gop_RshAssign: func (a *uint128) >>= (n untyped_uint)
func (u *Uint128) Gop_RshAssign(n Gop_ninteger) {
	*u = u.Gop_Rsh(n)
}

func (u Uint128) Gop_Add__1(n Uint128) (v Uint128) {
	var carry uint64
	v.lo, carry = bits.Add64(u.lo, n.lo, 0)
	v.hi, _ = bits.Add64(u.hi, n.hi, carry)
	return v
}

func (u Uint128) Gop_Add__0(n uint64) (v Uint128) {
	var carry uint64
	v.lo, carry = bits.Add64(u.lo, n, 0)
	v.hi = u.hi + carry
	return v
}

func (u Uint128) Gop_Sub__1(n Uint128) (v Uint128) {
	var borrowed uint64
	v.lo, borrowed = bits.Sub64(u.lo, n.lo, 0)
	v.hi, _ = bits.Sub64(u.hi, n.hi, borrowed)
	return v
}

func (u Uint128) Gop_Sub__0(n uint64) (v Uint128) {
	var borrowed uint64
	v.lo, borrowed = bits.Sub64(u.lo, n, 0)
	v.hi = u.hi - borrowed
	return v
}

func (u Uint128) Gop_EQ__1(n Uint128) bool {
	return u.hi == n.hi && u.lo == n.lo
}

func (u Uint128) Gop_EQ__0(n uint64) bool {
	return u.hi == 0 && u.lo == n
}

func (u Uint128) Gop_GT__1(n Uint128) bool {
	return u.hi > n.hi || (u.hi == n.hi && u.lo > n.lo)
}

func (u Uint128) Gop_GT__0(n uint64) bool {
	return u.hi > 0 || u.lo > n
}

func (u Uint128) Gop_GE__1(n Uint128) bool {
	return u.hi > n.hi || (u.hi == n.hi && u.lo >= n.lo)
}

func (u Uint128) Gop_GE__0(n uint64) bool {
	return u.hi > 0 || u.lo >= n
}

func (u Uint128) Gop_LT__1(n Uint128) bool {
	return u.hi < n.hi || (u.hi == n.hi && u.lo < n.lo)
}

func (u Uint128) Gop_LT__0(n uint64) bool {
	return u.hi == 0 && u.lo < n
}

func (u Uint128) Gop_LE__1(n Uint128) bool {
	return u.hi < n.hi || (u.hi == n.hi && u.lo <= n.lo)
}

func (u Uint128) Gop_LE__0(n uint64) bool {
	return u.hi == 0 && u.lo <= n
}

func (u Uint128) Gop_And__1(n Uint128) Uint128 {
	u.hi &= n.hi
	u.lo &= n.lo
	return u
}

func (u Uint128) Gop_And__0(n uint64) Uint128 {
	return Uint128{lo: u.lo & n}
}

func (u Uint128) Gop_AndNot(n Uint128) Uint128 {
	u.hi &^= n.hi
	u.lo &^= n.lo
	return u
}

func (u Uint128) Gop_Not() Uint128 {
	return Uint128{hi: ^u.hi, lo: ^u.lo}
}

func (u Uint128) Gop_Or__1(n Uint128) Uint128 {
	u.hi |= n.hi
	u.lo |= n.lo
	return u
}

func (u Uint128) Gop_Or__0(n uint64) Uint128 {
	u.lo |= n
	return u
}

func (u Uint128) Gop_Xor__1(v Uint128) Uint128 {
	u.hi ^= v.hi
	u.lo ^= v.lo
	return u
}

func (u Uint128) Gop_Xor__0(v uint64) Uint128 {
	u.lo ^= v
	return u
}

func (u Uint128) Gop_Lsh(n Gop_ninteger) Uint128 {
	if n < 64 {
		u.hi = (u.hi << n) | (u.lo >> (64 - n))
		u.lo <<= n
	} else {
		u.hi = u.lo << (n - 64)
		u.lo = 0
	}
	return u
}

func (u Uint128) Gop_Rsh(n Gop_ninteger) Uint128 {
	if n < 64 {
		u.lo = (u.lo >> n) | (u.hi << (64 - n))
		u.hi >>= n
	} else {
		u.lo = u.hi >> (n - 64)
		u.hi = 0
	}
	return u
}

func (u Uint128) Gop_Mul__1(n Uint128) Uint128 {
	hi, lo := bits.Mul64(u.lo, n.lo)
	hi += u.hi*n.lo + u.lo*n.hi
	return Uint128{hi, lo}
}

func (u Uint128) Gop_Mul__0(n uint64) (dest Uint128) {
	dest.hi, dest.lo = bits.Mul64(u.lo, n)
	dest.hi += u.hi * n
	return dest
}

const (
	divAlgoLeading0Spill = 16
)

func (u Uint128) Gop_Quo__1(by Uint128) (q Uint128) {
	if by.lo == 0 && by.hi == 0 {
		panic("division by zero")
	}

	if u.hi|by.hi == 0 {
		q.lo = u.lo / by.lo // FIXME: div/0 risk?
		return q
	}

	var byLoLeading0, byHiLeading0, byLeading0 uint
	if by.hi == 0 {
		byLoLeading0, byHiLeading0 = uint(bits.LeadingZeros64(by.lo)), 64
		byLeading0 = byLoLeading0 + 64
	} else {
		byHiLeading0 = uint(bits.LeadingZeros64(by.hi))
		byLeading0 = byHiLeading0
	}

	if byLeading0 == 127 {
		return u
	}

	byTrailing0 := uint(by.TrailingZeros())
	if (byLeading0 + byTrailing0) == 127 {
		return u.Gop_Rsh(byTrailing0)
	}

	if cmp := u.Cmp__1(by); cmp < 0 {
		return q // it's 100% remainder
	} else if cmp == 0 {
		q.lo = 1 // dividend and divisor are the same
		return q
	}

	uLeading0 := uint(u.LeadingZeros())
	if byLeading0-uLeading0 > divAlgoLeading0Spill {
		q, _ = quorem128by128(u, by, byHiLeading0, byLoLeading0)
		return q
	} else {
		return quo128bin(u, by, uLeading0, byLeading0)
	}
}

func (u Uint128) Gop_Quo__0(by uint64) (q Uint128) {
	if u.hi < by {
		q.lo, _ = bits.Div64(u.hi, u.lo, by)
	} else {
		q.hi = u.hi / by
		q.lo, _ = bits.Div64(u.hi%by, u.lo, by)
	}
	return q
}

func (u Uint128) QuoRem__1(by Uint128) (q, r Uint128) {
	if by.lo == 0 && by.hi == 0 {
		panic("division by zero")
	}

	if u.hi|by.hi == 0 {
		q.lo = u.lo / by.lo
		r.lo = u.lo % by.lo
		return q, r
	}

	var byLoLeading0, byHiLeading0, byLeading0 uint
	if by.hi == 0 {
		byLoLeading0, byHiLeading0 = uint(bits.LeadingZeros64(by.lo)), 64
		byLeading0 = byLoLeading0 + 64
	} else {
		byHiLeading0 = uint(bits.LeadingZeros64(by.hi))
		byLeading0 = byHiLeading0
	}

	if byLeading0 == 127 {
		return u, r
	}

	byTrailing0 := uint(by.TrailingZeros())
	if (byLeading0 + byTrailing0) == 127 {
		q = u.Gop_Rsh(byTrailing0)
		by.Gop_Dec()
		r = by.Gop_And__1(u)
		return
	}

	if cmp := u.Cmp__1(by); cmp < 0 {
		return q, u // it's 100% remainder

	} else if cmp == 0 {
		q.lo = 1 // dividend and divisor are the same
		return q, r
	}

	uLeading0 := uint(u.LeadingZeros())
	if byLeading0-uLeading0 > divAlgoLeading0Spill {
		return quorem128by128(u, by, byHiLeading0, byLoLeading0)
	} else {
		return quorem128bin(u, by, uLeading0, byLeading0)
	}
}

func (u Uint128) QuoRem__0(by uint64) (q, r Uint128) {
	if u.hi < by {
		q.lo, r.lo = bits.Div64(u.hi, u.lo, by)
	} else {
		q.hi, r.lo = bits.Div64(0, u.hi, by)
		q.lo, r.lo = bits.Div64(r.lo, u.lo, by)
	}
	return q, r
}

// Gop_Rem: func (a uint128) % (b uint128) uint128
func (u Uint128) Gop_Rem__1(by Uint128) (r Uint128) {
	// TODO: inline only the needed bits
	_, r = u.QuoRem__1(by)
	return r
}

func (u Uint128) Gop_Rem__0(by uint64) (r Uint128) {
	// https://github.com/golang/go/issues/28970
	// if u.hi < by {
	//     _, r.lo = bits.Rem64(u.hi, u.lo, by)
	// } else {
	//     _, r.lo = bits.Rem64(bits.Rem64(0, u.hi, by), u.lo, by)
	// }

	if u.hi < by {
		_, r.lo = bits.Div64(u.hi, u.lo, by)
	} else {
		_, r.lo = bits.Div64(0, u.hi, by)
		_, r.lo = bits.Div64(r.lo, u.lo, by)
	}
	return r
}

// Hacker's delight 9-4, divlu:
func quo128by64(u1, u0, v uint64, vLeading0 uint) (q uint64) {
	var b uint64 = 1 << 32
	var un1, un0, vn1, vn0, q1, q0, un32, un21, un10, rhat, vs, left, right uint64

	vs = v << vLeading0

	vn1 = vs >> 32
	vn0 = vs & 0xffffffff

	if vLeading0 > 0 {
		un32 = (u1 << vLeading0) | (u0 >> (64 - vLeading0))
		un10 = u0 << vLeading0
	} else {
		un32 = u1
		un10 = u0
	}

	un1 = un10 >> 32
	un0 = un10 & 0xffffffff

	q1 = un32 / vn1
	rhat = un32 % vn1

	left = q1 * vn0
	right = (rhat << 32) | un1

again1:
	if (q1 >= b) || (left > right) {
		q1--
		rhat += vn1
		if rhat < b {
			left -= vn0
			right = (rhat << 32) | un1
			goto again1
		}
	}

	un21 = (un32 << 32) + (un1 - (q1 * vs))

	q0 = un21 / vn1
	rhat = un21 % vn1

	left = q0 * vn0
	right = (rhat << 32) | un0

again2:
	if (q0 >= b) || (left > right) {
		q0--
		rhat += vn1
		if rhat < b {
			left -= vn0
			right = (rhat << 32) | un0
			goto again2
		}
	}

	return (q1 << 32) | q0
}

// Hacker's delight 9-4, divlu:
func quorem128by64(u1, u0, v uint64, vLeading0 uint) (q, r uint64) {
	var b uint64 = 1 << 32
	var un1, un0, vn1, vn0, q1, q0, un32, un21, un10, rhat, left, right uint64

	v <<= vLeading0

	vn1 = v >> 32
	vn0 = v & 0xffffffff

	if vLeading0 > 0 {
		un32 = (u1 << vLeading0) | (u0 >> (64 - vLeading0))
		un10 = u0 << vLeading0
	} else {
		un32 = u1
		un10 = u0
	}

	un1 = un10 >> 32
	un0 = un10 & 0xffffffff

	q1 = un32 / vn1
	rhat = un32 % vn1

	left = q1 * vn0
	right = (rhat << 32) + un1

again1:
	if (q1 >= b) || (left > right) {
		q1--
		rhat += vn1
		if rhat < b {
			left -= vn0
			right = (rhat << 32) | un1
			goto again1
		}
	}

	un21 = (un32 << 32) + (un1 - (q1 * v))

	q0 = un21 / vn1
	rhat = un21 % vn1

	left = q0 * vn0
	right = (rhat << 32) | un0

again2:
	if (q0 >= b) || (left > right) {
		q0--
		rhat += vn1
		if rhat < b {
			left -= vn0
			right = (rhat << 32) | un0
			goto again2
		}
	}

	return (q1 << 32) | q0, ((un21 << 32) + (un0 - (q0 * v))) >> vLeading0
}

func quorem128by128(m, v Uint128, vHiLeading0, vLoLeading0 uint) (q, r Uint128) {
	if v.hi == 0 {
		if m.hi < v.lo {
			q.lo, r.lo = quorem128by64(m.hi, m.lo, v.lo, vLoLeading0)
			return q, r

		} else {
			q.hi = m.hi / v.lo
			r.hi = m.hi % v.lo
			q.lo, r.lo = quorem128by64(r.hi, m.lo, v.lo, vLoLeading0)
			r.hi = 0
			return q, r
		}

	} else {
		v1 := v.Gop_Lsh(vHiLeading0)
		u1 := m.Gop_Rsh(1)

		var q1 Uint128
		q1.lo = quo128by64(u1.hi, u1.lo, v1.hi, vLoLeading0)
		q1 = q1.Gop_Rsh(63 - vHiLeading0)

		if q1.hi|q1.lo != 0 {
			q1.Gop_Dec()
		}
		q = q1
		q1 = q1.Gop_Mul__1(v)
		r = m.Gop_Sub__1(q1)

		if r.Cmp__1(v) >= 0 {
			q.Gop_Inc()
			r = r.Gop_Sub__1(v)
		}

		return q, r
	}
}

func quorem128bin(u, by Uint128, uLeading0, byLeading0 uint) (q, r Uint128) {
	shift := int(byLeading0 - uLeading0)
	by = by.Gop_Lsh(uint(shift))

	for {
		// q << 1
		q.hi = (q.hi << 1) | (q.lo >> 63)
		q.lo = q.lo << 1

		// performance tweak: simulate greater than or equal by hand-inlining "not less than".
		if u.hi > by.hi || (u.hi == by.hi && u.lo >= by.lo) {
			u = u.Gop_Sub__1(by)
			q.lo |= 1
		}

		// by >> 1
		by.lo = (by.lo >> 1) | (by.hi << 63)
		by.hi = by.hi >> 1

		if shift <= 0 {
			break
		}
		shift--
	}

	r = u
	return q, r
}

func quo128bin(u, by Uint128, uLeading0, byLeading0 uint) (q Uint128) {
	shift := int(byLeading0 - uLeading0)
	by = by.Gop_Lsh(uint(shift))

	for {
		// q << 1
		q.hi = (q.hi << 1) | (q.lo >> 63)
		q.lo = q.lo << 1

		// u >= by
		if u.hi > by.hi || (u.hi == by.hi && u.lo >= by.lo) {
			u = u.Gop_Sub__1(by)
			q.lo |= 1
		}

		// q >> 1
		by.lo = (by.lo >> 1) | (by.hi << 63)
		by.hi = by.hi >> 1

		if shift <= 0 {
			break
		}
		shift--
	}

	return q
}

// -----------------------------------------------------------------------------

func ParseUint128(s string, base int) (out Uint128, err error) {
	b, ok := new(big.Int).SetString(s, base)
	if !ok {
		err = fmt.Errorf("invalid uint128 string: %q", s)
		return
	}
	out, inRange := Uint128_Cast__9(b)
	if !inRange {
		err = fmt.Errorf("string %q was not in valid uint128 range", s)
	}
	return
}

func FormatUint128(i Uint128, base int) string {
	return i.Text(base)
}

// -----------------------------------------------------------------------------
