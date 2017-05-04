package webserver

import (
	"errors"
	"net/url"
)

// Req 请求体
type Req struct {
	Method  string
	URL     *url.URL
	Proto   string
	Header  map[string][]string
	Content []byte
}

// Resp 响应体
type Resp struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	Header     map[string][]string
	Content    []byte
}

// Handle 控制器
type Handle func(*Req) *Resp

func (s *Server) reg(route string, handle Handle) error {
	if _, ok := s.handleMap[route]; ok {
		return errors.New("Route has already registered.")
	}
	s.handleMap[route] = handle
	return nil
}

// Reg 注册路由控制器
func (s *Server) Reg(route string, handle Handle, methods []string) error {
	return s.reg(route, func(req *Req) *Resp {
		for _, m := range methods {
			if m == req.Method {
				return handle(req)
			}
		}
		return nil
	})
}

// Get 注册GET方法的路由控制器
func (s *Server) Get(route string, handle Handle) error {
	return s.reg(route, func(req *Req) *Resp {
		if req.Method == "GET" {
			return handle(req)
		}
		return nil
	})
}

// Post 注册POST方法的路由控制器
func (s *Server) Post(route string, handle Handle) error {
	return s.reg(route, func(req *Req) *Resp {
		if req.Method == "POST" {
			return handle(req)
		}
		return nil
	})
}

// Any 注册不限方法的路由控制器
func (s *Server) Any(route string, handle Handle) error {
	return s.reg(route, handle)
}
