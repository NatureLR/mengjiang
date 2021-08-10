package httpsocks

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"syscall"

	"golang.org/x/net/proxy"
)

type server struct {
	HttpAddr  string
	SocksAddr string
	Mode      string
	Conn      net.Conn
}

func New(http, socks, mode string) *server {
	return &server{
		HttpAddr:  http,
		SocksAddr: socks,
		Mode:      mode,
	}
}

func (s *server) Run() {
	log.Printf("http监听地址(%s),socks监听地址(%s),代理模式(%s)", s.HttpAddr, s.SocksAddr, s.Mode)

	l, err := net.Listen("tcp", s.HttpAddr)
	if err != nil {
		log.Println("监听失败：", err)
	}

	for {
		s.Conn, err = l.Accept()
		if err != nil {
			log.Println("接受数据错误:", err)
			continue
		}
		go func() {
			s.handleConnection()
		}()
	}
}

func (s *server) handleConnection() {
	defer func() {
		s.Conn.Close()
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	sc, ok := s.Conn.(*net.TCPConn)
	if !ok {
		log.Println("请求的不是tcp")
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	destAddr, err := s.getDestAddr(sc, buf)
	checkErr(err)

	dialer, err := proxy.SOCKS5("tcp", s.SocksAddr, nil, proxy.Direct)
	checkErr(err)
	dc, err := dialer.Dial("tcp", destAddr)
	checkErr(err)

	if s.Mode == "forward" {
		x := make([]byte, 1024)
		n, err := buf.Read(x)
		checkErr(err)
		_, err = dc.Write(x[:n])
		checkErr(err)
	}

	log.Printf("%s<--->(%s,%s)<--->%s<--->%s\n", sc.RemoteAddr().String(), sc.LocalAddr(), dc.LocalAddr().String(), dc.RemoteAddr().String(), destAddr)
	forward(dc, sc)
}

// 获取源目的地址,nat模式和非nat不太一样
func (s *server) getDestAddr(c *net.TCPConn, buf *bytes.Buffer) (addr string, err error) {
	switch s.Mode {
	case "nat":
		addr, err = getOrigAddr(c)
		fmt.Println("====================", addr)
	case "forward":
		// 从tcp中拿到http中的host字段
		hq, err := http.ReadRequest(bufio.NewReader(io.TeeReader(c, buf)))
		checkErr(err)
		if hq.URL.Port() == "" {
			return hq.Host + ":80", err
		}
		return hq.Host, err
	default:
		log.Println("不支持的模式")
	}
	return
}

func forward(dst, src net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	f := func(dst, src net.Conn) {
		defer wg.Done()
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Println(err)
		}
	}

	go f(dst, src)
	go f(src, dst)
	wg.Wait()
}

const SO_ORIGINAL_DST = 80

// 获取源目的地址
// 获取自 https://gist.github.com/time-river/210c730a66f5bf62b1fcc3cfc163335c
func getOrigAddr(c *net.TCPConn) (string, error) {
	f, err := c.File()
	if err != nil {
		log.Println(err)
	}
	addr, err := syscall.GetsockoptIPv6Mreq(int(f.Fd()), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		log.Println("syscall.GetsockoptIPv6Mreq error: %w", err)
		return "", err
	}

	remote := fmt.Sprintf("%d.%d.%d.%d:%d",
		addr.Multiaddr[4], addr.Multiaddr[5], addr.Multiaddr[6], addr.Multiaddr[7],
		uint16(addr.Multiaddr[2])<<8+uint16(addr.Multiaddr[3]))
	return remote, nil
}
