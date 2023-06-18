package ng

import (
	"fmt"
	"log"
	"math/big"
	"math/bits"
)

const (
	Int128_Max       = 1<<127 - 1
	Int128_Min       = -1 << 127
	Int128_IsUntyped = true
)

const (
	signBit  = 0x8000000000000000
	maxInt64 = 1<<63 - 1
)

var (
	big1       = new(big.Int).SetUint64(1)
	bigMaxU128 *big.Int
)

func init() {
	bigMaxU128 = new(big.Int).Lsh(big1, 128)
	bigMaxU128.Sub(bigMaxU128, big1) // 340282366920938463463374607431768211455
}

// -----------------------------------------------------------------------------

type Int128 struct {
	hi uint64
	lo uint64
}

// Int128_Init: func int128.init(v int) int128
func Int128_Init__0(v int) (out Int128) {
	return Int128_Cast__2(int64(v))
}

// Int128_Init: func int128.init(v untyped_bigint) int128
func Int128_Init__1(v UntypedBigint) (out Int128) {
	return Int128_Cast__1(v)
}

// Int128_Cast: func int128(v int) int128
func Int128_Cast__0(v int) Int128 {
	return Int128_Cast__2(int64(v))
}

// Int128_Cast: func int128(v untyped_bigint) int128
func Int128_Cast__1(v UntypedBigint) (out Int128) {
	out, inRange := Int128_Cast__9(v)
	if !inRange {
		log.Panicf("value %v was not in valid int128 range\n", v)
	}
	return
}

// Int128_Cast: func int128(v int64) int128
func Int128_Cast__2(v int64) (out Int128) {
	var hi uint64
	if v < 0 {
		hi = maxUint64
	}
	return Int128{hi: hi, lo: uint64(v)}
}

// Int128_Cast: func int128(v uint64) int128
func Int128_Cast__3(v uint64) Int128 {
	return Int128{lo: v}
}

// Int128_Cast: func int128(v int32) int128
func Int128_Cast__4(v int32) Int128 {
	return Int128_Cast__2(int64(v))
}

// Int128_Cast: func int128(v int16) int128
func Int128_Cast__5(v int16) Int128 {
	return Int128_Cast__2(int64(v))
}

// Int128_Cast: func int128(v int8) int128
func Int128_Cast__6(v int8) Int128 {
	return Int128_Cast__2(int64(v))
}

// Int128_Cast: func int128(v uint18) int128
func Int128_Cast__7(v Uint128) (out Int128) {
	return Int128{hi: v.hi, lo: v.lo}
}

// Int128_Cast: func int128(v *big.Int) int128
func Int128_Cast__8(v *big.Int) Int128 {
	out, _ := Int128_Cast__9(v)
	return out
}

func Int128_Cast__9(v *big.Int) (out Int128, inRange bool) {
	neg := v.Sign() < 0
	words := v.Bits()

	var u Uint128
	inRange = true
	switch intSize {
	case 64:
		lw := len(words)
		switch lw {
		case 0:
		case 1:
			u.lo = uint64(words[0])
		case 2:
			u.hi = uint64(words[1])
			u.lo = uint64(words[0])
		default:
			u, inRange = Uint128{hi: maxUint64, lo: maxUint64}, false
		}

	case 32:
		lw := len(words)
		switch lw {
		case 0:
		case 1:
			u.lo = uint64(words[0])
		case 2:
			u.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		case 3:
			u.hi = uint64(words[2])
			u.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		case 4:
			u.hi = (uint64(words[3]) << 32) | (uint64(words[2]))
			u.lo = (uint64(words[1]) << 32) | (uint64(words[0]))
		default:
			u, inRange = Uint128{hi: maxUint64, lo: maxUint64}, false
		}

	default:
		panic("unsupported bit size")
	}

	if neg {
		if cmp := u.Cmp__1(Uint128{hi: 0x8000000000000000, lo: 0}); cmp > 0 {
			out, inRange = Int128{hi: 0x8000000000000000, lo: 0}, false
		} else {
			out = Int128{hi: u.hi, lo: u.lo}.Gop_Neg()
		}
	} else {
		if cmp := u.Cmp__1(Uint128{hi: maxInt64, lo: maxUint64}); cmp > 0 {
			out, inRange = Int128{hi: maxInt64, lo: maxUint64}, false
		} else {
			out = Int128{hi: u.hi, lo: u.lo}
		}
	}
	return
}

// Int128_Cast: func int128() int128
func Int128_Cast__a() Int128 {
	return Int128{}
}

// Gop_Rcast: func uint128(v int128) uint128
func (i Int128) Gop_Rcast__0() Uint128 {
	return Uint128{lo: i.lo, hi: i.hi}
}

