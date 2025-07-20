package data

import (
	"backend-shorturl/internal/conf"
	"backend-shorturl/internal/data/ent"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewData, NewEntClient, NewShortUrlRepo)

// Data .
type Data struct {
	//这里如果使用ent导入，一定要先生成Client，才不会报错,否则会引入官方自带的ent包，但是官方自带的ent包具体是干嘛的需要研究
	// TODO wrapped database client
	client *ent.Client
}

// NewData .
func NewData(logger log.Logger, entClient *ent.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the client resources")
	}
	return &Data{
		client: entClient,
	}, cleanup, nil
}

func NewEntClient(c *conf.Data, logger log.Logger) (*ent.Client, error) {
	helper := log.NewHelper(log.With(logger, "module", "internal/ent"))
	client, err := ent.Open(
		c.Database.Driver,
		c.Database.Source,
		ent.Log(log.Info), // 启用日志
		ent.Debug(),       // 启用调试模式
	)
	if err != nil {
		helper.Fatalf("failed opening connection to db: %v", err)
	}
	fmt.Println("connect to mysql Success")
	// 自动创建表结构（开发环境推荐）
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client, nil
}
