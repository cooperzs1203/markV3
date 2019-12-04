package mnet

import (
	"log"
	"markV3/mface"
	"net"
)

type EntranceFunc func() error

const (
	Serve_Status_UnStarted = iota
	Serve_Status_Load
	Serve_Status_Reload
	Serve_Status_Starting
	Serve_Status_Running
	Serve_Status_Ending
	Serve_Status_Stopped
)

func NewServer() (mface.MServer, error) {
	config := newDefaultConfig()
	return NewServerWithConfig(config)
}

func NewServerWithConfig(config mface.MConfig) (mface.MServer, error) {
	s := &server{
		status:        Serve_Status_UnStarted,
		config:        config,
		connManager:   newConnManager(),
		msgManager:    newMsgManager(),
		routeManager:  newRouteManager(),
		listener:      nil,
		entranceFuncs: make([]EntranceFunc, 0),
	}

	s.connManager.SetServer(s)
	s.msgManager.SetServer(s)
	s.routeManager.SetServer(s)

	return s, s.Load()
}

type server struct {
	status int
	config mface.MConfig

	connManager  mface.MConnManager
	msgManager   mface.MConnManager
	routeManager mface.MRouteManager

	listener *net.Listener

	entranceFuncs []EntranceFunc
}

func (s *server) Load() error {
	s.status = Serve_Status_Load

	// 1. config load
	if err := s.config.Load(); err != nil {
		return err
	}

	// 2. connManager load
	if err := s.connManager.Load(); err != nil {
		return err
	}

	// 3. msgManager load
	if err := s.msgManager.Load(); err != nil {
		return err
	}

	// 4. routeManager load
	if err := s.routeManager.Load(); err != nil {
		return err
	}

	log.Printf("[Server] Load")

	// 5. server load
	if s.config.Host() == "" && s.config.Port() == "" {
		return nil
	}

	host, port := s.config.Host(), s.config.Port()
	if host == "" {
		host = "0.0.0.0"
	}

	if port == "" {
		port = "8888"
	}

	addr := net.JoinHostPort(host, port)
	l, err := net.Listen(s.config.NetType(), addr)
	if err != nil {
		return err
	}

	s.listener = &l

	return nil
}

func (s *server) Start() error {
	s.status = Serve_Status_Starting

	// 1.connManager start
	if err := s.connManager.Start(); err != nil {
		return err
	}

	// 2.msgManager start
	if err := s.msgManager.Start(); err != nil {
		return err
	}

	// 3.routeManager start
	if err := s.routeManager.Start(); err != nil {
		return err
	}

	log.Printf("[Server] Start")

	// 4.accept connection
	go s.startAcceptConnection()

	// 5.run entrance functions
	if len(s.entranceFuncs) > 0 {
		for _, eFunc := range s.entranceFuncs {
			if err := eFunc(); err != nil {
				return err
			}
		}
	}

	s.status = Serve_Status_Running

	// 6. todo:wait for signal
	//for {
	//	time.Sleep(time.Second * time.Duration(5))
	//}

	return nil
}

func (s *server) Stop() error {
	log.Printf("[Server] Stop")
	if s.status >= Serve_Status_Ending {
		return nil
	}

	// 1.start ending
	if err := s.StartEnding(); err != nil {
		return err
	}

	// 2.official ending
	if err := s.OfficialEnding(); err != nil {
		return err
	}

	return nil
}

func (s *server) StartEnding() error {
	log.Printf("[Server] Start Ending")
	s.status = Serve_Status_Ending

	// 1.close accept new connection
	if err := (*s.listener).Close(); err != nil {
		return err
	}

	// 2.notice connManager to start ending
	if err := s.connManager.StartEnding(); err != nil {
		return err
	}

	// 3.notice msgManager to start ending
	if err := s.msgManager.StartEnding(); err != nil {
		return err
	}

	// 4.notice routeManager to start ending
	if err := s.routeManager.StartEnding(); err != nil {
		return err
	}

	return nil
}

func (s *server) OfficialEnding() error {
	// 1.notice routeManager to official ending
	if err := s.routeManager.OfficialEnding(); err != nil {
		return err
	}

	// 2.notice msgManager to official ending
	if err := s.msgManager.OfficialEnding(); err != nil {
		return err
	}

	// 3.notice connManager to official ending
	if err := s.connManager.OfficialEnding(); err != nil {
		return err
	}

	// 4.official ending server
	s.listener = nil
	s.entranceFuncs = s.entranceFuncs[:0]
	s.entranceFuncs = nil

	s.status = Serve_Status_Stopped
	log.Printf("[Server] Official Ending")

	return nil
}

func (s *server) Reload() error {
	log.Printf("[Server] Reload")

	// 1. config reload
	if err := s.config.Reload(); err != nil {
		return err
	}

	// 2. connManager reload
	if err := s.connManager.Reload(); err != nil {
		return err
	}

	// 3. msgManager reload
	if err := s.msgManager.Reload(); err != nil {
		return err
	}

	// 4. routeManager reload
	if err := s.routeManager.Reload(); err != nil {
		return err
	}

	s.status = Serve_Status_Reload

	// 5. server reload
	if s.config.Host() == "" && s.config.Port() == "" {
		s.status = Serve_Status_Running
		return nil
	}

	host, port := s.config.Host(), s.config.Port()
	if host == "" {
		host = "0.0.0.0"
	}

	if port == "" {
		port = "8888"
	}

	addr := net.JoinHostPort(host, port)
	l, err := net.Listen(s.config.NetType(), addr)
	if err != nil {
		return err
	}

	if s.listener != nil {
		_ = (*s.listener).Close()
	}

	s.listener = &l

	s.status = Serve_Status_Running

	return nil
}

func (s *server) ConnManager() mface.MConnManager {
	return s.connManager
}

func (s *server) MsgManager() mface.MMsgManager {
	return s.msgManager
}

func (s *server) RouteManager() mface.MRouteManager {
	return s.routeManager
}

func (s *server) Config() mface.MConfig {
	return s.config
}

func (s *server) RunEntranceFunc(f func() error) {
	if f == nil {
		return
	}

	s.entranceFuncs = append(s.entranceFuncs, f)
}

func (s *server) AddRoute(routeId string, routeHandleFunc func(mface.MMessage, mface.MMessage) error) error {
	route := newRouteHandler(routeId, routeHandleFunc)
	if err := s.routeManager.AddRouteHandle(route); err != nil {
		return err
	}
	return nil
}

func (s *server) AddRoutes(routes map[string]func(mface.MMessage, mface.MMessage) error) error {
	if len(routes) == 0 {
		return nil
	}

	for routeId  , routeHandleFunc := range routes {
		route := newRouteHandler(routeId, routeHandleFunc)
		if err := s.routeManager.AddRouteHandle(route); err != nil {
			return err
		}
	}

	return nil
}

// ======= private functions ========

// start accept connection
func (s *server) startAcceptConnection() {
	if s.listener == nil {
		return
	}

	log.Printf("[Server] %s server running on %s", s.config.Name(), (*s.listener).Addr().String())

	for {
		conn, err := (*s.listener).Accept()
		if err != nil {
			continue
		}

		// todo:handle conn
		log.Println(conn)
	}
}
