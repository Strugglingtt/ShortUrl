package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/teris-io/shortid"
)

// ShortUrl is a ShortUrl model.
type ShortUrl struct {
	LongUrl    string
	ShortCode  string
	ExpireTime time.Time
}

type ShortUrlReply struct {
	ShortCode  string `json:"short_code"`
	ShortUrl   string `json:"short_url"`
	LongUrl    string `json:"long_url"`
	CreateTime string `json:"create_time"`
	ExpireTime string `json:"expire_time"`
}

// ShortUrlRepo is a ShortUrl repo.
type ShortUrlRepo interface {
	Save(ctx context.Context, url *ShortUrl) (int64, error)
	//FindByKey(ctx context.Context, key string) (*ShortUrl, error)
}

// ShortUrlUsecase is a ShortUrl usecase.
type ShortUrlUsecase struct {
	repo ShortUrlRepo
	log  *log.Helper
}

// NewShortUrlUsecase new a ShortUrl usecase.
func NewShortUrlUsecase(repo ShortUrlRepo, logger log.Logger) *ShortUrlUsecase {
	return &ShortUrlUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateShortUrl 创建短链接
func (uc *ShortUrlUsecase) CreateShortUrl(ctx context.Context, params *ShortUrl) (*ShortUrlReply, error) {
	// 1. 生成短码
	shortCode, err := shortid.Generate()
	if err != nil {
		uc.log.Errorf("generate short code error: %v", err)
		return nil, err
	}

	// 2. 构造实体
	now := time.Now()
	shortUrl := &ShortUrl{
		LongUrl:    params.LongUrl,
		ShortCode:  shortCode,
		ExpireTime: params.ExpireTime,
	}

	// 3. 存储到数据库
	urlInfo, err := uc.repo.Save(ctx, shortUrl)
	if err != nil {
		uc.log.Errorf("save short url error: %v", err)
		return nil, err
	}

	fmt.Printf("ShortUrl Info = %v", urlInfo)
	// 4. 构造返回结果
	return &ShortUrlReply{
		ShortCode:  shortCode,
		ShortUrl:   "http://localhost:8081/" + shortCode,
		LongUrl:    params.LongUrl,
		CreateTime: now.Format(time.RFC3339),
		ExpireTime: params.ExpireTime.Format(time.RFC3339),
	}, nil
}
