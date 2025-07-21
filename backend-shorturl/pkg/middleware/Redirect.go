package middleware

import (
	"errors"
	"net/http"
)

// RedirectMiddleware 短链接重定向中间件
func RedirectMiddleware(svc RedirectServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 获取短码
			code := getCodeFromRequest(r) // 需要根据路由实现

			// 如果是非短链接请求，直接放行
			if code == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 2. 查询长链接
			url, err := svc.GetOriginalURL(r.Context(), code)
			if err != nil {
				// 查询失败，可以记录日志或返回404
				if errors.Is(err, ErrNotFound) {
					http.Error(w, "Short URL not found", http.StatusNotFound)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// 3. 执行重定向
			w.Header().Set("Location", url)
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusFound)
		})
	}
}
