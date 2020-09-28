package connmgr

import (
	"encoding/binary"
	"errors"
	"github.com/stalker-loki/app/slog"
	"net"
)

const (
	torSucceeded         = 0x00
	torGeneralError      = 0x01
	torNotAllowed        = 0x02
	torNetUnreachable    = 0x03
	torHostUnreachable   = 0x04
	torConnectionRefused = 0x05
	torTTLExpired        = 0x06
	torCmdNotSupported   = 0x07
	torAddrNotSupported  = 0x08
)

var (
	// ErrTorInvalidAddressResponse indicates an invalid address was returned by the Tor DNS resolver.
	ErrTorInvalidAddressResponse = errors.New("invalid address response")
	// ErrTorInvalidProxyResponse indicates the Tor proxy returned a response in an unexpected format.
	ErrTorInvalidProxyResponse = errors.New("invalid proxy response")
	// ErrTorUnrecognizedAuthMethod indicates the authentication method provided is not recognized.
	ErrTorUnrecognizedAuthMethod = errors.New("invalid proxy authentication method")
	torStatusErrors              = map[byte]error{
		torSucceeded:         errors.New("tor succeeded"),
		torGeneralError:      errors.New("tor general error"),
		torNotAllowed:        errors.New("tor not allowed"),
		torNetUnreachable:    errors.New("tor network is unreachable"),
		torHostUnreachable:   errors.New("tor host is unreachable"),
		torConnectionRefused: errors.New("tor connection refused"),
		torTTLExpired:        errors.New("tor TTL expired"),
		torCmdNotSupported:   errors.New("tor command not supported"),
		torAddrNotSupported:  errors.New("tor address type not supported"),
	}
)

// TorLookupIP uses Tor to resolve DNS via the SOCKS extension they provide for resolution over the Tor network. Tor itself doesn't support ipv6 so this doesn't either.
func TorLookupIP(host, proxy string) (addr []net.IP, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", proxy); slog.Check(err) {
		return
	}
	defer conn.Close()
	buf := []byte{'\x05', '\x01', '\x00'}
	if _, err = conn.Write(buf); slog.Check(err) {
		return
	}
	buf = make([]byte, 2)
	if _, err = conn.Read(buf); slog.Check(err) {
		return
	}
	if buf[0] != '\x05' {
		err = ErrTorInvalidProxyResponse
		return
	}
	if buf[1] != '\x00' {
		err = ErrTorUnrecognizedAuthMethod
		return
	}
	buf = make([]byte, 7+len(host))
	buf[0] = 5      // protocol version
	buf[1] = '\xF0' // Tor Resolve
	buf[2] = 0      // reserved
	buf[3] = 3      // Tor Resolve
	buf[4] = byte(len(host))
	copy(buf[5:], host)
	buf[5+len(host)] = 0 // Port 0
	if _, err = conn.Write(buf); slog.Check(err) {
		return
	}
	buf = make([]byte, 4)
	if _, err = conn.Read(buf); slog.Check(err) {
		return
	}
	if buf[0] != 5 {
		err = ErrTorInvalidProxyResponse
		return
	}
	if buf[1] != 0 {
		if int(buf[1]) >= len(torStatusErrors) {
			err = ErrTorInvalidProxyResponse
			slog.Debug(err)
			return
		} else if err = torStatusErrors[buf[1]]; slog.Check(err) {
			return
		}
		err = ErrTorInvalidProxyResponse
		slog.Debug(err)
		return
	}
	if buf[3] != 1 {
		err = torStatusErrors[torGeneralError]
		slog.Debug(err)
		return
	}
	buf = make([]byte, 4)
	var bytes int
	if bytes, err = conn.Read(buf); slog.Check(err) {
		return
	}
	if bytes != 4 {
		err = ErrTorInvalidAddressResponse
		return
	}
	r := binary.BigEndian.Uint32(buf)
	addr = make([]net.IP, 1)
	addr[0] = net.IPv4(byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
	return
}
