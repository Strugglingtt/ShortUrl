package data

import (
	"context"

	"backend-shorturl/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.ShortUrlRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.ShortUrl) (*biz.ShortUrl, error) {
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.ShortUrl) (*biz.ShortUrl, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*biz.ShortUrl, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.ShortUrl, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*biz.ShortUrl, error) {
	return nil, nil
}
