package helper

import (
	"crypto/tls"
	"net"
	"reflect"
)

// 获取连接的FD
func SocketFD(conn net.Conn) int {

	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	if _, ok := conn.(*tls.Conn); ok {
		tcpConn = reflect.Indirect(tcpConn.Elem())
	}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

// 获取 Listener FD
func ListenerFd(l net.Listener) int {
	fdVal := reflect.Indirect(reflect.ValueOf(l)).FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
