package biz

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/teris-io/shortid"
)

// ShortUrl is a ShortUrl model.
type ShortUrl struct {
	LongUrl     string
	ShortCode   string
	ExpireTime  time.Time
	AccessCount int32
	ClientIp    string
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
	Save(ctx context.Context, url *ShortUrl) (*ShortUrl, error)
	GetOriginalURL(ctx context.Context, code string) (*ShortUrl, error)
	IsExit(ctx context.Context, url *ShortUrl) (bool, error)
	IsExitByCode(ctx context.Context, code string) (bool, error)
	GetShortStaticsInfo(ctx context.Context, code string) (*ShortStaticsInfo, error)
	GetAllShorStaticsInfos(ctx context.Context, page, size int) ([]*ShortStaticsInfo, int, int, error)
	SetCache(ctx context.Context, key string, shortUrl *ShortUrl) error
	GetCache(ctx context.Context, key string) (*ShortUrl, error)
	AddCache(ctx context.Context, key string) error
	SetEmptyCache(ctx context.Context, key string) error
	IncrementAccessCount(ctx context.Context, shortCode string) error
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

// GenerateShortCode 带递归的重试机制
func (uc *ShortUrlUsecase) GenerateShortCode(ctx context.Context, n int) (string, error) {
	if n > 5 {
		return "", errors.New("重试5次失败")
	}

	// 1. 生成短码
	shortCode, err := shortid.Generate()
	if err != nil {
		uc.log.Errorf("generate short code error: %v", err)
		return "", err
	}

	// 检查短码是否已存在（注意：这里isExist表示"已存在"，而非"可用"）
	isExist, err := uc.repo.IsExitByCode(ctx, shortCode)
	if err != nil {
		uc.log.Errorf("IsExistByCode error: %v", err)
		return "", err
	}

	// 如果已存在，递归重试
	if isExist {
		return uc.GenerateShortCode(ctx, n+1)
	}

	// 短码不存在，可用
	return shortCode, nil
}

// CreateShortUrl 创建短链接
func (uc *ShortUrlUsecase) CreateShortUrl(ctx context.Context, params *ShortUrl) (*ShortUrlReply, error) {
	//// 1. 生成短码
	//shortCode, err := shortid.Generate()
	//if err != nil {
	//	uc.log.Errorf("generate short code error: %v", err)
	//	return nil, err
	//}
	shortCode, err := uc.GenerateShortCode(ctx, 1)
	if err != nil {
		uc.log.Errorf("GenerateShortCode error: %v", err)
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
	const CachePrefix = "shortUrl:"

	key := CachePrefix + urlInfo.ShortCode
	//存入缓存
	if err := uc.repo.SetCache(ctx, key, urlInfo); err != nil {
		uc.log.Errorf("cache short url error: %v", err)
	}

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
// GetOriginalURL 获取原长链接，带缓存处理
func (uc *ShortUrlUsecase) GetOriginalURL(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.New("code is empty")
	}

	const CachePrefix = "shortUrl:"
	key := CachePrefix + code

	// 1. 优先从 Redis 获取

	shortUrlInfo, err := uc.repo.GetCache(ctx, key)
	if err != nil {
		// 处理 Redis 错误（非未命中）
		if !errors.Is(err, redis.Nil) {
			uc.log.Errorf("GetCache error: %v", err)
			return "", err
		}

		// ====== Redis 未命中（err == redis.Nil），进入数据库查询分支 ======
		// 2. 查询数据库（注意：这里应调用数据库查询方法，而非 GetCache）
		log.Info("code is : " + code)
		shortUrl, err := uc.repo.GetOriginalURL(ctx, code) // 修正：调用数据库查询方法
		if err != nil {
			uc.log.Errorf("GetOriginalURL error: %v", err)
			return "", err
		}

		// 3. 数据库也未找到（shortUrl 为 nil）
		if shortUrl == nil {
			// 设置空值缓存，防止缓存穿透
			_ = uc.repo.SetEmptyCache(ctx, key)
			return "", errors.New("short code not found")
		}

		// 4. 数据库找到，写入 Redis 缓存
		if err := uc.repo.SetCache(ctx, key, shortUrl); // 修正：传入 key 和 shortUrl
		err != nil {
			uc.log.Warnf("Failed to set cache: %v", err)
			// 缓存失败不影响主流程
		}

		// 5. 异步更新访问计数
		go func() {
			ctxBg := context.Background()
			if err := uc.repo.IncrementAccessCount(ctxBg, code); err != nil {
				uc.log.Warnf("Failed to increment access count: %v", err)
			}
		}()

		return shortUrl.LongUrl, nil
	}

	// ====== Redis 命中（err == nil），直接返回缓存数据 ======
	// 异步更新访问计数
	go func() {
		ctxBg := context.Background()
		if err := uc.repo.IncrementAccessCount(ctxBg, code); err != nil {
			uc.log.Warnf("Failed to increment access count: %v", err)
		}
	}()

	// 检查缓存数据是否为 nil（防止异常情况）
	if shortUrlInfo == nil {
		return "", errors.New("short code not found in cache")
	}

	return shortUrlInfo.LongUrl, nil
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
