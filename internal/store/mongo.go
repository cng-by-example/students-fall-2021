package store

import (
	"context"
	"fmt"

	"githuh.com/cng-by-example/students/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const Collection = "students"

type MongoDB struct {
	db *mongo.Database
}

func NewMongoDBStore(db *mongo.Database) MongoDB {
	return MongoDB{
		db: db,
	}
}

func (m MongoDB) Save(ctx context.Context, s model.Student) error {
	_, err := m.db.Collection(Collection).InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("mongodb insert failed %w", err)
	}

	return nil
}

func (m MongoDB) LoadByID(ctx context.Context, id string) (model.Student, error) {
	var student model.Student

	record := m.db.Collection(Collection).FindOne(ctx, bson.M{"id": id})

	if err := record.Decode(&student); err != nil {
		return model.Student{}, fmt.Errorf("reading from mongodb failed %w", err)
	}

	return student, nil
}

func (m MongoDB) Load(ctx context.Context) ([]model.Student, error) {
	var students []model.Student

	records, err := m.db.Collection(Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo find failed %w", err)
	}

	for records.Next(ctx) {
		var student model.Student

		if err := records.Decode(&student); err != nil {
			return students, fmt.Errorf("mongo record decoding failed %w", err)
		}

		students = append(students, student)
	}

	return students, nil
}
