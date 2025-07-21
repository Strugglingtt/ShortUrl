package service

import (
	"backend-shorturl/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "backend-shorturl/api/shorturl/public/v1"
)

type PublicService struct {
	pb.UnimplementedPublicServer
	log *log.Helper
	uc  *biz.ShortUrlUsecase
}

// CreateShortUrl 创建短链接
func (s *PublicService) CreateShortUrl(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenReply, error) {
	// 1. 校验请求参数
	if req.LongUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "long url is required")
	}

	// 2. 解析过期时间
	var expireTime time.Time
	var err error
	if req.ExpireTime != "" {
		expireTime, err = time.Parse(time.RFC3339, req.ExpireTime)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid expire time format, should be RFC3339")
		}
	}

	// 3. 调用biz层处理业务逻辑
	reply, err := s.uc.CreateShortUrl(ctx, &biz.ShortUrl{
		LongUrl:    req.LongUrl,
		ExpireTime: expireTime,
	})
	if err != nil {
		s.log.Errorf("CreateShortUrl error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create short url")
	}

	// 4. 构造响应
	return &pb.ShortenReply{
		Code:    "200",
		Message: "success",
		Data: &pb.Data{
			ShortCode:  reply.ShortCode,
			ShortUrl:   reply.ShortUrl,
			LongUrl:    reply.LongUrl,
			CreateAt:   reply.CreateTime,
			ExpireTime: reply.ExpireTime,
		},
	}, nil
}

// Redirect 短链接重定向跳转
func (s *PublicService) Redirect(ctx context.Context, req *pb.RedirectRequest) (*pb.RedirectReply, error) {
	if req.Code == "" {
		return nil, errors.BadRequest("REDIRECT", "code is required")
	}

	url, err := s.uc.GetOriginalURL(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	// 获取HTTP传输器
	if tr, ok := transport.FromServerContext(ctx); ok && tr.Kind() == transport.KindHTTP {
		if ht, ok := tr.(*http.Transport); ok {
			// 设置重定向头
			ht.ReplyHeader().Set("Location", url)
			// 设置缓存控制头
			ht.ReplyHeader().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			ht.ReplyHeader().Set("Pragma", "no-cache")
			ht.ReplyHeader().Set("Expires", "0")

			// 使用方式
			return nil, RedirectError(url)
		}
	}

	// 非HTTP请求或无法重定向时返回原始URL
	return &pb.RedirectReply{LongUrl: url}, nil
}
func RedirectError(url string) error {
	return errors.New(
		302,
		"REDIRECT",
		"Resource has moved",
	).WithMetadata(map[string]string{
		"Location": url,
	})
}
