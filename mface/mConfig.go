package mface

type MConfig interface {
	Load() error   // load config
	Reload() error // reload config

	Name() string            // name of server
	NetType() string         // net type of server
	Host() string            // host of server
	Port() string            // port of server
	ConnReadTimeOut() uint64 // time out of connection conn read
	ConnResponseCS() uint64  // connection response channel size
	DPCompletedCS() uint64   // data protocol completed message channel size
	CMMaxConnNumber() uint64 // max connections of connManager
	CMRequestCS() uint64     // connManager request channel size
	CMResponseCS() uint64    // connManager response channel size
	MMRequestCS() uint64     // msgManager request channel size
	MMResponseCS() uint64    // msgManager response channel size
	RMRequestCS() uint64     // routeManager request channel size
	RMResponseCS() uint64    // routeManager response channel size
}
