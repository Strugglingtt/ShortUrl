package data

import (
	"backend-shorturl/internal/conf"
	"backend-shorturl/internal/data/ent"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewEntClient,
	NewRedisClient,
	NewMongoDbClient,
	NewShortUrlRepo,
)

// Data .
type Data struct {
	//这里如果使用ent导入，一定要先生成Client，才不会报错,否则会引入官方自带的ent包，但是官方自带的ent包具体是干嘛的需要研究
	// TODO wrapped database client
	client *ent.Client
	redis  *redis.Client
	mgdb   *mongo.Client
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

func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	helper := log.NewHelper(log.With(logger, "module", "public/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Redis.ReadTimeout.AsDuration()))
	defer cancel()
	err := client.Ping(timeout).Err()
	if err != nil {
		helper.Fatalf("failed connecting to redis: %v", err)
	}
	return client, nil
}

func NewMongoDbClient(c *conf.Data, logger log.Logger) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err.Error())
	}
	defer client.Disconnect(context.Background())
	return client, nil
}
