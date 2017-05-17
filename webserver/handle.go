package webserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
	StatusText string
	StatusCode int
	Proto      string
	Header     map[string][]string
	Content    []byte
}

var statusMap = map[int]string{
	200: "OK",
	400: "BAD REQUEST",
	404: "NOT FOUND",
	500: "INNER ERROR",
}

func statusTextFromCode(code int) string {
	text, ok := statusMap[code]
	if !ok {
		return "ERROR"
	}
	return text
}

// FileResp 创建返回文件的响应体
func FileResp(path, contentType string) (*Resp, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	resp := NewResp(200, content)
	resp.Header["Content-Type"] = []string{contentType}
	return resp, nil
}

// HTMLResp 创建返回网页的响应体
func HTMLResp(path string) (*Resp, error) {
	return FileResp(path, "text/html")
}

// ErrResp 创建一个新的响应体
func ErrResp(statusCode int) *Resp {
	return NewResp(statusCode, []byte(fmt.Sprintf(
		"<html><head><title>%d</title></head><body>%s</body></html>",
		statusCode, statusTextFromCode(statusCode),
	)))
}

// NewResp 创建一个新的响应体
func NewResp(statusCode int, content []byte) *Resp {
	resp := Resp{
		StatusCode: statusCode,
		StatusText: statusTextFromCode(statusCode),
		Proto:      "HTTP/1.1",
		Header: map[string][]string{
			"Server":         []string{"Zheng's computer-network web server, v0,1."},
			"Date":           []string{time.Now().UTC().Format("Mon, 2 Jan 2006 15:04:05 GMT")},
			"Content-Length": []string{strconv.Itoa(len(content))},
		},
		Content: content,
	}
	return &resp
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

// Static 注册静态文件夹路由
func (s *Server) Static(route, path string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	s.staticMap[route] = dir + path
	return nil
}

// GetStaticHandle 获取静态文件夹路由控制器
func GetStaticHandle(path string) Handle {
	return func(req *Req) *Resp {
		var (
			resp *Resp
			err  error
		)
		suffix := path[strings.LastIndex(path, ".")+1:]
		switch suffix {
		case "html":
			resp, err = HTMLResp(path)
		case "js":
			fallthrough
		case "es":
			resp, err = FileResp(path, "text/javascript")
		default:
			resp, err = FileResp(path, "text/"+suffix)
		}
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				return ErrResp(404)
			}
			log.Println(err)
			return ErrResp(400)
		}
		return resp
	}
}
