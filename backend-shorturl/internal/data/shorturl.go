package data

import (
	"backend-shorturl/internal/biz"
	"backend-shorturl/internal/data/ent"
	"backend-shorturl/internal/data/ent/shorturl"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

const CachePrefix = "shortUrl:"

type ShortUrlRepo struct {
	data *Data
	log  *log.Helper
}

func NewShortUrlRepo(data *Data, logger log.Logger) biz.ShortUrlRepo {
	return &ShortUrlRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ShortUrlRepo) Save(ctx context.Context, su *biz.ShortUrl) (*biz.ShortUrl, error) {
	si, err := r.data.client.Shorturl.
		Create().
		SetShortCode(su.ShortCode).
		SetLongURL(su.LongUrl).
		SetExpireAt(su.ExpireTime).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, err
	}
	if err != nil {
		log.Errorf("插入短链接失败: %v", err) // 打印错误
		return nil, nil
	}
	return &biz.ShortUrl{
		ShortCode:   si.ShortCode,
		LongUrl:     si.LongURL,
		AccessCount: int32(si.AccessCount),
		ClientIp:    si.CreatorIP,
	}, nil
}

func (r *ShortUrlRepo) IsExit(ctx context.Context, url *biz.ShortUrl) (bool, error) {
	Exit, err := r.data.client.Shorturl.Query().Where(shorturl.LongURLEQ(url.LongUrl)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return Exit, nil
}

func (r *ShortUrlRepo) IsExitByCode(ctx context.Context, code string) (bool, error) {
	Exit, err := r.data.client.Shorturl.Query().Where(shorturl.ShortCodeEQ(code)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return Exit, nil
}

func (r *ShortUrlRepo) GetOriginalURL(ctx context.Context, code string) (*biz.ShortUrl, error) {
	// 使用Ent客户端查询
	url, err := r.data.client.Shorturl.
		Query().                           //开始一个查询构建
		Where(shorturl.ShortCodeEQ(code)). //添加过滤条件：短码等于传入的code
		First(ctx)                         //执行查询并返回第一条结果

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("short code not found")
		}
		return nil, err
	}

	return &biz.ShortUrl{
		LongUrl:   url.LongURL,
		ShortCode: url.ShortCode,
	}, nil
}

func (r *ShortUrlRepo) GetShortStaticsInfo(ctx context.Context, code string) (*biz.ShortStaticsInfo, error) {
	cc, err := r.data.client.Shorturl.Query().Where(shorturl.ShortCodeEQ(code)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.ShortStaticsInfo{
		ShortCode:   cc.ShortCode,
		OriginUrl:   cc.LongURL,
		TotalClicks: uint32(cc.AccessCount),
	}, nil
}

func (r *ShortUrlRepo) GetAllShorStaticsInfos(ctx context.Context, page, size int) ([]*biz.ShortStaticsInfo, int, int, error) {
	//ent orm的分页使用Order、limit、Offset进行分页
	// 计算偏移量
	offset := (page - 1) * size

	list, err := r.data.client.Shorturl.Query().
		Order().
		Limit(size).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	result := make([]*biz.ShortStaticsInfo, 0)
	for _, v := range list {
		result = append(result, &biz.ShortStaticsInfo{
			ShortCode:   v.ShortCode,
			OriginUrl:   v.LongURL,
			TotalClicks: uint32(v.AccessCount),
		})
	}
	//获取总记录数
	total, _ := r.GetTotalCount(ctx)
	//计算总页数
	total_page := r.CalculateTotalPages(total, size)

	return result, total_page, total, nil
}

// GetTotalCount 获取总记录数
func (r *ShortUrlRepo) GetTotalCount(ctx context.Context) (int, error) {
	return r.data.client.Shorturl.
		Query().
		Count(ctx)
}

// CalculateTotalPages 计算总页数
func (r *ShortUrlRepo) CalculateTotalPages(totalCount, pageSize int) int {
	if totalCount == 0 {
		return 1
	}
	return (totalCount + pageSize - 1) / pageSize // 向上取整
}

//redis相关

// SetCache 设置缓存短链信息
func (r *ShortUrlRepo) SetCache(ctx context.Context, key string, shortUrl *biz.ShortUrl) error {
	//序列化对象
	data, err := json.Marshal(shortUrl)
	if err != nil {
		r.log.Errorf("Failed to marshal short url: %v", err)
		return err
	}
	//设置缓存,有效期1小时
	expiration := time.Hour
	return r.data.redis.Set(ctx, key, data, expiration).Err()
}

// GetCache 获取短链信息
func (r *ShortUrlRepo) GetCache(ctx context.Context, key string) (*biz.ShortUrl, error) {
	value, err := r.data.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err //缓存未命中
		}
		r.log.Errorf("Failed to get short URL from cache: %v", err)
		return nil, err
	}
	var shortUrl biz.ShortUrl
	if err := json.Unmarshal([]byte(value), &shortUrl); err != nil {
		r.log.Errorf("Failed to unmarshal short URL: %v", err)
		return nil, err
	}
	return &shortUrl, nil
}

// AddCache 用于缓存添加数据
func (r *ShortUrlRepo) AddCache(ctx context.Context, key string) error {
	shortUrlInfo, err := r.GetCache(ctx, key)
	if err != nil {
		return err
	}
	shortUrlInfo.AccessCount++
	data, err := json.Marshal(shortUrlInfo)
	expiration := time.Hour
	return r.data.redis.Set(ctx, key, data, expiration).Err()
}

// SetEmptyCache 设置空值缓存（防止缓存穿透）
func (r *ShortUrlRepo) SetEmptyCache(ctx context.Context, key string) error {
	return r.data.redis.Set(ctx, key, "null", 10*time.Minute).Err()
}

// IncrementAccessCount 增加访问计数
func (r *ShortUrlRepo) IncrementAccessCount(ctx context.Context, shortCode string) error {
	// 1. 更新 Redis 计数
	key := CachePrefix + shortCode
	fmt.Println("key:", key)
	shortUrlInfo, err := r.GetCache(ctx, key)
	if err != nil {
		return err
	}
	shortUrlInfo.AccessCount++
	//序列化
	data, _ := json.Marshal(shortUrlInfo)
	expiration := time.Hour
	return r.data.redis.Set(ctx, key, data, expiration).Err()

	// 2. 异步持久化到数据库（使用 NSQ 或其他消息队列）
	// ... Tomorrow TODO

}
