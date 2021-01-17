// Package os provide Go+ "os" package, as "os" package in Go.
package os

import (
	os "os"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execChdir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.Chdir(args[0].(string))
	p.Ret(1, ret0)
}

func execChmod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Chmod(args[0].(string), args[1].(os.FileMode))
	p.Ret(2, ret0)
}

func execChown(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := os.Chown(args[0].(string), args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execChtimes(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := os.Chtimes(args[0].(string), args[1].(time.Time), args[2].(time.Time))
	p.Ret(3, ret0)
}

func execClearenv(_ int, p *gop.Context) {
	os.Clearenv()
	p.PopN(0)
}

func execCreate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.Create(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execEnviron(_ int, p *gop.Context) {
	ret0 := os.Environ()
	p.Ret(0, ret0)
}

func execExecutable(_ int, p *gop.Context) {
	ret0, ret1 := os.Executable()
	p.Ret(0, ret0, ret1)
}

func execExit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	os.Exit(args[0].(int))
	p.PopN(1)
}

func execExpand(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Expand(args[0].(string), args[1].(func(string) string))
	p.Ret(2, ret0)
}

func execExpandEnv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.ExpandEnv(args[0].(string))
	p.Ret(1, ret0)
}

func execmFileReaddir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*os.File).Readdir(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmFileReaddirnames(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*os.File).Readdirnames(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmFileName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.File).Name()
	p.Ret(1, ret0)
}

func execmFileRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*os.File).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmFileReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*os.File).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execmFileWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*os.File).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmFileWriteAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*os.File).WriteAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execmFileSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*os.File).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmFileWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*os.File).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmFileChmod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.File).Chmod(args[1].(os.FileMode))
	p.Ret(2, ret0)
}

func execmFileSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.File).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmFileSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.File).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmFileSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.File).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmFileSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*os.File).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmFileChown(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*os.File).Chown(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmFileTruncate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.File).Truncate(args[1].(int64))
	p.Ret(2, ret0)
}

func execmFileSync(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.File).Sync()
	p.Ret(1, ret0)
}

func execmFileChdir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.File).Chdir()
	p.Ret(1, ret0)
}

func execmFileFd(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.File).Fd()
	p.Ret(1, ret0)
}

func execmFileClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.File).Close()
	p.Ret(1, ret0)
}

func execmFileStat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*os.File).Stat()
	p.Ret(1, ret0, ret1)
}

func execiFileInfoIsDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).IsDir()
	p.Ret(1, ret0)
}

func execiFileInfoModTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).ModTime()
	p.Ret(1, ret0)
}

func execiFileInfoMode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).Mode()
	p.Ret(1, ret0)
}

func execiFileInfoName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).Name()
	p.Ret(1, ret0)
}

func execiFileInfoSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).Size()
	p.Ret(1, ret0)
}

func execiFileInfoSys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileInfo).Sys()
	p.Ret(1, ret0)
}

func execmFileModeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileMode).String()
	p.Ret(1, ret0)
}

func execmFileModeIsDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileMode).IsDir()
	p.Ret(1, ret0)
}

func execmFileModeIsRegular(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileMode).IsRegular()
	p.Ret(1, ret0)
}

func execmFileModePerm(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.FileMode).Perm()
	p.Ret(1, ret0)
}

func execFindProcess(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.FindProcess(args[0].(int))
	p.Ret(1, ret0, ret1)
}

func execGetegid(_ int, p *gop.Context) {
	ret0 := os.Getegid()
	p.Ret(0, ret0)
}

func execGetenv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.Getenv(args[0].(string))
	p.Ret(1, ret0)
}

func execGeteuid(_ int, p *gop.Context) {
	ret0 := os.Geteuid()
	p.Ret(0, ret0)
}

func execGetgid(_ int, p *gop.Context) {
	ret0 := os.Getgid()
	p.Ret(0, ret0)
}

