package service

import (
	"go-canal/libs/canal"
	"go-canal/libs/log"
)

type TService struct {
	Canal *canal.TCanal
}

var Service *TService

func (s *TService) Run() {
	log.Infof("Service start running")
	go s.Canal.Run()
}

func (s *TService) Close() {
	log.Warnf("Service stopping")
	s.Canal.Close()
	log.Warnf("Service stopped")
}

func Init() {
	Service = &TService{
		Canal: canal.Canal,
	}
}

func Run() {
	Service.Run()
}

func Stop() {
	Service.Close()
}
