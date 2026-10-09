use crate::models::{BalanceSnapshot, Category, Transaction};
use crate::storage::repository::{BalanceRepository, CategoryRepository, RepositoryError, TransactionRepository};
use rust_decimal::Decimal;
use std::sync::Arc;
use tokio::sync::Mutex;

/// Application state that holds all repositories and current data
pub struct AppState {
    transaction_repo: Box<dyn TransactionRepository>,
    category_repo: Box<dyn CategoryRepository>,
    balance_repo: Box<dyn BalanceRepository>,
}

impl AppState {
    pub fn new(
        transaction_repo: Box<dyn TransactionRepository>,
        category_repo: Box<dyn CategoryRepository>,
        balance_repo: Box<dyn BalanceRepository>,
    ) -> Arc<Mutex<Self>> {
        Arc::new(Mutex::new(Self {
            transaction_repo,
            category_repo,
            balance_repo,
        }))
    }

    // Delegation methods for transaction operations
    pub async fn save_transaction(&self, transaction: &Transaction) -> Result<(), RepositoryError> {
        self.transaction_repo.save(transaction).await
    }

    pub async fn find_transaction(&self, id: uuid::Uuid) -> Result<Option<Transaction>, RepositoryError> {
        self.transaction_repo.find_by_id(id).await
    }

    pub async fn all_transactions(&self) -> Result<Vec<Transaction>, RepositoryError> {
        self.transaction_repo.find_all().await
    }

    pub async fn net_balance(&self) -> Result<Decimal, RepositoryError> {
        self.transaction_repo.net_balance().await
    }

    // Delegation methods for category operations
    pub async fn save_category(&self, category: &Category) -> Result<(), RepositoryError> {
        self.category_repo.save(category).await
    }

    pub async fn all_categories(&self) -> Result<Vec<Category>, RepositoryError> {
        self.category_repo.find_all().await
    }

    // Delegation methods for balance operations
    pub async fn create_snapshot(&self) -> Result<BalanceSnapshot, RepositoryError> {
        self.balance_repo.create_snapshot().await
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::storage::{Database, repository::SqliteTransactionRepository};
    use crate::models::Transaction;
    use rust_decimal_macros::dec;
    use chrono::Utc;

    #[tokio::test]
    async fn should_create_app_state() {
        let db = Database::in_memory().await.unwrap();
        let tx_repo = Box::new(SqliteTransactionRepository::new(db.pool.clone()));
        let category_repo = Box::new(crate::storage::repository::SqliteCategoryRepository::new(db.pool.clone()));
        let balance_repo = Box::new(crate::storage::repository::SqliteBalanceRepository::new(
            db.pool.clone(),
            Box::new(SqliteTransactionRepository::new(db.pool.clone()))
        ));

        let state = AppState::new(tx_repo, category_repo, balance_repo);
        assert!(Arc::strong_count(&state) > 0);
    }
}
