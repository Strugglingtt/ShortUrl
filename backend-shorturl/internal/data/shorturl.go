package data

import (
	"backend-shorturl/internal/biz"
	"backend-shorturl/internal/data/ent"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

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

//
//CreateShortUrl(ctx context.Context, url *ShortUrl) (string, error)
//FindByKey(ctx context.Context, key string) (*ShortUrl, error)
//Exists(ctx context.Context, longURL string) (string, bool)

func (r *ShortUrlRepo) Save(ctx context.Context, su *biz.ShortUrl) (int64, error) {
	si, err := r.data.client.Shorturl.
		Create().
		SetShortCode(su.ShortCode).
		SetLongURL(su.LongUrl).
		SetExpireAt(su.ExpireTime).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return 666, err
	}
	if err != nil {
		log.Errorf("插入短链接失败: %v", err) // 打印错误
		return 666, nil
	}
	return si.ID, nil
}
