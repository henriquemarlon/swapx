package service

import "errors"

type GioHandlerFactory interface {
	NewGioHandler(domain uint16) (GioHandler, error)
}

type DefaultGioHandlerFactory struct {
	BaseUrl string `json:"base_url"`
}

func NewGioHandlerFactory(baseUrl string) *DefaultGioHandlerFactory {
	return &DefaultGioHandlerFactory{
		BaseUrl: baseUrl,
	}
}

func (f *DefaultGioHandlerFactory) NewGioHandler(domain uint16) (GioHandler, error) {
	switch domain {
	case 0x27:
		return &GioGetStorage{
			Domain:  domain,
			BaseUrl: f.BaseUrl,
		}, nil
	// more cases according with co-processor domains
	default:
		return nil, errors.New("domain not supported")
	}
}