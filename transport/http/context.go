package http

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

var _ Context = (*wrapper)(nil)

// http服务数据处理接口
// 请求参数的获取以及对应的响应
type Context interface {
	context.Context
	// url 获取参数值，针对不同的请求不同的参数格式
	Vars() url.Values
	Query() url.Values
	Form() url.Values

	// 请求相关
	Header() http.Header
	Request() *http.Request
	Response() http.ResponseWriter

	//参数绑定
	// Bind(interface{}) error
	// BindVars(interface{}) error
	// BindQuery(interface{}) error
	// BindForm(interface{}) error

	//

}

// 封装http.ResponseWriter
type responseWriter struct {
	code int
	w    http.ResponseWriter
}

func (w *responseWriter) reset(res http.ResponseWriter) {
	w.w = res
	w.code = http.StatusOK
}

func (w *responseWriter) Header() http.Header { return w.w.Header() }

func (w *responseWriter) WriterHeader(statuscode int) { w.code = statuscode }

func (w *responseWriter) Write(data []byte) (int error) {
	w.w.WriteHeader(w.code)
	return w.Write(data)
}

type wrapper struct {
	router *Router
	req    *http.Request
	// 这两个
	res http.ResponseWriter
	w   responseWriter
}

func (c *wrapper) Header() http.Header {
	return c.req.Header
}

func (c *wrapper) Vars() url.Values {
	raws := mux.Vars(c.req)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return vars
}

func (c *wrapper) Form() url.Values {
	if err := c.req.ParseForm(); err != nil {
		return url.Values{}
	}

	return c.req.Form
}

func (c *wrapper) Query() url.Values {
	return c.req.URL.Query()
}

func (c *wrapper) Request() *http.Request        { return c.req }
func (c *wrapper) Response() http.ResponseWriter { return c.res }

func (c *wrapper) Deadline() (time.Time, bool) {
	if c.req == nil {
		return time.Time{}, false
	}
	return c.req.Context().Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Done()
}

func (c *wrapper) Err() error {
	if c.req == nil {
		return context.Canceled
	}
	return c.req.Context().Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Value(key)
}