func execGetgroups(_ int, p *gop.Context) {
	ret0, ret1 := os.Getgroups()
	p.Ret(0, ret0, ret1)
}

func execGetpagesize(_ int, p *gop.Context) {
	ret0 := os.Getpagesize()
	p.Ret(0, ret0)
}

func execGetpid(_ int, p *gop.Context) {
	ret0 := os.Getpid()
	p.Ret(0, ret0)
}

func execGetppid(_ int, p *gop.Context) {
	ret0 := os.Getppid()
	p.Ret(0, ret0)
}

func execGetuid(_ int, p *gop.Context) {
	ret0 := os.Getuid()
	p.Ret(0, ret0)
}

func execGetwd(_ int, p *gop.Context) {
	ret0, ret1 := os.Getwd()
	p.Ret(0, ret0, ret1)
}

func execHostname(_ int, p *gop.Context) {
	ret0, ret1 := os.Hostname()
	p.Ret(0, ret0, ret1)
}

func toType0(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}

func execIsExist(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.IsExist(toType0(args[0]))
	p.Ret(1, ret0)
}

func execIsNotExist(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.IsNotExist(toType0(args[0]))
	p.Ret(1, ret0)
}

func execIsPathSeparator(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.IsPathSeparator(args[0].(uint8))
	p.Ret(1, ret0)
}

func execIsPermission(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.IsPermission(toType0(args[0]))
	p.Ret(1, ret0)
}

func execIsTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.IsTimeout(toType0(args[0]))
	p.Ret(1, ret0)
}

func execLchown(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := os.Lchown(args[0].(string), args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execLink(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Link(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execmLinkErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.LinkError).Error()
	p.Ret(1, ret0)
}

func execmLinkErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.LinkError).Unwrap()
	p.Ret(1, ret0)
}

func execLookupEnv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.LookupEnv(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLstat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.Lstat(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execMkdir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Mkdir(args[0].(string), args[1].(os.FileMode))
	p.Ret(2, ret0)
}

func execMkdirAll(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.MkdirAll(args[0].(string), args[1].(os.FileMode))
	p.Ret(2, ret0)
}

func execNewFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.NewFile(args[0].(uintptr), args[1].(string))
	p.Ret(2, ret0)
}

func execNewSyscallError(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.NewSyscallError(args[0].(string), toType0(args[1]))
	p.Ret(2, ret0)
}

func execOpen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.Open(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execOpenFile(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := os.OpenFile(args[0].(string), args[1].(int), args[2].(os.FileMode))
	p.Ret(3, ret0, ret1)
}

func execmPathErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.PathError).Error()
	p.Ret(1, ret0)
}

func execmPathErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.PathError).Unwrap()
	p.Ret(1, ret0)
}

func execmPathErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.PathError).Timeout()
	p.Ret(1, ret0)
}

func execPipe(_ int, p *gop.Context) {
	ret0, ret1, ret2 := os.Pipe()
	p.Ret(0, ret0, ret1, ret2)
}

func execmProcessRelease(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.Process).Release()
	p.Ret(1, ret0)
}

func execmProcessKill(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.Process).Kill()
	p.Ret(1, ret0)
}

func execmProcessWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*os.Process).Wait()
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) os.Signal {
	if v == nil {
		return nil
	}
	return v.(os.Signal)
}

func execmProcessSignal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*os.Process).Signal(toType1(args[1]))
	p.Ret(2, ret0)
}

func execmProcessStateUserTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).UserTime()
	p.Ret(1, ret0)
}

func execmProcessStateSystemTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).SystemTime()
	p.Ret(1, ret0)
}

func execmProcessStateExited(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).Exited()
	p.Ret(1, ret0)
}

func execmProcessStateSuccess(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).Success()
	p.Ret(1, ret0)
}

func execmProcessStateSys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).Sys()
	p.Ret(1, ret0)
}

