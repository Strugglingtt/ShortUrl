package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"log"
)

func MyMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				log.Printf("传输协议: %s", tr.Kind())
				log.Printf("操作名称: %s", tr.Operation())
				log.Printf("端点地址: %s", tr.Endpoint())

				// 只处理 HTTP 请求
				if tr.Kind() == transport.KindHTTP {
					if ht, ok := tr.(*http.Transport); ok {
						// 调用 Request() 方法获取 *http.Request
						httpReq := ht.Request()
						log.Printf("HTTP 方法: %s", httpReq.Method)
						log.Printf("请求路径: %s", httpReq.URL.Path)
						log.Printf("客户端 IP: %s", httpReq.RemoteAddr)

						//for k, v := range httpReq.Header {
						//	log.Printf("请求头 %s: %v", k, v)
						//}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
