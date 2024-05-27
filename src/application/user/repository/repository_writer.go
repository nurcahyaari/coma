package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryUserWrite struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserWriter(name string, db *database.Clover) repository.RepositoryUserWriter {
	db.DB.CreateCollection(name)
	return &RepositoryUserWrite{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserWrite) SaveUser(ctx context.Context, user entity.User) error {
	dataMap, err := user.MapStringInterface()
	if err != nil {
		return err
	}

	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	_, err = r.db.DB.InsertOne(r.name, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUserWrite) DeleteUser(ctx context.Context, filter entity.FilterUser) error {
	err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		Delete()
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUserWrite) UpdateUser(ctx context.Context, user entity.User) error {
	dataMap, err := user.MapStringInterface()
	if err != nil {
		return err
	}

	err = r.db.DB.Query(r.name).
		UpdateById(user.Id, dataMap)
	if err != nil {
		return err
	}

	return nil
}
