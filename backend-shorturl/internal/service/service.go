package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"backend-shorturl/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewPublicService)

// NewPublicService 创建PublicService实例
func NewPublicService(uc *biz.ShortUrlUsecase, logger log.Logger) *PublicService {
	return &PublicService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}
