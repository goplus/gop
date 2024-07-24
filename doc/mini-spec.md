Go+ STEM Education Minimum Specification
=====

## Comments

Comments serve as program documentation. There are three forms:

* __Line comments__ start with the character sequence `//` and stop at the end of the line.
* __Line comments__ start with the character sequence `#` and stop at the end of the line.
* __General comments__ start with the character sequence `/*` and stop with the first subsequent character sequence `*/`.

A __general comment__ containing no newlines acts like a space. Any other comment acts like a newline.

```
# this is a line comment
// this is another line comment
/* this is  a general comment */
```

## Literals

### Integer literals

An integer literal is a sequence of digits representing an [integer constant](). An optional prefix sets a non-decimal base: 0b or 0B for binary, 0, 0o, or 0O for octal, and 0x or 0X for hexadecimal. A single 0 is considered a decimal zero. In hexadecimal literals, letters a through f and A through F represent values 10 through 15.

For readability, an underscore character _ may appear after a base prefix or between successive digits; such underscores do not change the literal's value.

```go
42
4_2
0600
0_600
0o600
0O600       // second character is capital letter 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
170141183460469231731687303715884105727
170_141183_460469_231731_687303_715884_105727

_42         // an identifier, not an integer literal
42_         // invalid: _ must separate successive digits
4__2        // invalid: only one _ at a time
0_xBadFace  // invalid: _ must separate successive digits
```

### Floating-point literals

A floating-point literal is a decimal or hexadecimal representation of a [floating-point constant]().

A decimal floating-point literal consists of an integer part (decimal digits), a decimal point, a fractional part (decimal digits), and an exponent part (e or E followed by an optional sign and decimal digits). One of the integer part or the fractional part may be elided; one of the decimal point or the exponent part may be elided. An exponent value exp scales the mantissa (integer and fractional part) by 10<sup>exp</sup>.

A hexadecimal floating-point literal consists of a 0x or 0X prefix, an integer part (hexadecimal digits), a radix point, a fractional part (hexadecimal digits), and an exponent part (p or P followed by an optional sign and decimal digits). One of the integer part or the fractional part may be elided; the radix point may be elided as well, but the exponent part is required. (This syntax matches the one given in IEEE 754-2008 §5.12.3.) An exponent value exp scales the mantissa (integer and fractional part) by 2<sup>exp</sup>.

For readability, an underscore character _ may appear after a base prefix or between successive digits; such underscores do not change the literal value.

```go
0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375

0x15e-2      // == 0x15e - 2 (integer subtraction)

0x.p1        // invalid: mantissa has no digits
1p-2         // invalid: p exponent requires hexadecimal mantissa
0x1.5e-2     // invalid: hexadecimal mantissa requires p exponent
1_.5         // invalid: _ must separate successive digits
1._5         // invalid: _ must separate successive digits
1.5_e1       // invalid: _ must separate successive digits
1.5e_1       // invalid: _ must separate successive digits
1.5e1_       // invalid: _ must separate successive digits
```

### Imaginary literals

