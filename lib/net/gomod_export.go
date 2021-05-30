// Package net provide Go+ "net" package, as "net" package in Go.
package net

import (
	context "context"
	io "io"
	net "net"
	os "os"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execiAddrNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Addr).Network()
	p.Ret(1, ret0)
}

func execiAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Addr).String()
	p.Ret(1, ret0)
}

func execmAddrErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.AddrError).Error()
	p.Ret(1, ret0)
}

func execmAddrErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.AddrError).Timeout()
	p.Ret(1, ret0)
}

func execmAddrErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.AddrError).Temporary()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmBuffersWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*net.Buffers).WriteTo(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmBuffersRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*net.Buffers).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execCIDRMask(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := net.CIDRMask(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execiConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).Close()
	p.Ret(1, ret0)
}

func execiConnLocalAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).LocalAddr()
	p.Ret(1, ret0)
}

func execiConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(net.Conn).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiConnRemoteAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).RemoteAddr()
	p.Ret(1, ret0)
}

func execiConnSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(net.Conn).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmDNSConfigErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSConfigError).Unwrap()
	p.Ret(1, ret0)
}

func execmDNSConfigErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSConfigError).Error()
	p.Ret(1, ret0)
}

func execmDNSConfigErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSConfigError).Timeout()
	p.Ret(1, ret0)
}

func execmDNSConfigErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSConfigError).Temporary()
	p.Ret(1, ret0)
}

func execmDNSErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSError).Error()
	p.Ret(1, ret0)
}

func execmDNSErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSError).Timeout()
	p.Ret(1, ret0)
}

func execmDNSErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.DNSError).Temporary()
	p.Ret(1, ret0)
}

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.Dial(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execDialIP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.DialIP(args[0].(string), args[1].(*net.IPAddr), args[2].(*net.IPAddr))
	p.Ret(3, ret0, ret1)
}

func execDialTCP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.DialTCP(args[0].(string), args[1].(*net.TCPAddr), args[2].(*net.TCPAddr))
	p.Ret(3, ret0, ret1)
}

func execDialTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.DialTimeout(args[0].(string), args[1].(string), args[2].(time.Duration))
	p.Ret(3, ret0, ret1)
}

func execDialUDP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.DialUDP(args[0].(string), args[1].(*net.UDPAddr), args[2].(*net.UDPAddr))
	p.Ret(3, ret0, ret1)
}

func execDialUnix(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.DialUnix(args[0].(string), args[1].(*net.UnixAddr), args[2].(*net.UnixAddr))
	p.Ret(3, ret0, ret1)
}

func execmDialerDial(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Dialer).Dial(args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func toType1(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execmDialerDialContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*net.Dialer).DialContext(toType1(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execiErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Error).Error()
	p.Ret(1, ret0)
}

func execiErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Error).Temporary()
	p.Ret(1, ret0)
}

func execiErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Error).Timeout()
	p.Ret(1, ret0)
}

func execFileConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.FileConn(args[0].(*os.File))
	p.Ret(1, ret0, ret1)
}

func execFileListener(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.FileListener(args[0].(*os.File))
	p.Ret(1, ret0, ret1)
}

func execFilePacketConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.FilePacketConn(args[0].(*os.File))
	p.Ret(1, ret0, ret1)
}

func execmFlagsString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Flags).String()
	p.Ret(1, ret0)
}

func execmHardwareAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.HardwareAddr).String()
	p.Ret(1, ret0)
}

func execmIPIsUnspecified(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsUnspecified()
	p.Ret(1, ret0)
}

func execmIPIsLoopback(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsLoopback()
	p.Ret(1, ret0)
}

func execmIPIsMulticast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsMulticast()
	p.Ret(1, ret0)
}

func execmIPIsInterfaceLocalMulticast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsInterfaceLocalMulticast()
	p.Ret(1, ret0)
}

func execmIPIsLinkLocalMulticast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsLinkLocalMulticast()
	p.Ret(1, ret0)
}

func execmIPIsLinkLocalUnicast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsLinkLocalUnicast()
	p.Ret(1, ret0)
}

func execmIPIsGlobalUnicast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).IsGlobalUnicast()
	p.Ret(1, ret0)
}

func execmIPTo4(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).To4()
	p.Ret(1, ret0)
}

func execmIPTo16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).To16()
	p.Ret(1, ret0)
}

func execmIPDefaultMask(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).DefaultMask()
	p.Ret(1, ret0)
}

func execmIPMask(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.IP).Mask(args[1].(net.IPMask))
	p.Ret(2, ret0)
}

func execmIPString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IP).String()
	p.Ret(1, ret0)
}

func execmIPMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(net.IP).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmIPUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.IP).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmIPEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.IP).Equal(args[1].(net.IP))
	p.Ret(2, ret0)
}

func execmIPAddrNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.IPAddr).Network()
	p.Ret(1, ret0)
}

func execmIPAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.IPAddr).String()
	p.Ret(1, ret0)
}

func execmIPConnSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.IPConn).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmIPConnReadFromIP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.IPConn).ReadFromIP(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmIPConnReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.IPConn).ReadFrom(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmIPConnReadMsgIP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2, ret3, ret4 := args[0].(*net.IPConn).ReadMsgIP(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0, ret1, ret2, ret3, ret4)
}

