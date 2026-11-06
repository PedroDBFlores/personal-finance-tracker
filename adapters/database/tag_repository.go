package database

import (
	"encoding/json"
	"fmt"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/models"
)

type Tag struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

const tagPrefix = "tag"

func tagKey(tagId uuid.UUID) []byte {
	return fmt.Appendf(nil, "%s:%s", tagPrefix, tagId)
}

type TagRepository struct {
	db *badger.DB
}

func NewTagRepository(db *badger.DB) TagRepository {
	return TagRepository{db}
}

func (r *TagRepository) Create(tag models.Tag) (uuid.UUID, error) {
	id := uuid.NewV7()
	key := tagKey(id)
	tag.Id = id
	b, err := json.Marshal(toRepoTag(tag))
	if err != nil {
		return uuid.Nil(), err
	}
	if err := r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return uuid.Nil(), err
	}

	return id, nil
}

func (r *TagRepository) Update(tag models.Tag) error {
	key := tagKey(tag.Id)
	ok, err := IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Tag", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err != nil {
			return err
		}
		b, err := json.Marshal(toRepoTag(tag))
		if err != nil {
			return err
		}
		return txn.Set(key, b)
	})
}

func (r *TagRepository) Get(id uuid.UUID) (*models.Tag, error) {
	key := tagKey(id)
	out, err := GetByKey[Tag](r.db, key)
	if err != nil {
		return nil, err
	}
	modelTag := toModelTag(out)
	return &modelTag, nil
}

func (r *TagRepository) List() ([]models.Tag, error) {
	tags, err := ListAllByKey[Tag](r.db, []byte(tagPrefix))
	if err != nil {
		return []models.Tag{}, err
	}
	var modelTags []models.Tag
	for _, tag := range tags {
		modelTags = append(modelTags, toModelTag(tag))
	}

	return modelTags, nil
}

func (r *TagRepository) Delete(tagId uuid.UUID) error {
	key := tagKey(tagId)
	ok, err := IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Tag", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func toRepoTag(tag models.Tag) Tag {
	return Tag{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func toModelTag(tag Tag) models.Tag {
	return models.Tag{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
