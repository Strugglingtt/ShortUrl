package biz

import (
	"context"
	"errors"
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
type ShortStaticsInfo struct {
	ShortCode   string `json:"short_code"`
	OriginUrl   string `json:"origin_url"`
	TotalClicks uint32 `json:"total_clicks"`
}

// ShortUrlRepo is a ShortUrl repo.
type ShortUrlRepo interface {
	Save(ctx context.Context, url *ShortUrl) (int64, error)
	GetOriginalURL(ctx context.Context, code string) (string, error)
	IsExit(ctx context.Context, url *ShortUrl) (bool, error)
	IsExitByCode(ctx context.Context, code string) (bool, error)
	GetShortStaticsInfo(ctx context.Context, code string) (*ShortStaticsInfo, error)
	GetAllShorStaticsInfos(ctx context.Context, page, size int) ([]*ShortStaticsInfo, int, int, error)
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

	//判断这个长链接是否已经存储过了，如果存储过了就不用再写入数据库了对吧
	Exit, err := uc.repo.IsExit(ctx, shortUrl)
	if err != nil {
		uc.log.Errorf("check short code error: %v", err)
	}
	if Exit {
		uc.log.Errorf("short code is exit")
		return nil, errors.New("short code is exit")
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
		ShortUrl:   "api/" + shortCode,
		LongUrl:    params.LongUrl,
		CreateTime: now.Format(time.RFC3339),
		ExpireTime: params.ExpireTime.Format(time.RFC3339),
	}, nil
}

// GetOriginalURL 获取原长链接，这里应该加一个短链数据的处理
func (uc *ShortUrlUsecase) GetOriginalURL(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.New("code is empty")
	}

	return uc.repo.GetOriginalURL(ctx, code)
}

// GetStaticsInfo 获取短链统计的数据信息
func (uc *ShortUrlUsecase) GetStaticsInfo(ctx context.Context, code string) (*ShortStaticsInfo, error) {
	//第一步是去数据查询是否存在该code
	exit, err := uc.repo.IsExitByCode(ctx, code)
	if err != nil {
		uc.log.Errorf("check short code error: %v", err)
	}
	if !exit {
		uc.log.Errorf("short code is not exit")
	}
	//第二步才是去获取该code对应的数据并返回
	statics, err := uc.repo.GetShortStaticsInfo(ctx, code)
	if err != nil {
		uc.log.Errorf("get short code error: %v", err)
	}
	return statics, nil
}

// GetAllStaticsInfo 获取所有的短链统计信息
func (uc *ShortUrlUsecase) GetAllStaticsInfo(ctx context.Context, page, size int) ([]*ShortStaticsInfo, int, int, error) {
	AllStatis, total_page, totalCount, err := uc.repo.GetAllShorStaticsInfos(ctx, page, size)
	if err != nil {
		uc.log.Errorf("get all short code error: %v", err)
	}
	return AllStatis, total_page, totalCount, nil
}