func execmIPConnWriteToIP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.IPConn).WriteToIP(args[1].([]byte), args[2].(*net.IPAddr))
	p.Ret(3, ret0, ret1)
}

func toType2(v interface{}) net.Addr {
	if v == nil {
		return nil
	}
	return v.(net.Addr)
}

func execmIPConnWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.IPConn).WriteTo(args[1].([]byte), toType2(args[2]))
	p.Ret(3, ret0, ret1)
}

func execmIPConnWriteMsgIP(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1, ret2 := args[0].(*net.IPConn).WriteMsgIP(args[1].([]byte), args[2].([]byte), args[3].(*net.IPAddr))
	p.Ret(4, ret0, ret1, ret2)
}

func execmIPMaskSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(net.IPMask).Size()
	p.Ret(1, ret0, ret1)
}

func execmIPMaskString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.IPMask).String()
	p.Ret(1, ret0)
}

func execmIPNetContains(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.IPNet).Contains(args[1].(net.IP))
	p.Ret(2, ret0)
}

func execmIPNetNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.IPNet).Network()
	p.Ret(1, ret0)
}

func execmIPNetString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.IPNet).String()
	p.Ret(1, ret0)
}

func execIPv4(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := net.IPv4(args[0].(byte), args[1].(byte), args[2].(byte), args[3].(byte))
	p.Ret(4, ret0)
}

func execIPv4Mask(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := net.IPv4Mask(args[0].(byte), args[1].(byte), args[2].(byte), args[3].(byte))
	p.Ret(4, ret0)
}

func execmInterfaceAddrs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.Interface).Addrs()
	p.Ret(1, ret0, ret1)
}

func execmInterfaceMulticastAddrs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.Interface).MulticastAddrs()
	p.Ret(1, ret0, ret1)
}

func execInterfaceAddrs(_ int, p *gop.Context) {
	ret0, ret1 := net.InterfaceAddrs()
	p.Ret(0, ret0, ret1)
}

func execInterfaceByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.InterfaceByIndex(args[0].(int))
	p.Ret(1, ret0, ret1)
}

func execInterfaceByName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.InterfaceByName(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execInterfaces(_ int, p *gop.Context) {
	ret0, ret1 := net.Interfaces()
	p.Ret(0, ret0, ret1)
}

func execmInvalidAddrErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.InvalidAddrError).Error()
	p.Ret(1, ret0)
}

func execmInvalidAddrErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.InvalidAddrError).Timeout()
	p.Ret(1, ret0)
}

func execmInvalidAddrErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.InvalidAddrError).Temporary()
	p.Ret(1, ret0)
}

func execJoinHostPort(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := net.JoinHostPort(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execListen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.Listen(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmListenConfigListen(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*net.ListenConfig).Listen(toType1(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execmListenConfigListenPacket(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*net.ListenConfig).ListenPacket(toType1(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execListenIP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenIP(args[0].(string), args[1].(*net.IPAddr))
	p.Ret(2, ret0, ret1)
}

func execListenMulticastUDP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := net.ListenMulticastUDP(args[0].(string), args[1].(*net.Interface), args[2].(*net.UDPAddr))
	p.Ret(3, ret0, ret1)
}

func execListenPacket(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenPacket(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execListenTCP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenTCP(args[0].(string), args[1].(*net.TCPAddr))
	p.Ret(2, ret0, ret1)
}

func execListenUDP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenUDP(args[0].(string), args[1].(*net.UDPAddr))
	p.Ret(2, ret0, ret1)
}

func execListenUnix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenUnix(args[0].(string), args[1].(*net.UnixAddr))
	p.Ret(2, ret0, ret1)
}

func execListenUnixgram(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ListenUnixgram(args[0].(string), args[1].(*net.UnixAddr))
	p.Ret(2, ret0, ret1)
}

func execiListenerAccept(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(net.Listener).Accept()
	p.Ret(1, ret0, ret1)
}

func execiListenerAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Listener).Addr()
	p.Ret(1, ret0)
}

func execiListenerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Listener).Close()
	p.Ret(1, ret0)
}

func execLookupAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupAddr(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupCNAME(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupCNAME(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupHost(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupHost(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupIP(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupIP(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupMX(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupMX(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupNS(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupNS(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLookupPort(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.LookupPort(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execLookupSRV(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := net.LookupSRV(args[0].(string), args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1, ret2)
}

func execLookupTXT(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.LookupTXT(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmOpErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.OpError).Unwrap()
	p.Ret(1, ret0)
}

func execmOpErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.OpError).Error()
	p.Ret(1, ret0)
}

func execmOpErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.OpError).Timeout()
	p.Ret(1, ret0)
}

func execmOpErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.OpError).Temporary()
	p.Ret(1, ret0)
}

func execiPacketConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.PacketConn).Close()
	p.Ret(1, ret0)
}

func execiPacketConnLocalAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.PacketConn).LocalAddr()
	p.Ret(1, ret0)
}

func execiPacketConnReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(net.PacketConn).ReadFrom(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execiPacketConnSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.PacketConn).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiPacketConnSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.PacketConn).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiPacketConnSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.PacketConn).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiPacketConnWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(net.PacketConn).WriteTo(args[1].([]byte), toType2(args[2]))
	p.Ret(3, ret0, ret1)
}