// Gop_Rcast: func uint128(v int128) (uint128, bool)
func (i Int128) Gop_Rcast__1() (out Uint128, inRange bool) {
	return Uint128{lo: i.lo, hi: i.hi}, i.hi&signBit == 0
}

// Gop_Rcast: func int64(v int128) int64
func (i Int128) Gop_Rcast__2() int64 {
	if i.hi&signBit == 0 {
		return int64(i.lo)
	}
	return -int64(^(i.lo - 1))
}

// Gop_Rcast: func int64(v int128) (int64, bool)
func (i Int128) Gop_Rcast__3() (out int64, inRange bool) {
	if i.hi&signBit == 0 {
		return int64(i.lo), i.hi == 0 && i.lo <= maxInt64
	}
	return -int64(^(i.lo - 1)), i.hi == maxUint64 && i.lo >= 0x8000000000000000
}

// Gop_Rcast: func uint64(v int128) uint64
func (i Int128) Gop_Rcast__4() uint64 {
	return i.lo
}

// Gop_Rcast: func uint64(v int128) (uint64, bool)
func (i Int128) Gop_Rcast__5() (out uint64, inRange bool) {
	return i.lo, i.hi == 0
}

// -----------------------------------------------------------------------------

func (i Int128) IsZero() bool {
	return i.lo == 0 && i.hi == 0
}

func (i *Int128) Scan(state fmt.ScanState, verb rune) (err error) {
	t, err := state.Token(true, nil)
	if err != nil {
		return
	}
	v, err := ParseInt128(string(t), 10)
	if err == nil {
		*i = v
	}
	return
}

func (i Int128) Format(s fmt.State, c rune) {
	// TODO: not so good
	i.BigInt().Format(s, c)
}

func (i Int128) String() string {
	return i.Text(10)
}

func (i Int128) Text(base int) string {
	// TODO: not so good
	return i.BigInt().Text(base)
}

func (i Int128) BigInt() *big.Int {
	var v big.Int
	i.ToBigInt(&v)
	return &v
}

func (i Int128) ToBigInt(b *big.Int) {
	neg := i.hi&signBit != 0
	if i.hi > 0 {
		b.SetUint64(i.hi)
		b.Lsh(b, 64)
	}
	var lo big.Int
	lo.SetUint64(i.lo)
	b.Add(b, &lo)

	if neg {
		b.Xor(b, bigMaxU128).Add(b, big1).Neg(b)
	}
}

func (i Int128) Sign() int {
	if i.lo == 0 && i.hi == 0 {
		return 0
	} else if i.hi&signBit == 0 {
		return 1
	}
	return -1
}

func (i *Int128) Gop_Inc() {
	i.lo++
	if i.lo == 0 {
		i.hi++
	}
}

func (i *Int128) Gop_Dec() {
	if i.lo == 0 {
		i.hi--
	}
	i.lo--
}

// Gop_AddAssign: func (a *int128) += (b int128)
func (i *Int128) Gop_AddAssign(b Int128) {
	*i = i.Gop_Add__1(b)
}

// Gop_SubAssign: func (a *int128) -= (b int128)
func (i *Int128) Gop_SubAssign(b Int128) {
	*i = i.Gop_Sub__1(b)
}

// Gop_MulAssign: func (a *int128) *= (b int128)
func (i *Int128) Gop_MulAssign(b Int128) {
	*i = i.Gop_Mul__1(b)
}

// Gop_QuoAssign: func (a *int128) /= (b int128) {
func (i *Int128) Gop_QuoAssign(b Int128) {
	*i = i.Gop_Quo__1(b)
}

// Gop_RemAssign: func (a *int128) %= (b int128)
func (i *Int128) Gop_RemAssign(b Int128) {
	*i = i.Gop_Rem__1(b)
}

// Gop_OrAssign: func (a *int128) |= (b int128)
func (i *Int128) Gop_OrAssign(b Int128) {
	*i = i.Gop_Or(b)
}

// Gop_XorAssign: func (a *int128) ^= (b int128)
func (i *Int128) Gop_XorAssign(b Int128) {
	*i = i.Gop_Xor(b)
}

// Gop_AndAssign: func (a *int128) &= (b int128)
func (i *Int128) Gop_AndAssign(b Int128) {
	*i = i.Gop_And(b)
}

// Gop_AndNotAssign: func (a *int128) &^= (b int128)
func (i *Int128) Gop_AndNotAssign(b Int128) {
	*i = i.Gop_AndNot(b)
}

