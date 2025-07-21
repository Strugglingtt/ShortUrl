package data

import (
	"backend-shorturl/internal/biz"
	"backend-shorturl/internal/data/ent"
	"backend-shorturl/internal/data/ent/shorturl"
	"context"
	"errors"
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

func (r *ShortUrlRepo) GetOriginalURL(ctx context.Context, code string) (string, error) {
	// 使用Ent客户端查询
	url, err := r.data.client.Shorturl.
		Query().                           //开始一个查询构建
		Where(shorturl.ShortCodeEQ(code)). //添加过滤条件：短码等于传入的code
		Select(shorturl.FieldLongURL).     //只选择long_url字段返回
		First(ctx)                         //执行查询并返回第一条结果

	if err != nil {
		if ent.IsNotFound(err) {
			return "", errors.New("short code not found")
		}
		return "", err
	}

	return url.LongURL, nil
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
