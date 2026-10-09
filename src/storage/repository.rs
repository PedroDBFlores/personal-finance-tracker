use crate::models::{BalanceSnapshot, Category, Transaction};
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use thiserror::Error;
use uuid::Uuid;

pub mod transaction;
pub mod category;
pub mod balance;

pub use transaction::SqliteTransactionRepository;
pub use category::SqliteCategoryRepository;
pub use balance::SqliteBalanceRepository;

/// Common repository error type
#[derive(Error, Debug)]
pub enum RepositoryError {
    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),
    #[error("Not found")]
    NotFound,
    #[error("Duplicate entry")]
    DuplicateEntry,
}

/// Trait for transaction repository operations
#[async_trait]
pub trait TransactionRepository: Send + Sync {
    /// Saves a transaction to the repository
    async fn save(&self, transaction: &Transaction) -> Result<(), RepositoryError>;

    /// Finds a transaction by ID
    async fn find_by_id(&self, id: Uuid) -> Result<Option<Transaction>, RepositoryError>;

    /// Finds all transactions, optionally filtered
    async fn find_all(&self) -> Result<Vec<Transaction>, RepositoryError>;

    /// Finds transactions by category
    async fn find_by_category(&self, category_id: Uuid) -> Result<Vec<Transaction>, RepositoryError>;

    /// Finds transactions by type
    async fn find_by_type(&self, transaction_type: &crate::models::transaction::TransactionType) -> Result<Vec<Transaction>, RepositoryError>;

    /// Finds transactions within a date range
    async fn find_by_date_range(
        &self,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Result<Vec<Transaction>, RepositoryError>;

    /// Updates a transaction
    async fn update(&self, transaction: &Transaction) -> Result<(), RepositoryError>;

    /// Soft deletes a transaction
    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError>;

    /// Gets the total count of transactions
    async fn count(&self) -> Result<i64, RepositoryError>;

    /// Calculates the net balance from all transactions
    async fn net_balance(&self) -> Result<Decimal, RepositoryError>;

    /// Calculates total credits
    async fn total_credits(&self) -> Result<Decimal, RepositoryError>;

    /// Calculates total debits
    async fn total_debits(&self) -> Result<Decimal, RepositoryError>;
}

/// Trait for category repository operations
#[async_trait]
pub trait CategoryRepository: Send + Sync {
    /// Saves a category to the repository
    async fn save(&self, category: &Category) -> Result<(), RepositoryError>;

    /// Finds a category by ID
    async fn find_by_id(&self, id: Uuid) -> Result<Option<Category>, RepositoryError>;

    /// Finds all categories
    async fn find_all(&self) -> Result<Vec<Category>, RepositoryError>;

    /// Finds the default category
    async fn find_default(&self) -> Result<Option<Category>, RepositoryError>;

    /// Updates a category
    async fn update(&self, category: &Category) -> Result<(), RepositoryError>;

    /// Deletes a category (only if no transactions reference it)
    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError>;

    /// Gets the total count of categories
    async fn count(&self) -> Result<i64, RepositoryError>;
}

/// Trait for balance snapshot repository operations
#[async_trait]
pub trait BalanceRepository: Send + Sync {
    /// Saves a balance snapshot
    async fn save(&self, snapshot: &BalanceSnapshot) -> Result<(), RepositoryError>;

    /// Finds the latest balance snapshot
    async fn find_latest(&self) -> Result<Option<BalanceSnapshot>, RepositoryError>;

    /// Finds balance snapshots within a date range
    async fn find_by_date_range(
        &self,
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Result<Vec<BalanceSnapshot>, RepositoryError>;

    /// Calculates and saves a new snapshot based on current transactions
    async fn create_snapshot(&self) -> Result<BalanceSnapshot, RepositoryError>;
}