// Gop_LshAssign: func (a *int128) <<= (n untyped_uint)
func (i *Int128) Gop_LshAssign(n Gop_ninteger) {
	*i = i.Gop_Lsh(n)
}

// Gop_RshAssign: func (a *int128) >>= (n untyped_uint)
func (i *Int128) Gop_RshAssign(n Gop_ninteger) {
	*i = i.Gop_Rsh(n)
}

func (i Int128) Gop_Add__1(n Int128) (v Int128) {
	var carry uint64
	v.lo, carry = bits.Add64(i.lo, n.lo, 0)
	v.hi, _ = bits.Add64(i.hi, n.hi, carry)
	return v
}

func (i Int128) Gop_Add__0(n int64) (v Int128) {
	var carry uint64
	v.lo, carry = bits.Add64(i.lo, uint64(n), 0)
	if n < 0 {
		v.hi = i.hi + maxUint64 + carry
	} else {
		v.hi = i.hi + carry
	}
	return v
}

func (i Int128) Gop_Sub__1(n Int128) (v Int128) {
	var borrowed uint64
	v.lo, borrowed = bits.Sub64(i.lo, n.lo, 0)
	v.hi, _ = bits.Sub64(i.hi, n.hi, borrowed)
	return v
}

func (i Int128) Gop_Sub__0(n int64) (v Int128) {
	var borrowed uint64
	if n < 0 {
		v.lo, borrowed = bits.Sub64(i.lo, uint64(n), 0)
		v.hi = i.hi - maxUint64 - borrowed
	} else {
		v.lo, borrowed = bits.Sub64(i.lo, uint64(n), 0)
		v.hi = i.hi - borrowed
	}
	return v
}

func (i Int128) Gop_Neg() (v Int128) {
	if i.lo == 0 && i.hi == 0 {
		return
	}
	if i.hi&signBit == 0 {
		v.hi = ^i.hi
		v.lo = (^i.lo) + 1
	} else {
		v.hi = ^i.hi
		v.lo = ^(i.lo - 1)
	}
	if v.lo == 0 { // handle overflow
		v.hi++
	}
	return v
}

func (i Int128) Gop_Dup() (v Int128) {
	return i
}

// Abs returns the absolute value of i as a signed integer.
func (i Int128) Abs__0() Int128 {
	if i.hi&signBit != 0 {
		i.hi = ^i.hi
		i.lo = ^(i.lo - 1)
		if i.lo == 0 { // handle carry
			i.hi++
		}
	}
	return i
}

func (i Int128) Abs__1() (ret Int128, inRange bool) {
	return i.Abs__0(), i.hi != 0x8000000000000000 || i.lo != 0
}

// AbsU returns the absolute value of i as an unsigned integer. All
// values of i are representable using this function, but the type is
// changed.
func (i Int128) AbsU() Uint128 {
	if i.hi == 0x8000000000000000 && i.lo == 0 {
		return Uint128{hi: 0x8000000000000000}
	}
	if i.hi&signBit != 0 {
		i.hi = ^i.hi
		i.lo = ^(i.lo - 1)
		if i.lo == 0 { // handle carry
			i.hi++
		}
	}
	return Uint128{hi: i.hi, lo: i.lo}
}

// Cmp compares i to n and returns:
//
//	< 0 if i <  n
//	  0 if i == n
//	> 0 if i >  n
//
// The specific value returned by Cmp is undefined, but it is guaranteed to
// satisfy the above constraints.
func (i Int128) Cmp__1(n Int128) int {
	if i.hi == n.hi && i.lo == n.lo {
		return 0
	} else if i.hi&signBit == n.hi&signBit {
		if i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo) {
			return 1
		}
	} else if i.hi&signBit == 0 {
		return 1
	}
	return -1
}

// Cmp64 compares 'i' to 64-bit int 'n' and returns:
//
//	< 0 if i <  n
//	  0 if i == n
//	> 0 if i >  n
//
// The specific value returned by Cmp is undefined, but it is guaranteed to
// satisfy the above constraints.
func (i Int128) Cmp__0(n int64) int {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}
	if i.hi == nhi && i.lo == nlo {
		return 0
	} else if i.hi&signBit == nhi&signBit {
		if i.hi > nhi || (i.hi == nhi && i.lo > nlo) {
			return 1
		}
	} else if i.hi&signBit == 0 {
		return 1
	}
	return -1
}

func (i Int128) Gop_EQ__1(n Int128) bool {
	return i.hi == n.hi && i.lo == n.lo
}

func (i Int128) Gop_EQ__0(n int64) bool {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}
	return i.hi == nhi && i.lo == nlo
}

