use crate::models::Transaction;
use crate::storage::repository::{RepositoryError, TransactionRepository};
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use sqlx::{Pool, Sqlite, Row};
use std::str::FromStr;
use uuid::Uuid;

/// Creates a transaction from database row data
fn create_transaction_from_row(
    id: Uuid,
    amount: String,
    transaction_type: String,
    category_id: Option<Uuid>,
    description: String,
    date: DateTime<Utc>,
    created_at: DateTime<Utc>,
    is_deleted: bool,
) -> Result<Transaction, RepositoryError> {
    let amount = Decimal::from_str(&amount)
        .map_err(|e| RepositoryError::Database(sqlx::Error::Protocol(e.to_string())))?;
    
    let transaction_type = transaction_type.parse()
        .map_err(|e| RepositoryError::Database(sqlx::Error::Protocol(e)))?;

    Ok(Transaction {
        id,
        amount,
        transaction_type,
        category_id,
        description,
        date,
        created_at,
        is_deleted,
    })
}

/// Converts a list of rows to transactions
fn rows_to_transactions(rows: Vec<sqlx::sqlite::SqliteRow>) -> Result<Vec<Transaction>, RepositoryError> {
    let mut transactions = Vec::new();
    for row in rows {
        let transaction = create_transaction_from_row(
            row.get("id"),
            row.get("amount"),
            row.get("transaction_type"),
            row.get("category_id"),
            row.get("description"),
            row.get("date"),
            row.get("created_at"),
            row.get("is_deleted"),
        )?;
        transactions.push(transaction);
    }
    Ok(transactions)
}

/// SQLite implementation of TransactionRepository
#[derive(Debug, Clone)]
pub struct SqliteTransactionRepository {
    pool: Pool<Sqlite>,
}

