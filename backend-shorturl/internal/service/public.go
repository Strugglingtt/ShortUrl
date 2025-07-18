package service

import (
	"backend-shorturl/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "backend-shorturl/api/shorturl/public/v1"
)

type PublicService struct {
	pb.UnimplementedPublicServer
	log *log.Helper
	uc  biz.ShortUrlUsecase
}

func NewPublicService() *PublicService {
	return &PublicService{}
}

func (s *PublicService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{}, nil
}

// CreateShortUrl 这里的service层相当于gin中的Controller层
func (s *PublicService) CreateShortUrl(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenReply, error) {
	//1、 校验请求参数（只做简单校验，复杂逻辑下沉到biz层处理）
	if req.LongUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "long url is required")
	}

	//2、调用biz层处理业务逻辑
	Reply, err := s.uc.CreateShortUrl(ctx, req.LongUrl, req.ExpireTime)
	if err != nil {
		s.log.Errorf("CreateShortUrl error: %v", err)
		return nil, status.Errorf(codes.Internal, "CreateShortUrl error")
	}

	//3、应该是去返回响应，处理返回参数
	return Reply, nil
}