func (i Int128) Gop_GT__1(n Int128) bool {
	if i.hi&signBit == n.hi&signBit {
		return i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo)
	} else if i.hi&signBit == 0 {
		return true
	}
	return false
}

func (i Int128) Gop_GT__0(n int64) bool {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}

	if i.hi&signBit == nhi&signBit {
		return i.hi > nhi || (i.hi == nhi && i.lo > nlo)
	} else if i.hi&signBit == 0 {
		return true
	}
	return false
}

func (i Int128) Gop_GE__1(n Int128) bool {
	if i.hi == n.hi && i.lo == n.lo {
		return true
	}
	if i.hi&signBit == n.hi&signBit {
		return i.hi > n.hi || (i.hi == n.hi && i.lo > n.lo)
	} else if i.hi&signBit == 0 {
		return true
	}
	return false
}

func (i Int128) Gop_GE__0(n int64) bool {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}

	if i.hi == nhi && i.lo == nlo {
		return true
	}
	if i.hi&signBit == nhi&signBit {
		return i.hi > nhi || (i.hi == nhi && i.lo > nlo)
	} else if i.hi&signBit == 0 {
		return true
	}
	return false
}

func (i Int128) Gop_LT__1(n Int128) bool {
	if i.hi&signBit == n.hi&signBit {
		return i.hi < n.hi || (i.hi == n.hi && i.lo < n.lo)
	} else if i.hi&signBit != 0 {
		return true
	}
	return false
}

func (i Int128) Gop_LT__0(n int64) bool {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}

	if i.hi&signBit == nhi&signBit {
		return i.hi < nhi || (i.hi == nhi && i.lo < nlo)
	} else if i.hi&signBit != 0 {
		return true
	}
	return false
}

func (i Int128) Gop_LE__1(n Int128) bool {
	if i.hi == n.hi && i.lo == n.lo {
		return true
	}
	if i.hi&signBit == n.hi&signBit {
		return i.hi < n.hi || (i.hi == n.hi && i.lo < n.lo)
	} else if i.hi&signBit != 0 {
		return true
	}
	return false
}

func (i Int128) Gop_LE__0(n int64) bool {
	var nhi uint64
	var nlo = uint64(n)
	if n < 0 {
		nhi = maxUint64
	}

	if i.hi == nhi && i.lo == nlo {
		return true
	}
	if i.hi&signBit == nhi&signBit {
		return i.hi < nhi || (i.hi == nhi && i.lo < nlo)
	} else if i.hi&signBit != 0 {
		return true
	}
	return false
}

func (i Int128) Gop_And(n Int128) Int128 {
	i.hi &= n.hi
	i.lo &= n.lo
	return i
}

func (i Int128) Gop_AndNot(n Int128) Int128 {
	i.hi &^= n.hi
	i.lo &^= n.lo
	return i
}

func (i Int128) Gop_Not() Int128 {
	return Int128{hi: ^i.hi, lo: ^i.lo}
}

func (i Int128) Gop_Or(n Int128) Int128 {
	i.hi |= n.hi
	i.lo |= n.lo
	return i
}

func (i Int128) Gop_Xor(v Int128) Int128 {
	i.hi ^= v.hi
	i.lo ^= v.lo
	return i
}

func (i Int128) Gop_Lsh(n Gop_ninteger) Int128 {
	if n < 64 {
		i.hi = (i.hi << n) | (i.lo >> (64 - n))
		i.lo <<= n
	} else {
		i.hi = i.lo << (n - 64)
		i.lo = 0
	}
	return i
}

func (i Int128) Gop_Rsh(n Gop_ninteger) Int128 {
	if n < 64 {
		i.lo = (i.lo >> n) | (i.hi << (64 - n))
	} else {
		i.lo = i.hi >> (n - 64)
	}
	i.hi = uint64(int64(i.hi) >> n)
	return i
}

// Mul returns the product of two I128s.
//
// Overflow should wrap around, as per the Go spec.
func (i Int128) Gop_Mul__1(n Int128) (dest Int128) {
	hi, lo := bits.Mul64(i.lo, n.lo)
	hi += i.hi*n.lo + i.lo*n.hi
	return Int128{hi, lo}
}

func (i Int128) Gop_Mul__0(n int64) Int128 {
	nlo := uint64(n)
	var nhi uint64
	if n < 0 {
		nhi = maxUint64
	}
	hi, lo := bits.Mul64(i.lo, nlo)
	hi += i.hi*nlo + i.lo*nhi
	return Int128{hi, lo}
}

