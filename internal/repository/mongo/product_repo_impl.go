package mongo

import (
	"context"
	"strings"

	"clofi/internal/model"
	"clofi/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *MongoProductRepository {
	coll := db.Collection("products")
	ensureIndexes(context.Background(), coll)
	return &MongoProductRepository{collection: coll}
}

func ensureIndexes(ctx context.Context, coll *mongo.Collection) {

	textModel := mongo.IndexModel{
		Keys: bson.D{{Key: "$**", Value: "text"}},
	}
	coll.Indexes().CreateOne(ctx, textModel)

	priceModel := mongo.IndexModel{
		Keys: bson.D{{Key: "price", Value: 1}},
	}
	coll.Indexes().CreateOne(ctx, priceModel)

	coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"brand", 1}}})
	coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"availability", 1}}})
}

func (r *MongoProductRepository) Create(ctx context.Context, product *model.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *MongoProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	var product model.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *MongoProductRepository) FindAll(ctx context.Context, filters repository.ProductFilters, page, limit int) ([]*model.Product, error) {
	filter := bson.M{}

	// Фильтрация по бренду
	if filters.Brand != nil {
		filter["brand"] = *filters.Brand
	}

	// Фильтрация по наличию
	if filters.Availability != nil {
		filter["availability"] = *filters.Availability
	}

	// Настройка сортировки
	var sort bson.D
	switch strings.ToLower(filters.SortByPrice) {
	case "asc":
		sort = bson.D{{"price", 1}}
	case "desc":
		sort = bson.D{{"price", -1}}
	default:
		sort = bson.D{{"name", 1}}
	}

	opts := options.Find().SetSort(sort).SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *MongoProductRepository) SearchByText(ctx context.Context, query string, page, limit int) ([]*model.Product, error) {
	if query == "" {
		return r.FindAll(ctx, repository.ProductFilters{}, page, limit)
	}

	filter := bson.M{
		"$text": bson.M{"$search": query},
	}

	opts := options.Find().
		SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	for cursor.Next(ctx) {
		var p model.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, cursor.Err()
}
