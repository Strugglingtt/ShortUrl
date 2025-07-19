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

func (r *ShortUrlRepo) CreateShortUrl(ctx context.Context, su *biz.ShortUrl) (*ent.Shorturl, error) {
	return r.data.data.Shorturl.
		Create().
		SetShortCode(su.ShortCode).
		SetLongURL(su.LongUrl).
		SetExpireAt(su.ExpireTime).
		Save(ctx)
}