// QuoRem returns the quotient q and remainder r for y != 0. If y == 0, a
// division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//	q = x/y      with the result truncated to zero
//	r = x - y*q
//
// U128 does not support big.Int.DivMod()-style Euclidean division.
//
// Note: dividing MinI128 by -1 will overflow, returning MinI128, as
// per the Go spec (https://golang.org/ref/spec#Integer_operators):
//
//	The one exception to this rule is that if the dividend x is the most
//	negative value for the int type of x, the quotient q = x / -1 is equal to x
//	(and r = 0) due to two's-complement integer overflow.
func (i Int128) QuoRem__1(by Int128) (q, r Int128) {
	qSign, rSign := 1, 1
	if i.Gop_LT__0(0) {
		qSign, rSign = -1, -1
		i = i.Gop_Neg()
	}
	if by.Gop_LT__0(0) {
		qSign = -qSign
		by = by.Gop_Neg()
	}

	qu, ru := i.Gop_Rcast__0().QuoRem__1(by.Gop_Rcast__0())
	q, r = Int128_Cast__7(qu), Int128_Cast__7(ru)
	if qSign < 0 {
		q = q.Gop_Neg()
	}
	if rSign < 0 {
		r = r.Gop_Neg()
	}
	return q, r
}

func (i Int128) QuoRem__0(by int64) (q, r Int128) {
	ineg := i.hi&signBit != 0
	if ineg {
		i = i.Gop_Neg()
	}
	byneg := by < 0
	if byneg {
		by = -by
	}

	n := uint64(by)
	if i.hi < n {
		q.lo, r.lo = bits.Div64(i.hi, i.lo, n)
	} else {
		q.hi, r.lo = bits.Div64(0, i.hi, n)
		q.lo, r.lo = bits.Div64(r.lo, i.lo, n)
	}
	if ineg != byneg {
		q = q.Gop_Neg()
	}
	if ineg {
		r = r.Gop_Neg()
	}
	return q, r
}

// Quo returns the quotient x/y for y != 0. If y == 0, a division-by-zero
// run-time panic occurs. Quo implements truncated division (like Go); see
// QuoRem for more details.
func (i Int128) Gop_Quo__1(by Int128) (q Int128) {
	qSign := 1
	if i.Gop_LT__0(0) {
		qSign = -1
		i = i.Gop_Neg()
	}
	if by.Gop_LT__0(0) {
		qSign = -qSign
		by = by.Gop_Neg()
	}

	qu := i.Gop_Rcast__0().Gop_Quo__1(by.Gop_Rcast__0())
	q = Int128_Cast__7(qu)
	if qSign < 0 {
		q = q.Gop_Neg()
	}
	return q
}

func (i Int128) Gop_Quo__0(by int64) (q Int128) {
	ineg := i.hi&signBit != 0
	if ineg {
		i = i.Gop_Neg()
	}
	byneg := by < 0
	if byneg {
		by = -by
	}

	n := uint64(by)
	if i.hi < n {
		q.lo, _ = bits.Div64(i.hi, i.lo, n)
	} else {
		var rlo uint64
		q.hi, rlo = bits.Div64(0, i.hi, n)
		q.lo, _ = bits.Div64(rlo, i.lo, n)
	}
	if ineg != byneg {
		q = q.Gop_Neg()
	}
	return q
}

// Gop_Rem returns the remainder of x%y for y != 0. If y == 0, a division-by-zero
// run-time panic occurs. Gop_Rem implements truncated modulus (like Go); see
// QuoRem for more details.
func (i Int128) Gop_Rem__1(by Int128) (r Int128) {
	_, r = i.QuoRem__1(by)
	return r
}

func (i Int128) Gop_Rem__0(by int64) (r Int128) {
	ineg := i.hi&signBit != 0
	if ineg {
		i = i.Gop_Neg()
	}
	if by < 0 {
		by = -by
	}

	n := uint64(by)
	if i.hi < n {
		_, r.lo = bits.Div64(i.hi, i.lo, n)
	} else {
		_, r.lo = bits.Div64(0, i.hi, n)
		_, r.lo = bits.Div64(r.lo, i.lo, n)
	}
	if ineg {
		r = r.Gop_Neg()
	}
	return r
}

// -----------------------------------------------------------------------------

func ParseInt128(s string, base int) (out Int128, err error) {
	b, ok := new(big.Int).SetString(s, base)
	if !ok {
		err = fmt.Errorf("invalid int128 string: %q", s)
		return
	}
	out, inRange := Int128_Cast__9(b)
	if !inRange {
		err = fmt.Errorf("string %q was not in valid int128 range", s)
	}
	return
}

func FormatInt128(i Int128, base int) string {
	return i.Text(base)
}

// -----------------------------------------------------------------------------