func execParseCIDR(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := net.ParseCIDR(args[0].(string))
	p.Ret(1, ret0, ret1, ret2)
}

func execmParseErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.ParseError).Error()
	p.Ret(1, ret0)
}

func execParseIP(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := net.ParseIP(args[0].(string))
	p.Ret(1, ret0)
}

func execParseMAC(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := net.ParseMAC(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execPipe(_ int, p *gop.Context) {
	ret0, ret1 := net.Pipe()
	p.Ret(0, ret0, ret1)
}

func execResolveIPAddr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ResolveIPAddr(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execResolveTCPAddr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ResolveTCPAddr(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execResolveUDPAddr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ResolveUDPAddr(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execResolveUnixAddr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := net.ResolveUnixAddr(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmResolverLookupHost(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupHost(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupIPAddr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupIPAddr(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupPort(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*net.Resolver).LookupPort(toType1(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execmResolverLookupCNAME(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupCNAME(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupSRV(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1, ret2 := args[0].(*net.Resolver).LookupSRV(toType1(args[1]), args[2].(string), args[3].(string), args[4].(string))
	p.Ret(5, ret0, ret1, ret2)
}

func execmResolverLookupMX(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupMX(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupNS(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupNS(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupTXT(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupTXT(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmResolverLookupAddr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.Resolver).LookupAddr(toType1(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execSplitHostPort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := net.SplitHostPort(args[0].(string))
	p.Ret(1, ret0, ret1, ret2)
}

func execmTCPAddrNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPAddr).Network()
	p.Ret(1, ret0)
}

func execmTCPAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPAddr).String()
	p.Ret(1, ret0)
}

func execmTCPConnSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.TCPConn).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func toType3(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execmTCPConnReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*net.TCPConn).ReadFrom(toType3(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmTCPConnCloseRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPConn).CloseRead()
	p.Ret(1, ret0)
}

func execmTCPConnCloseWrite(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPConn).CloseWrite()
	p.Ret(1, ret0)
}

func execmTCPConnSetLinger(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.TCPConn).SetLinger(args[1].(int))
	p.Ret(2, ret0)
}

func execmTCPConnSetKeepAlive(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.TCPConn).SetKeepAlive(args[1].(bool))
	p.Ret(2, ret0)
}

func execmTCPConnSetKeepAlivePeriod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.TCPConn).SetKeepAlivePeriod(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execmTCPConnSetNoDelay(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.TCPConn).SetNoDelay(args[1].(bool))
	p.Ret(2, ret0)
}

func execmTCPListenerSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.TCPListener).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmTCPListenerAcceptTCP(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.TCPListener).AcceptTCP()
	p.Ret(1, ret0, ret1)
}

func execmTCPListenerAccept(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.TCPListener).Accept()
	p.Ret(1, ret0, ret1)
}

func execmTCPListenerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPListener).Close()
	p.Ret(1, ret0)
}

func execmTCPListenerAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.TCPListener).Addr()
	p.Ret(1, ret0)
}

func execmTCPListenerSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.TCPListener).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmTCPListenerFile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.TCPListener).File()
	p.Ret(1, ret0, ret1)
}

func execmUDPAddrNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UDPAddr).Network()
	p.Ret(1, ret0)
}

func execmUDPAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UDPAddr).String()
	p.Ret(1, ret0)
}

func execmUDPConnSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UDPConn).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmUDPConnReadFromUDP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.UDPConn).ReadFromUDP(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmUDPConnReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.UDPConn).ReadFrom(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmUDPConnReadMsgUDP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2, ret3, ret4 := args[0].(*net.UDPConn).ReadMsgUDP(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0, ret1, ret2, ret3, ret4)
}

func execmUDPConnWriteToUDP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.UDPConn).WriteToUDP(args[1].([]byte), args[2].(*net.UDPAddr))
	p.Ret(3, ret0, ret1)
}

func execmUDPConnWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.UDPConn).WriteTo(args[1].([]byte), toType2(args[2]))
	p.Ret(3, ret0, ret1)
}

func execmUDPConnWriteMsgUDP(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1, ret2 := args[0].(*net.UDPConn).WriteMsgUDP(args[1].([]byte), args[2].([]byte), args[3].(*net.UDPAddr))
	p.Ret(4, ret0, ret1, ret2)
}

func execmUnixAddrNetwork(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixAddr).Network()
	p.Ret(1, ret0)
}

func execmUnixAddrString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixAddr).String()
	p.Ret(1, ret0)
}

func execmUnixConnSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UnixConn).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmUnixConnCloseRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixConn).CloseRead()
	p.Ret(1, ret0)
}

func execmUnixConnCloseWrite(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixConn).CloseWrite()
	p.Ret(1, ret0)
}

func execmUnixConnReadFromUnix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.UnixConn).ReadFromUnix(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmUnixConnReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*net.UnixConn).ReadFrom(args[1].([]byte))
	p.Ret(2, ret0, ret1, ret2)
}

func execmUnixConnReadMsgUnix(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2, ret3, ret4 := args[0].(*net.UnixConn).ReadMsgUnix(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0, ret1, ret2, ret3, ret4)
}

