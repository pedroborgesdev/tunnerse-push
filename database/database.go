package database

import (
	"context"
	"fmt"
	"time"

	"tunnerse/config"
	"tunnerse/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// URI com autenticação
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		config.AppConfig.DBUser, // tunnerse_admin
		config.AppConfig.DBPwd,  // senhaForteAqui123
		config.AppConfig.DBHost, // normalmente "localhost"
		config.AppConfig.DBPort, // normalmente "27017"
		config.AppConfig.DBName) // authSource = tunnerse

	logger.Log("DEBUG", "Attempting MongoDB connection", []logger.LogDetail{
		{Key: "URI", Value: uri},
	})

	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Log("ERROR", "MongoDB connection failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "URI", Value: uri},
		})
		return
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err = client.Ping(pingCtx, nil); err != nil {
		logger.Log("ERROR", "MongoDB ping failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	MongoClient = client
	MongoDB = client.Database(config.AppConfig.DBName)

	logger.Log("INFO", "MongoDB connection established", []logger.LogDetail{
		{Key: "Database", Value: config.AppConfig.DBName},
		{Key: "Host", Value: config.AppConfig.DBHost},
	})
}

func CloseDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := MongoClient.Disconnect(ctx); err != nil {
			logger.Log("ERROR", "Failed to disconnect MongoDB", []logger.LogDetail{
				{Key: "Error", Value: err.Error()},
			})
		} else {
			logger.Log("INFO", "MongoDB connection closed", nil)
		}
	}
}