func execmProcessStateSysUsage(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).SysUsage()
	p.Ret(1, ret0)
}

func execmProcessStatePid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).Pid()
	p.Ret(1, ret0)
}

func execmProcessStateString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).String()
	p.Ret(1, ret0)
}

func execmProcessStateExitCode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.ProcessState).ExitCode()
	p.Ret(1, ret0)
}

func execReadlink(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.Readlink(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execRemove(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.Remove(args[0].(string))
	p.Ret(1, ret0)
}

func execRemoveAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.RemoveAll(args[0].(string))
	p.Ret(1, ret0)
}

func execRename(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Rename(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func toType2(v interface{}) os.FileInfo {
	if v == nil {
		return nil
	}
	return v.(os.FileInfo)
}

func execSameFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.SameFile(toType2(args[0]), toType2(args[1]))
	p.Ret(2, ret0)
}

func execSetenv(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Setenv(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execiSignalSignal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(os.Signal).Signal()
	p.PopN(1)
}

func execiSignalString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(os.Signal).String()
	p.Ret(1, ret0)
}

func execStartProcess(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := os.StartProcess(args[0].(string), args[1].([]string), args[2].(*os.ProcAttr))
	p.Ret(3, ret0, ret1)
}

func execStat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := os.Stat(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execSymlink(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Symlink(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execmSyscallErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.SyscallError).Error()
	p.Ret(1, ret0)
}

func execmSyscallErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.SyscallError).Unwrap()
	p.Ret(1, ret0)
}

func execmSyscallErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*os.SyscallError).Timeout()
	p.Ret(1, ret0)
}

func execTempDir(_ int, p *gop.Context) {
	ret0 := os.TempDir()
	p.Ret(0, ret0)
}

func execTruncate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := os.Truncate(args[0].(string), args[1].(int64))
	p.Ret(2, ret0)
}

func execUnsetenv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := os.Unsetenv(args[0].(string))
	p.Ret(1, ret0)
}

func execUserCacheDir(_ int, p *gop.Context) {
	ret0, ret1 := os.UserCacheDir()
	p.Ret(0, ret0, ret1)
}

func execUserConfigDir(_ int, p *gop.Context) {
	ret0, ret1 := os.UserConfigDir()
	p.Ret(0, ret0, ret1)
}

