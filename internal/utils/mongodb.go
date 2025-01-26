package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/LtePrince/Personal-Website-backend/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientLock sync.Mutex
	timer      *time.Timer
)

const (
	mongoURI       = "mongodb://localhost:27017" // 替换为你的 MongoDB URI
	databaseName   = "WebsiteBlog"               // 替换为你的数据库名
	collectionName = "blogs"                     // 替换为你的集合名
	idleTimeout    = 1 * time.Hour
)

// ConnectMongoDB 连接 MongoDB 并将客户端保存为全局变量
func ConnectMongoDB() error {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		resetIdleTimer()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	resetIdleTimer()
	return nil
}

// resetIdleTimer 重置空闲定时器
func resetIdleTimer() {
	if timer != nil {
		timer.Stop()
	}
	timer = time.AfterFunc(idleTimeout, func() {
		clientLock.Lock()
		defer clientLock.Unlock()
		if client != nil {
			_ = client.Disconnect(context.Background())
			client = nil
			log.Println("MongoDB connection closed due to inactivity.")
		}
	})
}

// GetBlogTitlesAndSummaries 从 MongoDB 获取所有 Blog 的标题和概述
func GetBlogTitlesAndSummaries() ([]api.BlogResponse, error) {
	if client == nil {
		if err := ConnectMongoDB(); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(databaseName).Collection(collectionName)

	// 定义查询条件（如果没有条件，可以使用 bson.D{}）
	filter := bson.D{}

	// 定义只返回 `title` 和 `summary` 的字段
	projection := bson.D{
		{"Title", 1},
		{"Summary", 1},
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %v", err)
	}
	defer cursor.Close(ctx)

	var blogs []api.BlogResponse
	for cursor.Next(ctx) {
		var blog api.BlogResponse
		if err := cursor.Decode(&blog); err != nil {
			return nil, fmt.Errorf("failed to decode blog: %v", err)
		}
		blogs = append(blogs, blog)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return blogs, nil
}
