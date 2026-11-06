package database

import (
	"encoding/json"
	"fmt"
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/domain"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
)

type tag struct {
	Id   domain.Id `json:"id"`
	Name string    `json:"name"`
}

const tagPrefix = "tag"

func tagKey(tagId domain.Id) []byte {
	return fmt.Appendf(nil, "%s:%s", tagPrefix, tagId)
}

type TagRepository struct {
	db *badger.DB
}

func NewTagRepository(db *badger.DB) TagRepository {
	return TagRepository{db}
}

func (r *TagRepository) Create(tag domain.Tag) (domain.Id, error) {
	id := domain.NewId()
	key := tagKey(id)
	tag.Id = id
	b, err := json.Marshal(toRepoTag(tag))
	if err != nil {
		return domain.Id(uuid.Nil()), err
	}
	if err := r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return domain.Id(uuid.Nil()), err
	}

	return id, nil
}

func (r *TagRepository) Update(tag domain.Tag) error {
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

func (r *TagRepository) Get(id domain.Id) (*domain.Tag, error) {
	key := tagKey(id)
	out, err := GetByKey[tag](r.db, key)
	if err != nil {
		return nil, err
	}
	modelTag := toModelTag(out)
	return &modelTag, nil
}

func (r *TagRepository) List() ([]domain.Tag, error) {
	tags, err := ListAllByKey[tag](r.db, []byte(tagPrefix))
	if err != nil {
		return []domain.Tag{}, err
	}
	var modelTags []domain.Tag
	for _, tag := range tags {
		modelTags = append(modelTags, toModelTag(tag))
	}

	return modelTags, nil
}

func (r *TagRepository) Delete(tagId domain.Id) error {
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

func toRepoTag(modelTag domain.Tag) tag {
	return tag{
		Id:   modelTag.Id,
		Name: modelTag.Name,
	}
}

func toModelTag(tag tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
