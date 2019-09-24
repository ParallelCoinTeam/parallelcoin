package rudp

import (
	"net"
	"time"
)

func NewConn(conn *net.UDPConn, rudp *Rudp) *RudpConn {
	con := &RudpConn{conn: conn, rudp: rudp,
		recvChan: make(chan []byte, 1<<16), recvErr: make(chan error, 2),
		sendChan: make(chan []byte, 1<<16), sendErr: make(chan error, 2),
		SendTick: make(chan int, 2),
	}
	go con.run()
	return con
}

func NewUnConn(conn *net.UDPConn, remoteAddr *net.UDPAddr, rudp *Rudp, close func(string)) *RudpConn {
	con := &RudpConn{conn: conn, rudp: rudp, SendTick: make(chan int, 2),
		recvChan: make(chan []byte, 1<<16), recvErr: make(chan error, 2),
		sendChan: make(chan []byte, 1<<16), sendErr: make(chan error, 2),
		closef: close, remoteAddr: remoteAddr, in: make(chan []byte, 1<<16),
	}
	go con.run()
	return con
}

type RudpConn struct {
	conn *net.UDPConn

	rudp *Rudp

	recvChan chan []byte
	recvErr  chan error

	sendChan chan []byte
	sendErr  chan error

	SendTick chan int

	//unconected
	remoteAddr *net.UDPAddr
	closef     func(addr string)
	in         chan []byte
}

func (rc *RudpConn) SetDeadline(t time.Time) error      { return nil }
func (rc *RudpConn) SetReadDeadline(t time.Time) error  { return nil }
func (rc *RudpConn) SetWriteDeadline(t time.Time) error { return nil }
func (rc *RudpConn) LocalAddr() net.Addr                { return rc.conn.LocalAddr() }
func (rc *RudpConn) Connected() bool                    { return rc.remoteAddr == nil }
func (rc *RudpConn) RemoteAddr() net.Addr {
	if rc.remoteAddr != nil {
		return rc.remoteAddr
	}
	return rc.conn.RemoteAddr()
}
func (rc *RudpConn) Close() error {
	var err error
	if rc.remoteAddr != nil {
		if rc.closef != nil {
			rc.closef(rc.remoteAddr.String())
		}
		_, err = rc.conn.WriteToUDP([]byte{TYPE_CORRUPT}, rc.remoteAddr)
		rc.in <- []byte{TYPE_EOF}
	} else {
		_, err = rc.conn.Write([]byte{TYPE_CORRUPT})
	}
	checkErr(err)
	return err
}
func (rc *RudpConn) Read(bts []byte) (n int, err error) {
	select {
	case data := <-rc.recvChan:
		copy(bts, data)
		return len(data), nil
	case err := <-rc.recvErr:
		return 0, err
	}
}

func (r *RudpConn) send(bts []byte) (err error) {
	select {
	case r.sendChan <- bts:
		return nil
	case err := <-r.sendErr:
		return err
	}
}
func (r *RudpConn) Write(bts []byte) (n int, err error) {
	sz := len(bts)
	for len(bts)+MAX_MSG_HEAD > GENERAL_PACKAGE {
		if err := r.send(bts[:GENERAL_PACKAGE-MAX_MSG_HEAD]); err != nil {
			return 0, err
		}
		bts = bts[GENERAL_PACKAGE-MAX_MSG_HEAD:]
	}
	return sz, r.send(bts)
}

func (r *RudpConn) rudpRecv(data []byte) error {
	for {
		n, err := r.rudp.Recv(data)
		if err != nil {
			r.recvErr <- err
			return err
		} else if n == 0 {
			break
		}
		bts := make([]byte, n)
		dataRead := data[:n]
		copy(bts, dataRead)
		r.recvChan <- bts
	}
	return nil
}
func (r *RudpConn) conectedRecvLoop() {
	data := make([]byte, MAX_PACKAGE)
	for {
		n, err := r.conn.Read(data)
		if err != nil {
			r.recvErr <- err
			return
		}
		dataRead := data[:n]
		r.rudp.Input(dataRead)
		if r.rudpRecv(data) != nil {
			return
		}
	}
}
func (r *RudpConn) unconectedRecvLoop() {
	data := make([]byte, MAX_PACKAGE)
	for {
		select {
		case bts := <-r.in:
			r.rudp.Input(bts)
			if r.rudpRecv(data) != nil {
				return
			}
		}
	}
}
func (r *RudpConn) sendLoop() {
	var sendNum int
	for {
		select {
		case tick := <-r.SendTick:
		sendOut:
			for {
				select {
				case bts := <-r.sendChan:
					_, err := r.rudp.Send(bts)
					if err != nil {
						r.sendErr <- err
						return
					}
					sendNum++
					if sendNum >= maxSendNumPerTick {
						break sendOut
					}
				default:
					break sendOut
				}
			}
			sendNum = 0
			p := r.rudp.Update(tick)
			var num, sz int
			for p != nil {
				n, err := int(0), error(nil)
				if r.Connected() {
					n, err = r.conn.Write(p.Bts)
				} else {
					n, err = r.conn.WriteToUDP(p.Bts, r.remoteAddr)
				}
				if err != nil {
					r.sendErr <- err
					return
				}
				sz, num = sz+n, num+1
				p = p.Next
			}
			if num > 1 {
				show := bitShow(sz * int(time.Second/sendTick))
				dbg("send package num %v,sz %v, %v/s,local %v,remote %v",
					num, show, show, r.LocalAddr(), r.RemoteAddr())
			}
		}
	}
}
func (r *RudpConn) run() {
	if autoSend && sendTick > 0 {
		go func() {
			tick := time.Tick(sendTick)
			for {
				select {
				case <-tick:
					r.SendTick <- 1
				}
			}
		}()
	}
	go func() {
		if r.Connected() {
			r.conectedRecvLoop()
		} else {
			r.unconectedRecvLoop()
		}
	}()
	r.sendLoop()
}
