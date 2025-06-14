package repository

import (
	"context"
	"tunnerse/database"
	"tunnerse/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TunnelRepository struct {
	Client *mongo.Client
	DBName string
}

func NewTunnelRepository() *TunnelRepository {
	return &TunnelRepository{
		Client: database.MongoClient,
		DBName: database.MongoDB.Name(),
	}
}

func (r *TunnelRepository) Register(tunnel *models.Tunnel) error {
	collection := r.Client.Database(r.DBName).Collection("tunnels")
	_, err := collection.InsertOne(context.TODO(), tunnel)
	if err != nil {
		return err
	}

	return nil
}

func (r *TunnelRepository) GetTunnelByName(name string) (*models.Tunnel, error) {
	collection := r.Client.Database(r.DBName).Collection("tunnels")

	filter := bson.D{{Key: "name", Value: name}}

	var tunnel models.Tunnel
	err := collection.FindOne(context.TODO(), filter).Decode(&tunnel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &tunnel, nil
}
