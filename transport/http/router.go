package http

// 路由的具体信息  方法+url
type RouterInfo struct {
	Path   string
	Method string
}

// Context 将所有请求以及响应操作通过Context完成
type HandleFunc func(Context) error

//该结构只是简单的封装，暴露给用户的注册路由结构
type Router struct {
	prefix string
	srv    *Server
}

func newRouter(prefix string) *Router {
	r := &Router{
		prefix: prefix,
	}

	return r
}

func (r *Router) Handle(method, relativePath string, h HandleFunc) {

}
