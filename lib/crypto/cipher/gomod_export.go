// Package cipher provide Go+ "crypto/cipher" package, as "crypto/cipher" package in Go.
package cipher

import (
	cipher "crypto/cipher"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiAEADNonceSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cipher.AEAD).NonceSize()
	p.Ret(1, ret0)
}

func execiAEADOpen(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := args[0].(cipher.AEAD).Open(args[1].([]byte), args[2].([]byte), args[3].([]byte), args[4].([]byte))
	p.Ret(5, ret0, ret1)
}

func execiAEADOverhead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cipher.AEAD).Overhead()
	p.Ret(1, ret0)
}

func execiAEADSeal(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(cipher.AEAD).Seal(args[1].([]byte), args[2].([]byte), args[3].([]byte), args[4].([]byte))
	p.Ret(5, ret0)
}

func execiBlockBlockSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cipher.Block).BlockSize()
	p.Ret(1, ret0)
}

func execiBlockDecrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(cipher.Block).Decrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execiBlockEncrypt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(cipher.Block).Encrypt(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execiBlockModeBlockSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cipher.BlockMode).BlockSize()
	p.Ret(1, ret0)
}

func execiBlockModeCryptBlocks(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(cipher.BlockMode).CryptBlocks(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func toType0(v interface{}) cipher.Block {
	if v == nil {
		return nil
	}
	return v.(cipher.Block)
}

func execNewCBCDecrypter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewCBCDecrypter(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewCBCEncrypter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewCBCEncrypter(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewCFBDecrypter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewCFBDecrypter(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewCFBEncrypter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewCFBEncrypter(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewCTR(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewCTR(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewGCM(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := cipher.NewGCM(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execNewGCMWithNonceSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := cipher.NewGCMWithNonceSize(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execNewGCMWithTagSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := cipher.NewGCMWithTagSize(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execNewOFB(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := cipher.NewOFB(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execiStreamXORKeyStream(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(cipher.Stream).XORKeyStream(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmStreamReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(cipher.StreamReader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmStreamWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(cipher.StreamWriter).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmStreamWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cipher.StreamWriter).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/cipher")

func init() {
	I.RegisterFuncs(
		I.Func("(AEAD).NonceSize", (cipher.AEAD).NonceSize, execiAEADNonceSize),
		I.Func("(AEAD).Open", (cipher.AEAD).Open, execiAEADOpen),
		I.Func("(AEAD).Overhead", (cipher.AEAD).Overhead, execiAEADOverhead),
		I.Func("(AEAD).Seal", (cipher.AEAD).Seal, execiAEADSeal),
		I.Func("(Block).BlockSize", (cipher.Block).BlockSize, execiBlockBlockSize),
		I.Func("(Block).Decrypt", (cipher.Block).Decrypt, execiBlockDecrypt),
		I.Func("(Block).Encrypt", (cipher.Block).Encrypt, execiBlockEncrypt),
		I.Func("(BlockMode).BlockSize", (cipher.BlockMode).BlockSize, execiBlockModeBlockSize),
		I.Func("(BlockMode).CryptBlocks", (cipher.BlockMode).CryptBlocks, execiBlockModeCryptBlocks),
		I.Func("NewCBCDecrypter", cipher.NewCBCDecrypter, execNewCBCDecrypter),
		I.Func("NewCBCEncrypter", cipher.NewCBCEncrypter, execNewCBCEncrypter),
		I.Func("NewCFBDecrypter", cipher.NewCFBDecrypter, execNewCFBDecrypter),
		I.Func("NewCFBEncrypter", cipher.NewCFBEncrypter, execNewCFBEncrypter),
		I.Func("NewCTR", cipher.NewCTR, execNewCTR),
		I.Func("NewGCM", cipher.NewGCM, execNewGCM),
		I.Func("NewGCMWithNonceSize", cipher.NewGCMWithNonceSize, execNewGCMWithNonceSize),
		I.Func("NewGCMWithTagSize", cipher.NewGCMWithTagSize, execNewGCMWithTagSize),
		I.Func("NewOFB", cipher.NewOFB, execNewOFB),
		I.Func("(Stream).XORKeyStream", (cipher.Stream).XORKeyStream, execiStreamXORKeyStream),
		I.Func("(StreamReader).Read", (cipher.StreamReader).Read, execmStreamReaderRead),
		I.Func("(StreamWriter).Write", (cipher.StreamWriter).Write, execmStreamWriterWrite),
		I.Func("(StreamWriter).Close", (cipher.StreamWriter).Close, execmStreamWriterClose),
	)
	I.RegisterTypes(
		I.Type("AEAD", reflect.TypeOf((*cipher.AEAD)(nil)).Elem()),
		I.Type("Block", reflect.TypeOf((*cipher.Block)(nil)).Elem()),
		I.Type("BlockMode", reflect.TypeOf((*cipher.BlockMode)(nil)).Elem()),
		I.Type("Stream", reflect.TypeOf((*cipher.Stream)(nil)).Elem()),
		I.Type("StreamReader", reflect.TypeOf((*cipher.StreamReader)(nil)).Elem()),
		I.Type("StreamWriter", reflect.TypeOf((*cipher.StreamWriter)(nil)).Elem()),
	)
}