func execmUnixConnWriteToUnix(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.UnixConn).WriteToUnix(args[1].([]byte), args[2].(*net.UnixAddr))
	p.Ret(3, ret0, ret1)
}

func execmUnixConnWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*net.UnixConn).WriteTo(args[1].([]byte), toType2(args[2]))
	p.Ret(3, ret0, ret1)
}

func execmUnixConnWriteMsgUnix(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1, ret2 := args[0].(*net.UnixConn).WriteMsgUnix(args[1].([]byte), args[2].([]byte), args[3].(*net.UnixAddr))
	p.Ret(4, ret0, ret1, ret2)
}

func execmUnixListenerSyscallConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UnixListener).SyscallConn()
	p.Ret(1, ret0, ret1)
}

func execmUnixListenerAcceptUnix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UnixListener).AcceptUnix()
	p.Ret(1, ret0, ret1)
}

func execmUnixListenerAccept(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UnixListener).Accept()
	p.Ret(1, ret0, ret1)
}

func execmUnixListenerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixListener).Close()
	p.Ret(1, ret0)
}

func execmUnixListenerAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*net.UnixListener).Addr()
	p.Ret(1, ret0)
}

func execmUnixListenerSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*net.UnixListener).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmUnixListenerFile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*net.UnixListener).File()
	p.Ret(1, ret0, ret1)
}

func execmUnixListenerSetUnlinkOnClose(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*net.UnixListener).SetUnlinkOnClose(args[1].(bool))
	p.PopN(2)
}

func execmUnknownNetworkErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.UnknownNetworkError).Error()
	p.Ret(1, ret0)
}

func execmUnknownNetworkErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.UnknownNetworkError).Timeout()
	p.Ret(1, ret0)
}

func execmUnknownNetworkErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.UnknownNetworkError).Temporary()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net")

