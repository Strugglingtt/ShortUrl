package biz

import (
	"backend-shorturl/api/shorturl/public/v1"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"

	pb "backend-shorturl/api/shorturl/public/v1"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "url not found")
)

// ShortUrl is a ShortUrl model.
type ShortUrl struct {
	longUrl    string
	expireTime string
}

// ShortUrlRepo is a Greater repo.
type ShortUrlRepo interface {
	Save(ctx context.Context, url *ShortUrl) (string, error)      // 保存URL并返回短码
	FindByKey(ctx context.Context, key string) (*ShortUrl, error) // 通过短码查找URL
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
func (uc *ShortUrlUsecase) CreateShortUrl(ctx context.Context, originalUrl, customKey string) (*pb.ShortenReply, error) {
	// 1. 业务逻辑：验证URL格式
	if !isValidURL(originalUrl) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid url format")
	}

	// 2. 业务逻辑：生成短码（如果没有自定义短码）
	key := customKey
	if key == "" {
		key = generateRandomKey(6) // 生成6位随机短码
	}

	// 3. 业务逻辑：检查短码是否已存在（唯一性校验）
	exists, err := uc.repo.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, status.Errorf(codes.AlreadyExists, "custom key already exists")
	}

	// 4. 创建领域对象
	url := &ShortUrl{}

	// 5. 调用仓储保存
	return uc.repo.Save(ctx, url)
}

// 辅助方法：验证URL格式
func isValidURL(url string) bool {
	_, err := url.ParseRequestURI(url)
	return err == nil
}

// 辅助方法：生成随机短码
func generateRandomKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
