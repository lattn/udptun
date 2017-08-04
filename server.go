package udptun

import (
	"net"
	"time"
)

type Server struct {
	config    *Config
	isRunning bool
	stopChan  chan bool
}

func NewServer(configPath string) (*Server, error) {
	config, err := parseConfig(configPath)
	if err != nil {
		return nil, err
	}
	var server Server
	server.config = config
	return &server, nil
}

func (s *Server) Run() error {
	addr, err := net.ResolveUDPAddr("udp", s.config.LocalAddr)
	if err != nil {
		return err
	}

	logger.Printf("try to listen on: \"%s\", allowed IPs: \"%s\"\n", addr, s.config.AllowedIps)

	// 监听
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	logger.Println("try to handle request on:", addr)

	for {
		s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *net.UDPConn) {
	// 读取请求
	request := make([]byte, s.config.BufferSize)
	n, remoteAddr, err := conn.ReadFromUDP(request)
	if err != nil {
		logger.Println("fail to read from UDP:", err)
		return
	}

	logger.Printf("%d bytes read from %s: \"%s\"", n, remoteAddr.String(), string(request))

	// 转发
	if s.allowAddr(remoteAddr) {
		logger.Println("try to forward request")
		go s.forward(conn, remoteAddr, request[:n])
	} else {
		logger.Println("invalid request from ip", remoteAddr)
		go s.sendError(conn, remoteAddr)
	}
}

func (s *Server) allowAddr(addr *net.UDPAddr) bool {
	for _, ip := range s.config.AllowedIps {
		if ip == addr.IP.String() {
			return true
		}
	}
	return false
}

func (s *Server) forward(conn *net.UDPConn, remoteAddr *net.UDPAddr, request []byte) {
	// 转发请求到目标地址
	logger.Println("try to connect target service")
	targetConn, err := net.DialTimeout("udp", s.config.TargetAddr, time.Second*1)
	if err != nil {
		logger.Println("unable to connect target service with error:", err)
		return
	}

	defer targetConn.Close()

	targetConn.SetWriteDeadline(time.Now().Add(time.Second * 1))

	logger.Printf("try to write request %s to target\n", string(request))
	n, err := targetConn.Write(request)
	if err != nil {
		logger.Printf("unable to write request to target: %s\n", err.Error())
		return
	}
	logger.Printf("write %d bytes of request to target\n", n)

	data := make([]byte, s.config.BufferSize)
	n, err = targetConn.Read(data)
	if err != nil {
		logger.Printf("unable to read from target: %s\n", err.Error())
		return
	}

	logger.Printf("read %d bytes from target: \"%s\"\n", n, string(data))
	n, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		logger.Println("fail to write to remote addr:", remoteAddr, string(request), err)
		return
	}
	logger.Println(n, "bytes write to remote addr:", remoteAddr)
}

func (s *Server) sendError(conn *net.UDPConn, remoteAddr *net.UDPAddr) {
	conn.WriteToUDP([]byte(s.config.ErrorMsg), remoteAddr)
}