func init() {
	I.RegisterFuncs(
		I.Func("(Addr).Network", (net.Addr).Network, execiAddrNetwork),
		I.Func("(Addr).String", (net.Addr).String, execiAddrString),
		I.Func("(*AddrError).Error", (*net.AddrError).Error, execmAddrErrorError),
		I.Func("(*AddrError).Timeout", (*net.AddrError).Timeout, execmAddrErrorTimeout),
		I.Func("(*AddrError).Temporary", (*net.AddrError).Temporary, execmAddrErrorTemporary),
		I.Func("(*Buffers).WriteTo", (*net.Buffers).WriteTo, execmBuffersWriteTo),
		I.Func("(*Buffers).Read", (*net.Buffers).Read, execmBuffersRead),
		I.Func("CIDRMask", net.CIDRMask, execCIDRMask),
		I.Func("(Conn).Close", (net.Conn).Close, execiConnClose),
		I.Func("(Conn).LocalAddr", (net.Conn).LocalAddr, execiConnLocalAddr),
		I.Func("(Conn).Read", (net.Conn).Read, execiConnRead),
		I.Func("(Conn).RemoteAddr", (net.Conn).RemoteAddr, execiConnRemoteAddr),
		I.Func("(Conn).SetDeadline", (net.Conn).SetDeadline, execiConnSetDeadline),
		I.Func("(Conn).SetReadDeadline", (net.Conn).SetReadDeadline, execiConnSetReadDeadline),
		I.Func("(Conn).SetWriteDeadline", (net.Conn).SetWriteDeadline, execiConnSetWriteDeadline),
		I.Func("(Conn).Write", (net.Conn).Write, execiConnWrite),
		I.Func("(*DNSConfigError).Unwrap", (*net.DNSConfigError).Unwrap, execmDNSConfigErrorUnwrap),
		I.Func("(*DNSConfigError).Error", (*net.DNSConfigError).Error, execmDNSConfigErrorError),
		I.Func("(*DNSConfigError).Timeout", (*net.DNSConfigError).Timeout, execmDNSConfigErrorTimeout),
		I.Func("(*DNSConfigError).Temporary", (*net.DNSConfigError).Temporary, execmDNSConfigErrorTemporary),
		I.Func("(*DNSError).Error", (*net.DNSError).Error, execmDNSErrorError),
		I.Func("(*DNSError).Timeout", (*net.DNSError).Timeout, execmDNSErrorTimeout),
		I.Func("(*DNSError).Temporary", (*net.DNSError).Temporary, execmDNSErrorTemporary),
		I.Func("Dial", net.Dial, execDial),
		I.Func("DialIP", net.DialIP, execDialIP),
		I.Func("DialTCP", net.DialTCP, execDialTCP),
		I.Func("DialTimeout", net.DialTimeout, execDialTimeout),
		I.Func("DialUDP", net.DialUDP, execDialUDP),
		I.Func("DialUnix", net.DialUnix, execDialUnix),
		I.Func("(*Dialer).Dial", (*net.Dialer).Dial, execmDialerDial),
		I.Func("(*Dialer).DialContext", (*net.Dialer).DialContext, execmDialerDialContext),
		I.Func("(Error).Error", (net.Error).Error, execiErrorError),
		I.Func("(Error).Temporary", (net.Error).Temporary, execiErrorTemporary),
		I.Func("(Error).Timeout", (net.Error).Timeout, execiErrorTimeout),
		I.Func("FileConn", net.FileConn, execFileConn),
		I.Func("FileListener", net.FileListener, execFileListener),
		I.Func("FilePacketConn", net.FilePacketConn, execFilePacketConn),
		I.Func("(Flags).String", (net.Flags).String, execmFlagsString),
		I.Func("(HardwareAddr).String", (net.HardwareAddr).String, execmHardwareAddrString),
		I.Func("(IP).IsUnspecified", (net.IP).IsUnspecified, execmIPIsUnspecified),
		I.Func("(IP).IsLoopback", (net.IP).IsLoopback, execmIPIsLoopback),
		I.Func("(IP).IsMulticast", (net.IP).IsMulticast, execmIPIsMulticast),
		I.Func("(IP).IsInterfaceLocalMulticast", (net.IP).IsInterfaceLocalMulticast, execmIPIsInterfaceLocalMulticast),
		I.Func("(IP).IsLinkLocalMulticast", (net.IP).IsLinkLocalMulticast, execmIPIsLinkLocalMulticast),
		I.Func("(IP).IsLinkLocalUnicast", (net.IP).IsLinkLocalUnicast, execmIPIsLinkLocalUnicast),
		I.Func("(IP).IsGlobalUnicast", (net.IP).IsGlobalUnicast, execmIPIsGlobalUnicast),
		I.Func("(IP).To4", (net.IP).To4, execmIPTo4),
		I.Func("(IP).To16", (net.IP).To16, execmIPTo16),
		I.Func("(IP).DefaultMask", (net.IP).DefaultMask, execmIPDefaultMask),
		I.Func("(IP).Mask", (net.IP).Mask, execmIPMask),
		I.Func("(IP).String", (net.IP).String, execmIPString),
		I.Func("(IP).MarshalText", (net.IP).MarshalText, execmIPMarshalText),
		I.Func("(*IP).UnmarshalText", (*net.IP).UnmarshalText, execmIPUnmarshalText),
		I.Func("(IP).Equal", (net.IP).Equal, execmIPEqual),
		I.Func("(*IPAddr).Network", (*net.IPAddr).Network, execmIPAddrNetwork),
		I.Func("(*IPAddr).String", (*net.IPAddr).String, execmIPAddrString),
		I.Func("(*IPConn).SyscallConn", (*net.IPConn).SyscallConn, execmIPConnSyscallConn),
		I.Func("(*IPConn).ReadFromIP", (*net.IPConn).ReadFromIP, execmIPConnReadFromIP),
		I.Func("(*IPConn).ReadFrom", (*net.IPConn).ReadFrom, execmIPConnReadFrom),
		I.Func("(*IPConn).ReadMsgIP", (*net.IPConn).ReadMsgIP, execmIPConnReadMsgIP),
		I.Func("(*IPConn).WriteToIP", (*net.IPConn).WriteToIP, execmIPConnWriteToIP),
		I.Func("(*IPConn).WriteTo", (*net.IPConn).WriteTo, execmIPConnWriteTo),
		I.Func("(*IPConn).WriteMsgIP", (*net.IPConn).WriteMsgIP, execmIPConnWriteMsgIP),
		I.Func("(IPMask).Size", (net.IPMask).Size, execmIPMaskSize),
		I.Func("(IPMask).String", (net.IPMask).String, execmIPMaskString),
		I.Func("(*IPNet).Contains", (*net.IPNet).Contains, execmIPNetContains),
		I.Func("(*IPNet).Network", (*net.IPNet).Network, execmIPNetNetwork),
		I.Func("(*IPNet).String", (*net.IPNet).String, execmIPNetString),
		I.Func("IPv4", net.IPv4, execIPv4),
		I.Func("IPv4Mask", net.IPv4Mask, execIPv4Mask),
		I.Func("(*Interface).Addrs", (*net.Interface).Addrs, execmInterfaceAddrs),
		I.Func("(*Interface).MulticastAddrs", (*net.Interface).MulticastAddrs, execmInterfaceMulticastAddrs),
		I.Func("InterfaceAddrs", net.InterfaceAddrs, execInterfaceAddrs),
		I.Func("InterfaceByIndex", net.InterfaceByIndex, execInterfaceByIndex),
		I.Func("InterfaceByName", net.InterfaceByName, execInterfaceByName),
		I.Func("Interfaces", net.Interfaces, execInterfaces),
		I.Func("(InvalidAddrError).Error", (net.InvalidAddrError).Error, execmInvalidAddrErrorError),
		I.Func("(InvalidAddrError).Timeout", (net.InvalidAddrError).Timeout, execmInvalidAddrErrorTimeout),
		I.Func("(InvalidAddrError).Temporary", (net.InvalidAddrError).Temporary, execmInvalidAddrErrorTemporary),
		I.Func("JoinHostPort", net.JoinHostPort, execJoinHostPort),
		I.Func("Listen", net.Listen, execListen),
		I.Func("(*ListenConfig).Listen", (*net.ListenConfig).Listen, execmListenConfigListen),
		I.Func("(*ListenConfig).ListenPacket", (*net.ListenConfig).ListenPacket, execmListenConfigListenPacket),
		I.Func("ListenIP", net.ListenIP, execListenIP),
		I.Func("ListenMulticastUDP", net.ListenMulticastUDP, execListenMulticastUDP),
		I.Func("ListenPacket", net.ListenPacket, execListenPacket),
		I.Func("ListenTCP", net.ListenTCP, execListenTCP),
		I.Func("ListenUDP", net.ListenUDP, execListenUDP),
		I.Func("ListenUnix", net.ListenUnix, execListenUnix),
		I.Func("ListenUnixgram", net.ListenUnixgram, execListenUnixgram),
		I.Func("(Listener).Accept", (net.Listener).Accept, execiListenerAccept),
		I.Func("(Listener).Addr", (net.Listener).Addr, execiListenerAddr),
		I.Func("(Listener).Close", (net.Listener).Close, execiListenerClose),
		I.Func("LookupAddr", net.LookupAddr, execLookupAddr),
		I.Func("LookupCNAME", net.LookupCNAME, execLookupCNAME),
		I.Func("LookupHost", net.LookupHost, execLookupHost),
		I.Func("LookupIP", net.LookupIP, execLookupIP),
		I.Func("LookupMX", net.LookupMX, execLookupMX),
		I.Func("LookupNS", net.LookupNS, execLookupNS),
		I.Func("LookupPort", net.LookupPort, execLookupPort),
		I.Func("LookupSRV", net.LookupSRV, execLookupSRV),
		I.Func("LookupTXT", net.LookupTXT, execLookupTXT),
		I.Func("(*OpError).Unwrap", (*net.OpError).Unwrap, execmOpErrorUnwrap),
		I.Func("(*OpError).Error", (*net.OpError).Error, execmOpErrorError),
		I.Func("(*OpError).Timeout", (*net.OpError).Timeout, execmOpErrorTimeout),
		I.Func("(*OpError).Temporary", (*net.OpError).Temporary, execmOpErrorTemporary),
		I.Func("(PacketConn).Close", (net.PacketConn).Close, execiPacketConnClose),
		I.Func("(PacketConn).LocalAddr", (net.PacketConn).LocalAddr, execiPacketConnLocalAddr),
		I.Func("(PacketConn).ReadFrom", (net.PacketConn).ReadFrom, execiPacketConnReadFrom),
		I.Func("(PacketConn).SetDeadline", (net.PacketConn).SetDeadline, execiPacketConnSetDeadline),
		I.Func("(PacketConn).SetReadDeadline", (net.PacketConn).SetReadDeadline, execiPacketConnSetReadDeadline),
		I.Func("(PacketConn).SetWriteDeadline", (net.PacketConn).SetWriteDeadline, execiPacketConnSetWriteDeadline),
		I.Func("(PacketConn).WriteTo", (net.PacketConn).WriteTo, execiPacketConnWriteTo),
		I.Func("ParseCIDR", net.ParseCIDR, execParseCIDR),
		I.Func("(*ParseError).Error", (*net.ParseError).Error, execmParseErrorError),
		I.Func("ParseIP", net.ParseIP, execParseIP),
		I.Func("ParseMAC", net.ParseMAC, execParseMAC),
		I.Func("Pipe", net.Pipe, execPipe),
		I.Func("ResolveIPAddr", net.ResolveIPAddr, execResolveIPAddr),
		I.Func("ResolveTCPAddr", net.ResolveTCPAddr, execResolveTCPAddr),
		I.Func("ResolveUDPAddr", net.ResolveUDPAddr, execResolveUDPAddr),
		I.Func("ResolveUnixAddr", net.ResolveUnixAddr, execResolveUnixAddr),
		I.Func("(*Resolver).LookupHost", (*net.Resolver).LookupHost, execmResolverLookupHost),
		I.Func("(*Resolver).LookupIPAddr", (*net.Resolver).LookupIPAddr, execmResolverLookupIPAddr),
		I.Func("(*Resolver).LookupPort", (*net.Resolver).LookupPort, execmResolverLookupPort),
		I.Func("(*Resolver).LookupCNAME", (*net.Resolver).LookupCNAME, execmResolverLookupCNAME),
		I.Func("(*Resolver).LookupSRV", (*net.Resolver).LookupSRV, execmResolverLookupSRV),
		I.Func("(*Resolver).LookupMX", (*net.Resolver).LookupMX, execmResolverLookupMX),
		I.Func("(*Resolver).LookupNS", (*net.Resolver).LookupNS, execmResolverLookupNS),
		I.Func("(*Resolver).LookupTXT", (*net.Resolver).LookupTXT, execmResolverLookupTXT),
		I.Func("(*Resolver).LookupAddr", (*net.Resolver).LookupAddr, execmResolverLookupAddr),
		I.Func("SplitHostPort", net.SplitHostPort, execSplitHostPort),
		I.Func("(*TCPAddr).Network", (*net.TCPAddr).Network, execmTCPAddrNetwork),
		I.Func("(*TCPAddr).String", (*net.TCPAddr).String, execmTCPAddrString),
		I.Func("(*TCPConn).SyscallConn", (*net.TCPConn).SyscallConn, execmTCPConnSyscallConn),
		I.Func("(*TCPConn).ReadFrom", (*net.TCPConn).ReadFrom, execmTCPConnReadFrom),
		I.Func("(*TCPConn).CloseRead", (*net.TCPConn).CloseRead, execmTCPConnCloseRead),
		I.Func("(*TCPConn).CloseWrite", (*net.TCPConn).CloseWrite, execmTCPConnCloseWrite),
		I.Func("(*TCPConn).SetLinger", (*net.TCPConn).SetLinger, execmTCPConnSetLinger),
		I.Func("(*TCPConn).SetKeepAlive", (*net.TCPConn).SetKeepAlive, execmTCPConnSetKeepAlive),
		I.Func("(*TCPConn).SetKeepAlivePeriod", (*net.TCPConn).SetKeepAlivePeriod, execmTCPConnSetKeepAlivePeriod),
		I.Func("(*TCPConn).SetNoDelay", (*net.TCPConn).SetNoDelay, execmTCPConnSetNoDelay),
		I.Func("(*TCPListener).SyscallConn", (*net.TCPListener).SyscallConn, execmTCPListenerSyscallConn),
		I.Func("(*TCPListener).AcceptTCP", (*net.TCPListener).AcceptTCP, execmTCPListenerAcceptTCP),
		I.Func("(*TCPListener).Accept", (*net.TCPListener).Accept, execmTCPListenerAccept),
		I.Func("(*TCPListener).Close", (*net.TCPListener).Close, execmTCPListenerClose),
		I.Func("(*TCPListener).Addr", (*net.TCPListener).Addr, execmTCPListenerAddr),
		I.Func("(*TCPListener).SetDeadline", (*net.TCPListener).SetDeadline, execmTCPListenerSetDeadline),
		I.Func("(*TCPListener).File", (*net.TCPListener).File, execmTCPListenerFile),
		I.Func("(*UDPAddr).Network", (*net.UDPAddr).Network, execmUDPAddrNetwork),
		I.Func("(*UDPAddr).String", (*net.UDPAddr).String, execmUDPAddrString),
		I.Func("(*UDPConn).SyscallConn", (*net.UDPConn).SyscallConn, execmUDPConnSyscallConn),
		I.Func("(*UDPConn).ReadFromUDP", (*net.UDPConn).ReadFromUDP, execmUDPConnReadFromUDP),
		I.Func("(*UDPConn).ReadFrom", (*net.UDPConn).ReadFrom, execmUDPConnReadFrom),
		I.Func("(*UDPConn).ReadMsgUDP", (*net.UDPConn).ReadMsgUDP, execmUDPConnReadMsgUDP),
		I.Func("(*UDPConn).WriteToUDP", (*net.UDPConn).WriteToUDP, execmUDPConnWriteToUDP),
		I.Func("(*UDPConn).WriteTo", (*net.UDPConn).WriteTo, execmUDPConnWriteTo),
		I.Func("(*UDPConn).WriteMsgUDP", (*net.UDPConn).WriteMsgUDP, execmUDPConnWriteMsgUDP),
		I.Func("(*UnixAddr).Network", (*net.UnixAddr).Network, execmUnixAddrNetwork),
		I.Func("(*UnixAddr).String", (*net.UnixAddr).String, execmUnixAddrString),
		I.Func("(*UnixConn).SyscallConn", (*net.UnixConn).SyscallConn, execmUnixConnSyscallConn),
		I.Func("(*UnixConn).CloseRead", (*net.UnixConn).CloseRead, execmUnixConnCloseRead),
		I.Func("(*UnixConn).CloseWrite", (*net.UnixConn).CloseWrite, execmUnixConnCloseWrite),
		I.Func("(*UnixConn).ReadFromUnix", (*net.UnixConn).ReadFromUnix, execmUnixConnReadFromUnix),
		I.Func("(*UnixConn).ReadFrom", (*net.UnixConn).ReadFrom, execmUnixConnReadFrom),
		I.Func("(*UnixConn).ReadMsgUnix", (*net.UnixConn).ReadMsgUnix, execmUnixConnReadMsgUnix),
		I.Func("(*UnixConn).WriteToUnix", (*net.UnixConn).WriteToUnix, execmUnixConnWriteToUnix),
		I.Func("(*UnixConn).WriteTo", (*net.UnixConn).WriteTo, execmUnixConnWriteTo),
		I.Func("(*UnixConn).WriteMsgUnix", (*net.UnixConn).WriteMsgUnix, execmUnixConnWriteMsgUnix),
		I.Func("(*UnixListener).SyscallConn", (*net.UnixListener).SyscallConn, execmUnixListenerSyscallConn),
		I.Func("(*UnixListener).AcceptUnix", (*net.UnixListener).AcceptUnix, execmUnixListenerAcceptUnix),
		I.Func("(*UnixListener).Accept", (*net.UnixListener).Accept, execmUnixListenerAccept),
		I.Func("(*UnixListener).Close", (*net.UnixListener).Close, execmUnixListenerClose),
		I.Func("(*UnixListener).Addr", (*net.UnixListener).Addr, execmUnixListenerAddr),
		I.Func("(*UnixListener).SetDeadline", (*net.UnixListener).SetDeadline, execmUnixListenerSetDeadline),
		I.Func("(*UnixListener).File", (*net.UnixListener).File, execmUnixListenerFile),
		I.Func("(*UnixListener).SetUnlinkOnClose", (*net.UnixListener).SetUnlinkOnClose, execmUnixListenerSetUnlinkOnClose),
		I.Func("(UnknownNetworkError).Error", (net.UnknownNetworkError).Error, execmUnknownNetworkErrorError),
		I.Func("(UnknownNetworkError).Timeout", (net.UnknownNetworkError).Timeout, execmUnknownNetworkErrorTimeout),
		I.Func("(UnknownNetworkError).Temporary", (net.UnknownNetworkError).Temporary, execmUnknownNetworkErrorTemporary),
	)
	I.RegisterVars(
		I.Var("DefaultResolver", &net.DefaultResolver),
		I.Var("ErrWriteToConnected", &net.ErrWriteToConnected),
		I.Var("IPv4allrouter", &net.IPv4allrouter),
		I.Var("IPv4allsys", &net.IPv4allsys),
		I.Var("IPv4bcast", &net.IPv4bcast),
		I.Var("IPv4zero", &net.IPv4zero),
		I.Var("IPv6interfacelocalallnodes", &net.IPv6interfacelocalallnodes),
		I.Var("IPv6linklocalallnodes", &net.IPv6linklocalallnodes),
		I.Var("IPv6linklocalallrouters", &net.IPv6linklocalallrouters),
		I.Var("IPv6loopback", &net.IPv6loopback),
		I.Var("IPv6unspecified", &net.IPv6unspecified),
		I.Var("IPv6zero", &net.IPv6zero),
	)
	I.RegisterTypes(
		I.Type("Addr", reflect.TypeOf((*net.Addr)(nil)).Elem()),
		I.Type("AddrError", reflect.TypeOf((*net.AddrError)(nil)).Elem()),
		I.Type("Buffers", reflect.TypeOf((*net.Buffers)(nil)).Elem()),
		I.Type("Conn", reflect.TypeOf((*net.Conn)(nil)).Elem()),
		I.Type("DNSConfigError", reflect.TypeOf((*net.DNSConfigError)(nil)).Elem()),
		I.Type("DNSError", reflect.TypeOf((*net.DNSError)(nil)).Elem()),
		I.Type("Dialer", reflect.TypeOf((*net.Dialer)(nil)).Elem()),
		I.Type("Error", reflect.TypeOf((*net.Error)(nil)).Elem()),
		I.Type("Flags", reflect.TypeOf((*net.Flags)(nil)).Elem()),
		I.Type("HardwareAddr", reflect.TypeOf((*net.HardwareAddr)(nil)).Elem()),
		I.Type("IP", reflect.TypeOf((*net.IP)(nil)).Elem()),
		I.Type("IPAddr", reflect.TypeOf((*net.IPAddr)(nil)).Elem()),
		I.Type("IPConn", reflect.TypeOf((*net.IPConn)(nil)).Elem()),
		I.Type("IPMask", reflect.TypeOf((*net.IPMask)(nil)).Elem()),
		I.Type("IPNet", reflect.TypeOf((*net.IPNet)(nil)).Elem()),
		I.Type("Interface", reflect.TypeOf((*net.Interface)(nil)).Elem()),
		I.Type("InvalidAddrError", reflect.TypeOf((*net.InvalidAddrError)(nil)).Elem()),
		I.Type("ListenConfig", reflect.TypeOf((*net.ListenConfig)(nil)).Elem()),
		I.Type("Listener", reflect.TypeOf((*net.Listener)(nil)).Elem()),
		I.Type("MX", reflect.TypeOf((*net.MX)(nil)).Elem()),
		I.Type("NS", reflect.TypeOf((*net.NS)(nil)).Elem()),
		I.Type("OpError", reflect.TypeOf((*net.OpError)(nil)).Elem()),
		I.Type("PacketConn", reflect.TypeOf((*net.PacketConn)(nil)).Elem()),
		I.Type("ParseError", reflect.TypeOf((*net.ParseError)(nil)).Elem()),
		I.Type("Resolver", reflect.TypeOf((*net.Resolver)(nil)).Elem()),
		I.Type("SRV", reflect.TypeOf((*net.SRV)(nil)).Elem()),
		I.Type("TCPAddr", reflect.TypeOf((*net.TCPAddr)(nil)).Elem()),
		I.Type("TCPConn", reflect.TypeOf((*net.TCPConn)(nil)).Elem()),
		I.Type("TCPListener", reflect.TypeOf((*net.TCPListener)(nil)).Elem()),
		I.Type("UDPAddr", reflect.TypeOf((*net.UDPAddr)(nil)).Elem()),
		I.Type("UDPConn", reflect.TypeOf((*net.UDPConn)(nil)).Elem()),
		I.Type("UnixAddr", reflect.TypeOf((*net.UnixAddr)(nil)).Elem()),
		I.Type("UnixConn", reflect.TypeOf((*net.UnixConn)(nil)).Elem()),
		I.Type("UnixListener", reflect.TypeOf((*net.UnixListener)(nil)).Elem()),
		I.Type("UnknownNetworkError", reflect.TypeOf((*net.UnknownNetworkError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("FlagBroadcast", qspec.Uint, net.FlagBroadcast),
		I.Const("FlagLoopback", qspec.Uint, net.FlagLoopback),
		I.Const("FlagMulticast", qspec.Uint, net.FlagMulticast),
		I.Const("FlagPointToPoint", qspec.Uint, net.FlagPointToPoint),
		I.Const("FlagUp", qspec.Uint, net.FlagUp),
		I.Const("IPv4len", qspec.ConstUnboundInt, net.IPv4len),
		I.Const("IPv6len", qspec.ConstUnboundInt, net.IPv6len),
	)
}