func execUserHomeDir(_ int, p *gop.Context) {
	ret0, ret1 := os.UserHomeDir()
	p.Ret(0, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("os")

func init() {
	I.RegisterFuncs(
		I.Func("Chdir", os.Chdir, execChdir),
		I.Func("Chmod", os.Chmod, execChmod),
		I.Func("Chown", os.Chown, execChown),
		I.Func("Chtimes", os.Chtimes, execChtimes),
		I.Func("Clearenv", os.Clearenv, execClearenv),
		I.Func("Create", os.Create, execCreate),
		I.Func("Environ", os.Environ, execEnviron),
		I.Func("Executable", os.Executable, execExecutable),
		I.Func("Exit", os.Exit, execExit),
		I.Func("Expand", os.Expand, execExpand),
		I.Func("ExpandEnv", os.ExpandEnv, execExpandEnv),
		I.Func("(*File).Readdir", (*os.File).Readdir, execmFileReaddir),
		I.Func("(*File).Readdirnames", (*os.File).Readdirnames, execmFileReaddirnames),
		I.Func("(*File).Name", (*os.File).Name, execmFileName),
		I.Func("(*File).Read", (*os.File).Read, execmFileRead),
		I.Func("(*File).ReadAt", (*os.File).ReadAt, execmFileReadAt),
		I.Func("(*File).Write", (*os.File).Write, execmFileWrite),
		I.Func("(*File).WriteAt", (*os.File).WriteAt, execmFileWriteAt),
		I.Func("(*File).Seek", (*os.File).Seek, execmFileSeek),
		I.Func("(*File).WriteString", (*os.File).WriteString, execmFileWriteString),
		I.Func("(*File).Chmod", (*os.File).Chmod, execmFileChmod),
		I.Func("(*File).SetDeadline", (*os.File).SetDeadline, execmFileSetDeadline),
		I.Func("(*File).SetReadDeadline", (*os.File).SetReadDeadline, execmFileSetReadDeadline),
		I.Func("(*File).SetWriteDeadline", (*os.File).SetWriteDeadline, execmFileSetWriteDeadline),
		I.Func("(*File).SyscallConn", (*os.File).SyscallConn, execmFileSyscallConn),
		I.Func("(*File).Chown", (*os.File).Chown, execmFileChown),
		I.Func("(*File).Truncate", (*os.File).Truncate, execmFileTruncate),
		I.Func("(*File).Sync", (*os.File).Sync, execmFileSync),
		I.Func("(*File).Chdir", (*os.File).Chdir, execmFileChdir),
		I.Func("(*File).Fd", (*os.File).Fd, execmFileFd),
		I.Func("(*File).Close", (*os.File).Close, execmFileClose),
		I.Func("(*File).Stat", (*os.File).Stat, execmFileStat),
		I.Func("(FileInfo).IsDir", (os.FileInfo).IsDir, execiFileInfoIsDir),
		I.Func("(FileInfo).ModTime", (os.FileInfo).ModTime, execiFileInfoModTime),
		I.Func("(FileInfo).Mode", (os.FileInfo).Mode, execiFileInfoMode),
		I.Func("(FileInfo).Name", (os.FileInfo).Name, execiFileInfoName),
		I.Func("(FileInfo).Size", (os.FileInfo).Size, execiFileInfoSize),
		I.Func("(FileInfo).Sys", (os.FileInfo).Sys, execiFileInfoSys),
		I.Func("(FileMode).String", (os.FileMode).String, execmFileModeString),
		I.Func("(FileMode).IsDir", (os.FileMode).IsDir, execmFileModeIsDir),
		I.Func("(FileMode).IsRegular", (os.FileMode).IsRegular, execmFileModeIsRegular),
		I.Func("(FileMode).Perm", (os.FileMode).Perm, execmFileModePerm),
		I.Func("FindProcess", os.FindProcess, execFindProcess),
		I.Func("Getegid", os.Getegid, execGetegid),
		I.Func("Getenv", os.Getenv, execGetenv),
		I.Func("Geteuid", os.Geteuid, execGeteuid),
		I.Func("Getgid", os.Getgid, execGetgid),
		I.Func("Getgroups", os.Getgroups, execGetgroups),
		I.Func("Getpagesize", os.Getpagesize, execGetpagesize),
		I.Func("Getpid", os.Getpid, execGetpid),
		I.Func("Getppid", os.Getppid, execGetppid),
		I.Func("Getuid", os.Getuid, execGetuid),
		I.Func("Getwd", os.Getwd, execGetwd),
		I.Func("Hostname", os.Hostname, execHostname),
		I.Func("IsExist", os.IsExist, execIsExist),
		I.Func("IsNotExist", os.IsNotExist, execIsNotExist),
		I.Func("IsPathSeparator", os.IsPathSeparator, execIsPathSeparator),
		I.Func("IsPermission", os.IsPermission, execIsPermission),
		I.Func("IsTimeout", os.IsTimeout, execIsTimeout),
		I.Func("Lchown", os.Lchown, execLchown),
		I.Func("Link", os.Link, execLink),
		I.Func("(*LinkError).Error", (*os.LinkError).Error, execmLinkErrorError),
		I.Func("(*LinkError).Unwrap", (*os.LinkError).Unwrap, execmLinkErrorUnwrap),
		I.Func("LookupEnv", os.LookupEnv, execLookupEnv),
		I.Func("Lstat", os.Lstat, execLstat),
		I.Func("Mkdir", os.Mkdir, execMkdir),
		I.Func("MkdirAll", os.MkdirAll, execMkdirAll),
		I.Func("NewFile", os.NewFile, execNewFile),
		I.Func("NewSyscallError", os.NewSyscallError, execNewSyscallError),
		I.Func("Open", os.Open, execOpen),
		I.Func("OpenFile", os.OpenFile, execOpenFile),
		I.Func("(*PathError).Error", (*os.PathError).Error, execmPathErrorError),
		I.Func("(*PathError).Unwrap", (*os.PathError).Unwrap, execmPathErrorUnwrap),
		I.Func("(*PathError).Timeout", (*os.PathError).Timeout, execmPathErrorTimeout),
		I.Func("Pipe", os.Pipe, execPipe),
		I.Func("(*Process).Release", (*os.Process).Release, execmProcessRelease),
		I.Func("(*Process).Kill", (*os.Process).Kill, execmProcessKill),
		I.Func("(*Process).Wait", (*os.Process).Wait, execmProcessWait),
		I.Func("(*Process).Signal", (*os.Process).Signal, execmProcessSignal),
		I.Func("(*ProcessState).UserTime", (*os.ProcessState).UserTime, execmProcessStateUserTime),
		I.Func("(*ProcessState).SystemTime", (*os.ProcessState).SystemTime, execmProcessStateSystemTime),
		I.Func("(*ProcessState).Exited", (*os.ProcessState).Exited, execmProcessStateExited),
		I.Func("(*ProcessState).Success", (*os.ProcessState).Success, execmProcessStateSuccess),
		I.Func("(*ProcessState).Sys", (*os.ProcessState).Sys, execmProcessStateSys),
		I.Func("(*ProcessState).SysUsage", (*os.ProcessState).SysUsage, execmProcessStateSysUsage),
		I.Func("(*ProcessState).Pid", (*os.ProcessState).Pid, execmProcessStatePid),
		I.Func("(*ProcessState).String", (*os.ProcessState).String, execmProcessStateString),
		I.Func("(*ProcessState).ExitCode", (*os.ProcessState).ExitCode, execmProcessStateExitCode),
		I.Func("Readlink", os.Readlink, execReadlink),
		I.Func("Remove", os.Remove, execRemove),
		I.Func("RemoveAll", os.RemoveAll, execRemoveAll),
		I.Func("Rename", os.Rename, execRename),
		I.Func("SameFile", os.SameFile, execSameFile),
		I.Func("Setenv", os.Setenv, execSetenv),
		I.Func("(Signal).Signal", (os.Signal).Signal, execiSignalSignal),
		I.Func("(Signal).String", (os.Signal).String, execiSignalString),
		I.Func("StartProcess", os.StartProcess, execStartProcess),
		I.Func("Stat", os.Stat, execStat),
		I.Func("Symlink", os.Symlink, execSymlink),
		I.Func("(*SyscallError).Error", (*os.SyscallError).Error, execmSyscallErrorError),
		I.Func("(*SyscallError).Unwrap", (*os.SyscallError).Unwrap, execmSyscallErrorUnwrap),
		I.Func("(*SyscallError).Timeout", (*os.SyscallError).Timeout, execmSyscallErrorTimeout),
		I.Func("TempDir", os.TempDir, execTempDir),
		I.Func("Truncate", os.Truncate, execTruncate),
		I.Func("Unsetenv", os.Unsetenv, execUnsetenv),
		I.Func("UserCacheDir", os.UserCacheDir, execUserCacheDir),
		I.Func("UserConfigDir", os.UserConfigDir, execUserConfigDir),
		I.Func("UserHomeDir", os.UserHomeDir, execUserHomeDir),
	)
	I.RegisterVars(
		I.Var("Args", &os.Args),
		I.Var("ErrClosed", &os.ErrClosed),
		I.Var("ErrExist", &os.ErrExist),
		I.Var("ErrInvalid", &os.ErrInvalid),
		I.Var("ErrNoDeadline", &os.ErrNoDeadline),
		I.Var("ErrNotExist", &os.ErrNotExist),
		I.Var("ErrPermission", &os.ErrPermission),
		I.Var("Interrupt", &os.Interrupt),
		I.Var("Kill", &os.Kill),
		I.Var("Stderr", &os.Stderr),
		I.Var("Stdin", &os.Stdin),
		I.Var("Stdout", &os.Stdout),
	)
	I.RegisterTypes(
		I.Type("File", reflect.TypeOf((*os.File)(nil)).Elem()),
		I.Type("FileInfo", reflect.TypeOf((*os.FileInfo)(nil)).Elem()),
		I.Type("FileMode", reflect.TypeOf((*os.FileMode)(nil)).Elem()),
		I.Type("LinkError", reflect.TypeOf((*os.LinkError)(nil)).Elem()),
		I.Type("PathError", reflect.TypeOf((*os.PathError)(nil)).Elem()),
		I.Type("ProcAttr", reflect.TypeOf((*os.ProcAttr)(nil)).Elem()),
		I.Type("Process", reflect.TypeOf((*os.Process)(nil)).Elem()),
		I.Type("ProcessState", reflect.TypeOf((*os.ProcessState)(nil)).Elem()),
		I.Type("Signal", reflect.TypeOf((*os.Signal)(nil)).Elem()),
		I.Type("SyscallError", reflect.TypeOf((*os.SyscallError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DevNull", qspec.ConstBoundString, os.DevNull),
		I.Const("ModeAppend", qspec.Uint32, os.ModeAppend),
		I.Const("ModeCharDevice", qspec.Uint32, os.ModeCharDevice),
		I.Const("ModeDevice", qspec.Uint32, os.ModeDevice),
		I.Const("ModeDir", qspec.Uint32, os.ModeDir),
		I.Const("ModeExclusive", qspec.Uint32, os.ModeExclusive),
		I.Const("ModeIrregular", qspec.Uint32, os.ModeIrregular),
		I.Const("ModeNamedPipe", qspec.Uint32, os.ModeNamedPipe),
		I.Const("ModePerm", qspec.Uint32, os.ModePerm),
		I.Const("ModeSetgid", qspec.Uint32, os.ModeSetgid),
		I.Const("ModeSetuid", qspec.Uint32, os.ModeSetuid),
		I.Const("ModeSocket", qspec.Uint32, os.ModeSocket),
		I.Const("ModeSticky", qspec.Uint32, os.ModeSticky),
		I.Const("ModeSymlink", qspec.Uint32, os.ModeSymlink),
		I.Const("ModeTemporary", qspec.Uint32, os.ModeTemporary),
		I.Const("ModeType", qspec.Uint32, os.ModeType),
		I.Const("O_APPEND", qspec.Int, os.O_APPEND),
		I.Const("O_CREATE", qspec.Int, os.O_CREATE),
		I.Const("O_EXCL", qspec.Int, os.O_EXCL),
		I.Const("O_RDONLY", qspec.Int, os.O_RDONLY),
		I.Const("O_RDWR", qspec.Int, os.O_RDWR),
		I.Const("O_SYNC", qspec.Int, os.O_SYNC),
		I.Const("O_TRUNC", qspec.Int, os.O_TRUNC),
		I.Const("O_WRONLY", qspec.Int, os.O_WRONLY),
		I.Const("PathListSeparator", qspec.ConstBoundRune, os.PathListSeparator),
		I.Const("PathSeparator", qspec.ConstBoundRune, os.PathSeparator),
		I.Const("SEEK_CUR", qspec.Int, os.SEEK_CUR),
		I.Const("SEEK_END", qspec.Int, os.SEEK_END),
		I.Const("SEEK_SET", qspec.Int, os.SEEK_SET),
	)
}
