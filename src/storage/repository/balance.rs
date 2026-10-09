use crate::models::BalanceSnapshot;
use crate::storage::repository::{BalanceRepository, RepositoryError, TransactionRepository};
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use sqlx::{Pool, Row, Sqlite};

/// SQLite implementation of BalanceRepository
pub struct SqliteBalanceRepository {
    pool: Pool<Sqlite>,
    transaction_repo: Box<dyn TransactionRepository>,
}

impl SqliteBalanceRepository {
    pub fn new(pool: Pool<Sqlite>, transaction_repo: Box<dyn TransactionRepository>) -> Self {
        Self {
            pool,
            transaction_repo,
        }
    }
}

#[async_trait]
impl BalanceRepository for SqliteBalanceRepository {
    async fn save(&self, snapshot: &BalanceSnapshot) -> Result<(), RepositoryError> {
        let (net_str, credit_str, debit_str) = snapshot.to_strings();

        let query = sqlx::query(
            r#"
            INSERT INTO balance_snapshots (id, net_balance, credit_total, debit_total, timestamp)
            VALUES ($1, $2, $3, $4, $5)
            "#,
        )
        .bind(snapshot.id)
        .bind(net_str)
        .bind(credit_str)
        .bind(debit_str)
        .bind(snapshot.timestamp);

        query.execute(&self.pool).await?;
        Ok(())
    }

    async fn find_latest(&self) -> Result<Option<BalanceSnapshot>, RepositoryError> {
        let row = sqlx::query(
            r#"
            SELECT id, net_balance, credit_total, debit_total, timestamp
            FROM balance_snapshots
            ORDER BY timestamp DESC
            LIMIT 1
            "#,
        )
        .fetch_optional(&self.pool)
        .await?;

        match row {
            Some(row) => {
                let snapshot = BalanceSnapshot::from_strings(
                    row.get("id"),
                    row.get("net_balance"),
                    row.get("credit_total"),
                    row.get("debit_total"),
                    row.get("timestamp"),
                )
                .map_err(|e| RepositoryError::Database(sqlx::Error::Protocol(e.to_string())))?;
                Ok(Some(snapshot))
            }
            None => Ok(None),
        }
    }

    async fn find_by_date_range(
        &self,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Result<Vec<BalanceSnapshot>, RepositoryError> {
        let rows = sqlx::query(
            r#"
            SELECT id, net_balance, credit_total, debit_total, timestamp
            FROM balance_snapshots
            WHERE timestamp >= $1 AND timestamp <= $2
            ORDER BY timestamp DESC
            "#,
        )
        .bind(start)
        .bind(end)
        .fetch_all(&self.pool)
        .await?;

        let mut snapshots = Vec::new();
        for row in rows {
            let snapshot = BalanceSnapshot::from_strings(
                row.get("id"),
                row.get("net_balance"),
                row.get("credit_total"),
                row.get("debit_total"),
                row.get("timestamp"),
            )
            .map_err(|e| RepositoryError::Database(sqlx::Error::Protocol(e.to_string())))?;
            snapshots.push(snapshot);
        }
        Ok(snapshots)
    }

    async fn create_snapshot(&self) -> Result<BalanceSnapshot, RepositoryError> {
        let net_balance = self.transaction_repo.net_balance().await?;
        let credit_total = self.transaction_repo.total_credits().await?;
        let debit_total = self.transaction_repo.total_debits().await?;

        let snapshot = BalanceSnapshot::new(net_balance, credit_total, debit_total);
        self.save(&snapshot).await?;
        Ok(snapshot)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::models::Transaction;
    use crate::storage::{Database, repository::SqliteTransactionRepository};
    use rust_decimal_macros::dec;

    async fn get_test_db() -> (Database, SqliteBalanceRepository) {
        let db = Database::in_memory()
            .await
            .expect("Failed to create test database");
        let tx_repo = SqliteTransactionRepository::new(db.pool.clone());
        let repo = SqliteBalanceRepository::new(db.pool.clone(), Box::new(tx_repo));
        (db, repo)
    }

    #[tokio::test]
    async fn should_create_snapshot() {
        let (_db, repo) = get_test_db().await;

        // Add some transactions first
        let tx_repo = SqliteTransactionRepository::new(_db.pool.clone());

        let credit = Transaction::new(
            dec!(200.00),
            TransactionType::Credit,
            None,
            "".to_string(),
            Utc::now(),
        );
        let debit = Transaction::new(
            dec!(100.00),
            TransactionType::Debit,
            None,
            "".to_string(),
            Utc::now(),
        );

        tx_repo.save(&credit).await.unwrap();
        tx_repo.save(&debit).await.unwrap();

        let snapshot = repo
            .create_snapshot()
            .await
            .expect("Failed to create snapshot");

        assert_eq!(snapshot.net_balance, dec!(100.00));
        assert_eq!(snapshot.credit_total, dec!(200.00));
        assert_eq!(snapshot.debit_total, dec!(100.00));
    }

    #[tokio::test]
    async fn should_save_and_find_latest_snapshot() {
        let (_db, repo) = get_test_db().await;

        let snapshot = BalanceSnapshot::new(dec!(500.00), dec!(800.00), dec!(300.00));
        repo.save(&snapshot).await.expect("Failed to save");

        let latest = repo.find_latest().await.expect("Failed to find latest");
        assert!(latest.is_some());
        assert_eq!(latest.unwrap().net_balance, dec!(500.00));
    }
}
