package service

import (
	"backend-shorturl/internal/biz"
	"context"
	"fmt"
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
		return nil, err
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
	//获取原长连接
	fmt.Println("req.Code:", req.Code)
	url, err := s.uc.GetOriginalURL(ctx, req.Code)
	if err != nil {
		log.Info("长连接获取失败")
		return nil, err
	}
	//实际这里还要增加短链的数据，对短链的一个点击详情进行统计

	// 获取HTTP传输器 进行一个重定向跳转
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
		302, //302和308   == 308是永久重定向，浏览器会缓存
		"REDIRECT",
		"Resource has moved",
	).WithMetadata(map[string]string{
		"Location": url,
	})
}

// GetStatics 短链数据统计信息获取
func (s *PublicService) GetStatics(ctx context.Context, req *pb.GetStaticsRequest) (*pb.GetStaticsReply, error) {
	//如果是短链数据的获取的话，首先应该是对短链参数是否传递进行校验
	if req.GetShortCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "short code is required")
	}
	//调用biz层进行业务处理，实际也是查询相关参数，这里简化查询，只对数据进行一个简单的处理，颗粒度暂不考虑
	statis, err := s.uc.GetStaticsInfo(ctx, req.GetShortCode())
	if err != nil {
		return nil, err
	}

	data := pb.GetStaticsReply_Data{}
	data.TotalClicks = statis.TotalClicks
	data.ShortCode = statis.ShortCode
	data.OriginalUrl = statis.OriginUrl
	return &pb.GetStaticsReply{
		Code:    "200",
		Data:    &data,
		Message: "success",
	}, nil
}

func (s *PublicService) GetAllStatics(ctx context.Context, req *pb.GetAllStaticsRequest) (*pb.GetAllStaticsReply, error) {
	Allstatis, total_page, totalCount, err := s.uc.GetAllStaticsInfo(ctx, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	data := make([]*pb.GetAllStaticsReply_Data, 0)
	for _, statis := range Allstatis {
		data = append(data, &pb.GetAllStaticsReply_Data{
			ShortCode:   statis.ShortCode,
			OriginalUrl: statis.OriginUrl,
			TotalClicks: statis.TotalClicks,
		})
	}
	return &pb.GetAllStaticsReply{
		Code:       "200",
		Data:       data,
		Message:    "success",
		Total:      int32(totalCount),
		TotalPages: int32(total_page),
		Page:       req.GetPage(),
	}, nil
}
