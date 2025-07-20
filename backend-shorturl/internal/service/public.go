package service

import (
	"backend-shorturl/internal/biz"
	"context"
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