An imaginary literal represents the imaginary part of a [complex constant](). It consists of an [integer](#integer-literals) or [floating-point](#floating-point-literals) literal followed by the lowercase letter __i__. The value of an imaginary literal is the value of the respective integer or floating-point literal multiplied by the imaginary unit __i__.

For backward compatibility, an imaginary literal's integer part consisting entirely of decimal digits (and possibly underscores) is considered a decimal integer, even if it starts with a leading 0.

```go
0i
0123i         // == 123i for backward-compatibility
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i
```

### Rational literals

TODO


### Rune literals

A rune literal represents a [rune constant](), an integer value identifying a Unicode code point. A rune literal is expressed as one or more characters enclosed in single quotes, as in `'x'` or `'\n'`. Within the quotes, any character may appear except newline and unescaped single quote. A single quoted character represents the Unicode value of the character itself, while multi-character sequences beginning with a backslash encode values in various formats.

The simplest form represents the single character within the quotes; since Go+ source text is Unicode characters encoded in UTF-8, multiple UTF-8-encoded bytes may represent a single integer value. For instance, the literal `'a'` holds a single byte representing a literal a, Unicode `U+0061`, value 0x61, while `'ä'` holds two bytes (0xc3 0xa4) representing a literal a-dieresis, `U+00E4`, value 0xe4.

Several backslash escapes allow arbitrary values to be encoded as ASCII text. There are four ways to represent the integer value as a numeric constant: `\x` followed by exactly two hexadecimal digits; `\u` followed by exactly four hexadecimal digits; `\U` followed by exactly eight hexadecimal digits, and a plain backslash `\` followed by exactly three octal digits. In each case the value of the literal is the value represented by the digits in the corresponding base.

Although these representations all result in an integer, they have different valid ranges. Octal escapes must represent a value between 0 and 255 inclusive. Hexadecimal escapes satisfy this condition by construction. The escapes `\u` and `\U` represent Unicode code points so within them some values are illegal, in particular those above 0x10FFFF and surrogate halves.

After a backslash, certain single-character escapes represent special values:

```
\a   U+0007 alert or bell
\b   U+0008 backspace
\f   U+000C form feed
\n   U+000A line feed or newline
\r   U+000D carriage return
\t   U+0009 horizontal tab
\v   U+000B vertical tab
\\   U+005C backslash
\'   U+0027 single quote  (valid escape only within rune literals)
\"   U+0022 double quote  (valid escape only within string literals)
```

An unrecognized character following a backslash in a rune literal is illegal.

```go
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
'\''         // rune literal containing single quote character
'aa'         // illegal: too many characters
'\k'         // illegal: k is not recognized after a backslash
'\xa'        // illegal: too few hexadecimal digits
'\0'         // illegal: too few octal digits
'\400'       // illegal: octal value over 255
'\uDFFF'     // illegal: surrogate half
'\U00110000' // illegal: invalid Unicode code point
```

### String literals

A string literal represents a [string constant]() obtained from concatenating a sequence of characters. There are two forms: raw string literals and interpreted string literals.

Raw string literals are character sequences between back quotes, as in \`foo\`. Within the quotes, any character may appear except back quote. The value of a raw string literal is the string composed of the uninterpreted (implicitly UTF-8-encoded) characters between the quotes; in particular, backslashes have no special meaning and the string may contain newlines. Carriage return characters (`'\r'`) inside raw string literals are discarded from the raw string value.

Interpreted string literals are character sequences between double quotes, as in `"bar"`. Within the quotes, any character may appear except newline and unescaped double quote. The text between the quotes forms the value of the literal, with backslash escapes interpreted as they are in [rune literals](#rune-literals) (except that `\'` is illegal and `\"` is legal), with the same restrictions. The three-digit octal (`\nnn`) and two-digit hexadecimal (`\xnn`) escapes represent individual bytes of the resulting string; all other escapes represent the (possibly multi-byte) UTF-8 encoding of individual characters. Thus inside a string literal `\377` and `\xFF` represent a single byte of value 0xFF=255, while `ÿ`, `\u00FF`, `\U000000FF` and `\xc3\xbf` represent the two bytes 0xc3 0xbf of the UTF-8 encoding of character `U+00FF`.

```go
`abc`                // same as "abc"
`\n
\n`                  // same as "\\n\n\\n"
"\n"
"\""                 // same as `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
"\uD800"             // illegal: surrogate half
"\U00110000"         // illegal: invalid Unicode code point
```

These examples all represent the same string:

```go
"日本語"                                 // UTF-8 input text
`日本語`                                 // UTF-8 input text as a raw literal
"\u65e5\u672c\u8a9e"                    // the explicit Unicode code points
"\U000065e5\U0000672c\U00008a9e"        // the explicit Unicode code points
"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // the explicit UTF-8 bytes
```

If the source code represents a character as two code points, such as a combining form involving an accent and a letter, the result will be an error if placed in a rune literal (it is not a single code point), and will appear as two code points if placed in a string literal.

#### C style string literals

TODO

```go
c"Hello, world!\n"
```

#### Python string literals

TODO

```go
py"Hello, world!\n"
```
