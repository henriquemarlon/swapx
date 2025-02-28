package gio

import "errors"

type GioHandlerFactory interface {
	NewGioHandler(domain uint16) (GioHandler, error)
}

type DefaultGioHandlerFactory struct {
	BaseUrl string
}

func NewGioHandlerFactory(baseUrl string) GioHandlerFactory {
	return &DefaultGioHandlerFactory{
		BaseUrl: baseUrl,
	}
}

func (f *DefaultGioHandlerFactory) NewGioHandler(domain uint16) (GioHandler, error) {
	switch domain {
	case 0x27:
		return NewGioGetStorage(f.BaseUrl, domain), nil
	// Adicione novos cases aqui conforme necessário
	default:
		return nil, errors.New("domain not supported")
	}
}
