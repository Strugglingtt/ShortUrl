package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"log"
)

func RedirectMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if tr.Kind() == transport.KindHTTP {
					if ht, ok := tr.(*http.Transport); ok {
						httpReq := ht.Request()
						log.Printf("URL: %s", httpReq.URL.Path)
						if httpReq.URL.Path == "/api/123" {
							http.NewRedirect("https://baidu.com", 301)
							return handler(ctx, req)
						}
					}
				}
			}
			return handler(ctx, req)

		}
	}
}
