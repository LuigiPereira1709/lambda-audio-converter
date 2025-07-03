package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client      *mongo.Client
	db          *mongo.Database
	connectOnce sync.Once
	connectErr  error
)

const (
	disconnectionTimeout = 10 * time.Second // Timeout for closing the connection
	pingTimeout          = 5 * time.Second  // Timeout for the ping command
)

func GetDatabase() (*mongo.Database, error) {
	newClient()

	if connectErr != nil {
		return nil, connectErr
	}

	return db, connectErr
}

func CloseConnection() error {
	if client == nil {
		return fmt.Errorf("no MongoDB client to close: %w", errors.New("client is nil"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), disconnectionTimeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to close MongoDB connection: %w", err)
	}

	slog.Info("MongoDB connection closed successfully")
	return nil
}

func newClient() {
	connectOnce.Do(func() {
		uri := os.Getenv("MONGO_URI")
		dbName := os.Getenv("MONGO_DB")

		if uri == "" || dbName == "" {
			connectErr = fmt.Errorf("MONGO_URI and MONGO_DB environment variables must be set")
			return
		}

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		c, err := mongo.Connect(opts)
		if err != nil {
			connectErr = fmt.Errorf("failed to connect to MongoDB: %w", err)
			return
		}

		if err := testConnection(c); err != nil {
			connectErr = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
		}

		client = c
		db = client.Database(dbName)
		slog.Info("Connected to MongoDB successfully", "dbName", dbName)
	})
}

func testConnection(c *mongo.Client) error {
	pingCtx, pingCancel := context.WithTimeout(context.Background(), pingTimeout)
	defer pingCancel()

	if err := c.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return nil
}
