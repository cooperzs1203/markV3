package mnet

import (
	"log"
	"markV3/mface"
)

func newDefaultConfig() mface.MConfig {
	dc := &defaultConfig{}
	return dc
}

type defaultConfig struct {
	name    string
	netType string
	host    string
	port    string

	cmRequestCS  uint64
	cmResponseCS uint64

	mmRequestCS  uint64
	mmResponseCS uint64

	rmRequestCS  uint64
	rmResponseCS uint64
}

func (dc *defaultConfig) Load() error {
	log.Printf("[DefaultConfig] Load")

	dc.name = "Pay"
	dc.netType = "tcp"
	dc.host = "0.0.0.0"
	dc.port = "8888"

	dc.cmRequestCS = 1000
	dc.cmResponseCS = 1000
	dc.mmRequestCS = 1000
	dc.mmResponseCS = 1000
	dc.rmRequestCS = 1000
	dc.rmResponseCS = 1000

	return nil
}

func (dc *defaultConfig) Reload() error {
	log.Printf("[DefaultConfig] Reload")

	dc.name = "Pay"
	dc.netType = "tcp"
	dc.host = "0.0.0.0"
	dc.port = "8888"

	dc.cmRequestCS = 1000
	dc.cmResponseCS = 1000
	dc.mmRequestCS = 1000
	dc.mmResponseCS = 1000
	dc.rmRequestCS = 1000
	dc.rmResponseCS = 1000

	return nil
}

func (dc *defaultConfig) Name() string {
	return dc.name
}

func (dc *defaultConfig) NetType() string {
	return dc.netType
}

func (dc *defaultConfig) Host() string {
	return dc.host
}

func (dc *defaultConfig) Port() string {
	return dc.port
}

func (dc *defaultConfig) CMRequestCS() uint64 {
	return dc.cmRequestCS
}
func (dc *defaultConfig) CMResponseCS() uint64 {
	return dc.cmResponseCS
}
func (dc *defaultConfig) MMRequestCS() uint64 {
	return dc.mmRequestCS
}
func (dc *defaultConfig) MMResponseCS() uint64 {
	return dc.mmResponseCS
}
func (dc *defaultConfig) RMRequestCS() uint64 {
	return dc.rmRequestCS
}
func (dc *defaultConfig) RMResponseCS() uint64 {
	return dc.rmResponseCS
}
