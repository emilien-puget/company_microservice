package company

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName         = "your_database_name"
	collectionName = "companies"
)

type MongoModel struct {
	ID          string      `bson:"id"`
	Name        string      `bson:"name"`
	Description string      `bson:"description,omitempty"`
	Employees   int         `bson:"employees"`
	Registered  bool        `bson:"registered"`
	Type        CompanyType `bson:"type"`
}

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewCompanyRepository(client *mongo.Client) *Repository {
	collection := client.Database(dbName).Collection(collectionName)

	return &Repository{
		client:     client,
		collection: collection,
	}
}

var ErrCompanyNotFound = errors.New("company not found")

func (r *Repository) Create(ctx context.Context, company *MongoModel) error {
	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FetchByID(ctx context.Context, id uuid.UUID) (*MongoModel, error) {
	var company MongoModel

	filter := bson.M{"id": id.String()}

	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}

	return &company, nil
}

func (r *Repository) FetchByName(ctx context.Context, name string) (*MongoModel, error) {
	var company MongoModel

	filter := bson.M{"name": name}

	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}

	return &company, nil
}

func (r *Repository) Patch(ctx context.Context, id uuid.UUID, updatedCompany *MongoModel) error {
	filter := bson.M{"id": id.String()}
	update := bson.M{
		"$set": bson.M{
			"name":        updatedCompany.Name,
			"description": updatedCompany.Description,
			"employees":   updatedCompany.Employees,
			"registered":  updatedCompany.Registered,
			"type":        updatedCompany.Type,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id.String()})
	if err != nil {
		return err
	}
	return nil
}
