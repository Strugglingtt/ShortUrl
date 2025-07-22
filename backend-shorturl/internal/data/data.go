package data

import (
	"backend-shorturl/internal/conf"
	"backend-shorturl/internal/data/ent"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	NewRedisClient,
	NewEntClient,
	NewShortUrlRepo,
	//NewMongoDbClient,    不知道为什么连接不上，暂时先不用mongodb写一下
)

// Data .
type Data struct {
	//这里如果使用ent导入，一定要先生成Client，才不会报错,否则会引入官方自带的ent包，但是官方自带的ent包具体是干嘛的需要研究
	// TODO wrapped database client
	client *ent.Client
	redis  *redis.Client
	//mongo  *mongo.Client
}

// NewData .
func NewData(logger log.Logger, entClient *ent.Client, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the client resources")
	}
	return &Data{
		client: entClient,
		redis:  redis,
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
	//// 自动创建表结构（开发环境推荐）
	//if err := client.Schema.Create(context.Background()); err != nil {
	//	log.Fatalf("failed creating schema resources: %v", err)
	//}
	return client, nil
}

func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	helper := log.NewHelper(log.With(logger, "module", "public/redis"))
	Rdbclient := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Redis.ReadTimeout.AsDuration()))
	defer cancel()
	err := Rdbclient.Ping(timeout).Err()
	if err != nil {
		helper.Fatalf("failed connecting to redis: %v", err)
	}
	log.NewHelper(logger).Info("connected to redis")
	return Rdbclient, nil
}

func NewMongoDbClient(c *conf.Data, logger log.Logger) (*mongo.Client, error) {
	// 从配置中读取 URI，而非硬编码
	uri := c.Mongodb.Uri
	if uri == "" {
		uri = "mongodb://localhost:27017" // 默认值
	}

	// 正确设置超时时间（直接使用配置的秒数）
	timeout := c.Mongodb.Timeout.AsDuration() * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建 MongoDB 客户端
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(timeout).        // 设置连接超时
		SetServerSelectionTimeout(timeout) // 设置服务器选择超时

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 验证连接
	if err := mongoClient.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		// 关闭客户端连接
		_ = mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.NewHelper(logger).Info("connected to mongodb")
	return mongoClient, nil
}
