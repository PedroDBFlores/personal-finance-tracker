package database_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestTagRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("it creates a new tag", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewTagRepository(db)
		tag := domain.Tag{
			Name: "Test Tag",
		}
		id, err := repo.Create(tag)
		require.NoError(t, err)

		found, err := repo.Get(id)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, id, found.Id)
		require.Equal(t, tag.Name, found.Name)
	})

	t.Run("get", func(t *testing.T) {
		t.Run("it gets a tag by id", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)
			tag := domain.Tag{
				Name: "Test Tag",
			}
			id, err := repo.Create(tag)
			require.NoError(t, err)

			found, err := repo.Get(id)
			require.NoError(t, err)
			require.Equal(t, id, found.Id)
			require.Equal(t, tag.Name, found.Name)
		})

		t.Run("returns ResourceNotFoundError if the tag doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)

			id := domain.NewId()
			acc, err := repo.Get(id)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Tag", fmt.Sprintf("tag:%s", id)))
			require.Nil(t, acc)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("it updates a tag", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)
			tag := domain.Tag{
				Name: "Test Tag",
			}
			id, err := repo.Create(tag)
			require.NoError(t, err)

			tag.Id = id
			tag.Name = "Updated Tag"
			err = repo.Update(tag)
			require.NoError(t, err)

			found, err := repo.Get(id)
			require.NoError(t, err)
			require.Equal(t, id, found.Id)
			require.Equal(t, "Updated Tag", found.Name)
		})

		t.Run("returns ResourceNotFoundError if the tag doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)

			tag := domain.Tag{
				Id:   domain.NewId(),
				Name: "Non-existent Tag",
			}
			err = repo.Update(tag)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Tag", fmt.Sprintf("tag:%s", tag.Id)))
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("it lists all tags", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)

			tag1 := domain.Tag{Name: "Tag 1"}
			tag2 := domain.Tag{Name: "Tag 2"}
			_, err = repo.Create(tag1)
			require.NoError(t, err)
			_, err = repo.Create(tag2)
			require.NoError(t, err)

			tags, err := repo.List()
			require.NoError(t, err)
			require.Len(t, tags, 2)
			require.Contains(t, []string{tags[0].Name, tags[1].Name}, "Tag 1")
			require.Contains(t, []string{tags[0].Name, tags[1].Name}, "Tag 2")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("deletes a tag by id", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)
			tag := domain.Tag{
				Name: "Random Tag",
			}
			id, err := repo.Create(tag)
			require.NoError(t, err)

			err = repo.Delete(id)
			require.NoError(t, err)

			found, err := repo.Get(id)
			require.Error(t, err)
			require.Nil(t, found)
		})

		t.Run("returns ResourceNotFoundError if the tag doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewTagRepository(db)

			id := domain.NewId()
			err = repo.Delete(id)

			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Tag", fmt.Sprintf("tag:%s", id)))
		})
	})
}
