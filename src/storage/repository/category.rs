use crate::models::Category;
use crate::storage::repository::{RepositoryError, CategoryRepository};
use async_trait::async_trait;
use sqlx::{Pool, Sqlite, Row};
use uuid::Uuid;

/// SQLite implementation of CategoryRepository
#[derive(Debug, Clone)]
pub struct SqliteCategoryRepository {
    pool: Pool<Sqlite>,
}

impl SqliteCategoryRepository {
    pub fn new(pool: Pool<Sqlite>) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl CategoryRepository for SqliteCategoryRepository {
    async fn save(&self, category: &Category) -> Result<(), RepositoryError> {
        let query = sqlx::query(
            r#"
            INSERT INTO categories (id, name, color, is_default)
            VALUES ($1, $2, $3, $4)
            "#
        )
        .bind(category.id)
        .bind(&category.name)
        .bind(&category.color)
        .bind(category.is_default);

        query.execute(&self.pool).await?;
        Ok(())
    }

    async fn find_by_id(&self, id: Uuid) -> Result<Option<Category>, RepositoryError> {
        let row = sqlx::query(
            r#"
            SELECT id, name, color, is_default
            FROM categories
            WHERE id = $1
            "#
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;

        match row {
            Some(row) => {
                let category = Category {
                    id: row.get("id"),
                    name: row.get("name"),
                    color: row.get("color"),
                    is_default: row.get("is_default"),
                };
                Ok(Some(category))
            }
            None => Ok(None),
        }
    }

    async fn find_all(&self) -> Result<Vec<Category>, RepositoryError> {
        let rows = sqlx::query(
            r#"
            SELECT id, name, color, is_default
            FROM categories
            ORDER BY is_default DESC, name ASC
            "#
        )
        .fetch_all(&self.pool)
        .await?;

        let mut categories = Vec::new();
        for row in rows {
            let category = Category {
                id: row.get("id"),
                name: row.get("name"),
                color: row.get("color"),
                is_default: row.get("is_default"),
            };
            categories.push(category);
        }
        Ok(categories)
    }

    async fn find_default(&self) -> Result<Option<Category>, RepositoryError> {
        let row = sqlx::query(
            r#"
            SELECT id, name, color, is_default
            FROM categories
            WHERE is_default = true
            LIMIT 1
            "#
        )
        .fetch_optional(&self.pool)
        .await?;

        match row {
            Some(row) => {
                let category = Category {
                    id: row.get("id"),
                    name: row.get("name"),
                    color: row.get("color"),
                    is_default: row.get("is_default"),
                };
                Ok(Some(category))
            }
            None => Ok(None),
        }
    }

    async fn update(&self, category: &Category) -> Result<(), RepositoryError> {
        let query = sqlx::query(
            r#"
            UPDATE categories
            SET name = $1, color = $2, is_default = $3
            WHERE id = $4
            "#
        )
        .bind(&category.name)
        .bind(&category.color)
        .bind(category.is_default)
        .bind(category.id);

        let rows_affected = query.execute(&self.pool).await?.rows_affected();
        
        if rows_affected == 0 {
            return Err(RepositoryError::NotFound);
        }
        
        Ok(())
    }

    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError> {
        // First check if any transactions reference this category
        let check_query = sqlx::query_scalar::<_, Option<i64>>(
            r#"SELECT COUNT(*) FROM transactions WHERE category_id = $1"#
        )
        .bind(id);

        let count = check_query.fetch_one(&self.pool).await?;
        
        if count > Some(0) {
            return Err(RepositoryError::NotFound); // Using this as "cannot delete - has references"
        }

        let delete_query = sqlx::query(
            r#"DELETE FROM categories WHERE id = $1"#
        )
        .bind(id);

        let rows_affected = delete_query.execute(&self.pool).await?.rows_affected();
        
        if rows_affected == 0 {
            return Err(RepositoryError::NotFound);
        }
        
        Ok(())
    }

    async fn count(&self) -> Result<i64, RepositoryError> {
        let query = sqlx::query_scalar::<_, Option<i64>>(
            r#"SELECT COUNT(*) FROM categories"#
        );

        let count = query.fetch_one(&self.pool).await?;
        Ok(count.unwrap_or(0))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::storage::Database;

    async fn get_test_db() -> (Database, SqliteCategoryRepository) {
        let db = Database::in_memory().await.expect("Failed to create test database");
        let repo = SqliteCategoryRepository::new(db.pool.clone());
        (db, repo)
    }

    #[tokio::test]
    async fn should_save_and_find_category() {
        let (_db, repo) = get_test_db().await;

        let category = Category::new("Food".to_string(), Some("#ff0000".to_string()));
        repo.save(&category).await.expect("Failed to save");
        
        let found = repo.find_by_id(category.id).await.expect("Failed to find");
        assert!(found.is_some());
        assert_eq!(found.unwrap().name, "Food");
    }

    #[tokio::test]
    async fn should_find_all_categories() {
        let (_db, repo) = get_test_db().await;

        let c1 = Category::new("Food".to_string(), None);
        let c2 = Category::new("Transport".to_string(), None);

        repo.save(&c1).await.unwrap();
        repo.save(&c2).await.unwrap();

        let all = repo.find_all().await.expect("Failed to find all");
        assert!(all.len() >= 2); // At least 2 + default
    }

    #[tokio::test]
    async fn should_find_default_category() {
        let (_db, repo) = get_test_db().await;

        let default = repo.find_default().await.expect("Failed to find default");
        assert!(default.is_some());
        assert_eq!(default.unwrap().name, "Uncategorized");
    }

    #[tokio::test]
    async fn should_count_categories() {
        let (_db, repo) = get_test_db().await;

        let c1 = Category::new("Food".to_string(), None);
        let c2 = Category::new("Transport".to_string(), None);

        repo.save(&c1).await.unwrap();
        repo.save(&c2).await.unwrap();

        let count = repo.count().await.expect("Failed to count");
        assert!(count >= 2);
    }

    #[tokio::test]
    async fn should_delete_category() {
        let (_db, repo) = get_test_db().await;

        let category = Category::new("TestCategory".to_string(), None);
        repo.save(&category).await.unwrap();

        repo.delete(category.id).await.expect("Failed to delete");
        let found = repo.find_by_id(category.id).await.expect("Failed to find");

        assert!(found.is_none());
    }
}