impl SqliteTransactionRepository {
    pub fn new(pool: Pool<Sqlite>) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl TransactionRepository for SqliteTransactionRepository {
    async fn save(&self, transaction: &Transaction) -> Result<(), RepositoryError> {
        let query = sqlx::query(
            r#"
            INSERT INTO transactions (id, amount, transaction_type, category_id, description, date, created_at, is_deleted)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            "#
        )
        .bind(transaction.id)
        .bind(transaction.amount_to_string())
        .bind(transaction.transaction_type.to_string())
        .bind(transaction.category_id)
        .bind(&transaction.description)
        .bind(transaction.date)
        .bind(transaction.created_at)
        .bind(transaction.is_deleted);

        query.execute(&self.pool).await?;
        Ok(())
    }

    async fn find_by_id(&self, id: Uuid) -> Result<Option<Transaction>, RepositoryError> {
        let row = sqlx::query(
            r#"
            SELECT id, amount, transaction_type, category_id, description, date, created_at, is_deleted
            FROM transactions
            WHERE id = $1 AND is_deleted = false
            "#
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;

        match row {
            Some(row) => {
                let transaction = create_transaction_from_row(
                    row.get("id"),
                    row.get("amount"),
                    row.get("transaction_type"),
                    row.get("category_id"),
                    row.get("description"),
                    row.get("date"),
                    row.get("created_at"),
                    row.get("is_deleted"),
                )?;
                Ok(Some(transaction))
            }
            None => Ok(None),
        }
    }

    async fn find_all(&self) -> Result<Vec<Transaction>, RepositoryError> {
        let rows = sqlx::query(
            r#"
            SELECT id, amount, transaction_type, category_id, description, date, created_at, is_deleted
            FROM transactions
            WHERE is_deleted = false
            ORDER BY date DESC, created_at DESC
            "#
        )
        .fetch_all(&self.pool)
        .await?;

        rows_to_transactions(rows)
    }

    async fn find_by_category(&self, category_id: Uuid) -> Result<Vec<Transaction>, RepositoryError> {
        let rows = sqlx::query(
            r#"
            SELECT id, amount, transaction_type, category_id, description, date, created_at, is_deleted
            FROM transactions
            WHERE category_id = $1 AND is_deleted = false
            ORDER BY date DESC, created_at DESC
            "#
        )
        .bind(category_id)
        .fetch_all(&self.pool)
        .await?;

        rows_to_transactions(rows)
    }

    async fn find_by_type(&self, transaction_type: &crate::models::transaction::TransactionType) -> Result<Vec<Transaction>, RepositoryError> {
        let type_str = transaction_type.to_string();
        let rows = sqlx::query(
            r#"
            SELECT id, amount, transaction_type, category_id, description, date, created_at, is_deleted
            FROM transactions
            WHERE transaction_type = $1 AND is_deleted = false
            ORDER BY date DESC, created_at DESC
            "#
        )
        .bind(type_str)
        .fetch_all(&self.pool)
        .await?;

        rows_to_transactions(rows)
    }

    async fn find_by_date_range(
        &self,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Result<Vec<Transaction>, RepositoryError> {
        let rows = sqlx::query(
            r#"
            SELECT id, amount, transaction_type, category_id, description, date, created_at, is_deleted
            FROM transactions
            WHERE date >= $1 AND date <= $2 AND is_deleted = false
            ORDER BY date DESC, created_at DESC
            "#
        )
        .bind(start)
        .bind(end)
        .fetch_all(&self.pool)
        .await?;

        rows_to_transactions(rows)
    }

    async fn update(&self, transaction: &Transaction) -> Result<(), RepositoryError> {
        let query = sqlx::query(
            r#"
            UPDATE transactions
            SET amount = $1, transaction_type = $2, category_id = $3, description = $4, date = $5, is_deleted = $6
            WHERE id = $7
            "#
        )
        .bind(transaction.amount_to_string())
        .bind(transaction.transaction_type.to_string())
        .bind(transaction.category_id)
        .bind(&transaction.description)
        .bind(transaction.date)
        .bind(transaction.is_deleted)
        .bind(transaction.id);

        let rows_affected = query.execute(&self.pool).await?.rows_affected();
        
        if rows_affected == 0 {
            return Err(RepositoryError::NotFound);
        }
        
        Ok(())
    }

    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError> {
        let query = sqlx::query(
            r#"UPDATE transactions SET is_deleted = true WHERE id = $1"#
        )
        .bind(id);

        let rows_affected = query.execute(&self.pool).await?.rows_affected();
        
        if rows_affected == 0 {
            return Err(RepositoryError::NotFound);
        }
        
        Ok(())
    }

    async fn count(&self) -> Result<i64, RepositoryError> {
        let query = sqlx::query_scalar::<_, Option<i64>>(
            r#"SELECT COUNT(*) FROM transactions WHERE is_deleted = false"#
        );

        let count = query.fetch_one(&self.pool).await?;
        Ok(count.unwrap_or(0))
    }

    async fn net_balance(&self) -> Result<Decimal, RepositoryError> {
        // Get all transactions and calculate net balance in memory
        let transactions = self.find_all().await?;
        let net = transactions.iter()
            .map(|t| t.effective_amount())
            .fold(Decimal::ZERO, |acc, amount| acc + amount);
        Ok(net)
    }

    async fn total_credits(&self) -> Result<Decimal, RepositoryError> {
        let transactions = self.find_by_type(&crate::models::transaction::TransactionType::Credit).await?;
        let total = transactions.iter()
            .map(|t| t.amount)
            .fold(Decimal::ZERO, |acc, amount| acc + amount);
        Ok(total)
    }

    async fn total_debits(&self) -> Result<Decimal, RepositoryError> {
        let transactions = self.find_by_type(&crate::models::transaction::TransactionType::Debit).await?;
        let total = transactions.iter()
            .map(|t| t.amount)
            .fold(Decimal::ZERO, |acc, amount| acc + amount);
        Ok(total)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::storage::Database;
    use rust_decimal_macros::dec;

    async fn get_test_db() -> (Database, SqliteTransactionRepository) {
        let db = Database::in_memory().await.expect("Failed to create test database");
        let repo = SqliteTransactionRepository::new(db.pool.clone());
        (db, repo)
    }

    #[tokio::test]
    async fn should_save_and_find_transaction() {
        let (_db, repo) = get_test_db().await;

        let transaction = Transaction::new(
            dec!(100.00),
            crate::models::transaction::TransactionType::Credit,
            None,
            "Test transaction".to_string(),
            Utc::now(),
        );

        repo.save(&transaction).await.expect("Failed to save");
        let found = repo.find_by_id(transaction.id).await.expect("Failed to find");

        assert!(found.is_some());
        assert_eq!(found.unwrap().amount, dec!(100.00));
    }

    #[tokio::test]
    async fn should_find_all_transactions() {
        let (_db, repo) = get_test_db().await;

        let t1 = Transaction::new(dec!(100.00), crate::models::transaction::TransactionType::Credit, None, "T1".to_string(), Utc::now());
        let t2 = Transaction::new(dec!(50.00), crate::models::transaction::TransactionType::Debit, None, "T2".to_string(), Utc::now());

        repo.save(&t1).await.unwrap();
        repo.save(&t2).await.unwrap();

        let all = repo.find_all().await.expect("Failed to find all");
        assert_eq!(all.len(), 2);
    }

    #[tokio::test]
    async fn should_calculate_net_balance() {
        let (_db, repo) = get_test_db().await;

        let credit = Transaction::new(dec!(200.00), crate::models::transaction::TransactionType::Credit, None, "".to_string(), Utc::now());
        let debit = Transaction::new(dec!(100.00), crate::models::transaction::TransactionType::Debit, None, "".to_string(), Utc::now());

        repo.save(&credit).await.unwrap();
        repo.save(&debit).await.unwrap();

        let net = repo.net_balance().await.expect("Failed to get net balance");
        assert_eq!(net, dec!(100.00));
    }

    #[tokio::test]
    async fn should_count_transactions() {
        let (_db, repo) = get_test_db().await;

        let t1 = Transaction::new(dec!(10.00), crate::models::transaction::TransactionType::Credit, None, "".to_string(), Utc::now());
        let t2 = Transaction::new(dec!(20.00), crate::models::transaction::TransactionType::Debit, None, "".to_string(), Utc::now());

        repo.save(&t1).await.unwrap();
        repo.save(&t2).await.unwrap();

        let count = repo.count().await.expect("Failed to count");
        assert_eq!(count, 2);
    }

    #[tokio::test]
    async fn should_delete_transaction() {
        let (_db, repo) = get_test_db().await;

        let transaction = Transaction::new(dec!(100.00), crate::models::transaction::TransactionType::Credit, None, "".to_string(), Utc::now());
        repo.save(&transaction).await.unwrap();

        repo.delete(transaction.id).await.expect("Failed to delete");
        let found = repo.find_by_id(transaction.id).await.expect("Failed to find");

        assert!(found.is_none());
    }

    #[tokio::test]
    async fn should_find_by_type() {
        let (_db, repo) = get_test_db().await;

        let credit = Transaction::new(dec!(100.00), crate::models::transaction::TransactionType::Credit, None, "".to_string(), Utc::now());
        let debit = Transaction::new(dec!(50.00), crate::models::transaction::TransactionType::Debit, None, "".to_string(), Utc::now());

        repo.save(&credit).await.unwrap();
        repo.save(&debit).await.unwrap();

        let credits = repo.find_by_type(&crate::models::transaction::TransactionType::Credit).await.unwrap();
        let debits = repo.find_by_type(&crate::models::transaction::TransactionType::Debit).await.unwrap();

        assert_eq!(credits.len(), 1);
        assert_eq!(debits.len(), 1);
    }
}
