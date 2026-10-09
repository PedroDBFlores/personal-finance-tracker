use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode, SqlitePoolOptions, SqliteSynchronous};
use sqlx::{Pool, Sqlite};
use std::path::Path;
use thiserror::Error;

/// Configuration for the SQLite database
#[derive(Debug, Clone)]
pub struct DatabaseConfig {
    /// Path to the SQLite database file
    pub path: String,
}

impl Default for DatabaseConfig {
    fn default() -> Self {
        Self {
            path: "personal_finance_tracker.db".to_string(),
        }
    }
}

impl DatabaseConfig {
    /// Creates a new database configuration with the specified path
    pub fn new<P: AsRef<Path>>(path: P) -> Self {
        Self {
            path: path.as_ref().to_string_lossy().into_owned(),
        }
    }

    /// Returns the connection string for SQLite
    pub fn connection_string(&self) -> String {
        format!("sqlite:{}", self.path)
    }
}

/// Main database struct that holds the connection pool
#[derive(Debug)]
pub struct Database {
    pub pool: Pool<Sqlite>,
}

impl Database {
    /// Creates a new database instance with the given configuration
    pub async fn new(config: &DatabaseConfig) -> Result<Self, DatabaseError> {
        // Ensure the parent directory exists
        if let Some(parent) = Path::new(&config.path).parent() {
            if !parent.as_os_str().is_empty() {
                tokio::fs::create_dir_all(parent).await?;
            }
        }

        // Configure SQLite connection options for better performance
        let connect_options = SqliteConnectOptions::new()
            .filename(&config.path)
            .create_if_missing(true)
            .journal_mode(SqliteJournalMode::Wal)
            .synchronous(SqliteSynchronous::Normal)
            .foreign_keys(true);

        // Create connection pool
        let pool = SqlitePoolOptions::new()
            .max_connections(5)
            .connect_with(connect_options)
            .await?;

        // Run migrations
        Self::run_migrations(&pool, &config.path).await?;

        Ok(Self { pool })
    }

    /// Runs database migrations using runtime SQL
    async fn run_migrations(pool: &Pool<Sqlite>, db_path: &str) -> Result<(), DatabaseError> {
        // For in-memory databases, always run migrations
        // For file databases, check if tables exist
        let force_migrations = db_path == ":memory:";
        
        let categories_exist: Option<i64> = if !force_migrations {
            match sqlx::query_scalar::<_, Option<i64>>(
                "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'categories'"
            )
            .fetch_optional(pool)
            .await {
                Ok(result) => result.flatten(),
                Err(_) => None,
            }
        } else {
            None
        };

        if force_migrations || categories_exist != Some(1) {
            // Create categories table first (transactions references it)
            sqlx::query(
                r#"
                CREATE TABLE categories (
                    id BLOB PRIMARY KEY NOT NULL,
                    name TEXT NOT NULL UNIQUE,
                    color TEXT,
                    is_default INTEGER NOT NULL DEFAULT 0
                )"#
            )
            .execute(pool)
            .await?;

            // Create transactions table (foreign key will be added after categories exists)
            sqlx::query(
                r#"
                CREATE TABLE transactions (
                    id BLOB PRIMARY KEY NOT NULL,
                    amount TEXT NOT NULL,
                    transaction_type TEXT NOT NULL CHECK(transaction_type IN ('Credit', 'Debit')),
                    category_id BLOB,
                    description TEXT NOT NULL,
                    date TEXT NOT NULL,
                    created_at TEXT NOT NULL,
                    is_deleted INTEGER NOT NULL DEFAULT 0
                )"#
            )
            .execute(pool)
            .await?;

            // Create balance_snapshots table
            sqlx::query(
                r#"
                CREATE TABLE balance_snapshots (
                    id BLOB PRIMARY KEY NOT NULL,
                    net_balance TEXT NOT NULL,
                    credit_total TEXT NOT NULL,
                    debit_total TEXT NOT NULL,
                    timestamp TEXT NOT NULL
                )"#
            )
            .execute(pool)
            .await?;

            // Create indexes
            sqlx::query("CREATE INDEX idx_transactions_category ON transactions(category_id)")
                .execute(pool)
                .await?;

            sqlx::query("CREATE INDEX idx_transactions_type ON transactions(transaction_type)")
                .execute(pool)
                .await?;

            sqlx::query("CREATE INDEX idx_transactions_date ON transactions(date)")
                .execute(pool)
                .await?;

            sqlx::query("CREATE INDEX idx_balance_snapshots_timestamp ON balance_snapshots(timestamp)")
                .execute(pool)
                .await?;

            // Insert default category
            let default_uuid = uuid::Uuid::parse_str("00000000-0000-0000-0000-000000000000").unwrap();
            sqlx::query(
                r#"
                INSERT INTO categories (id, name, is_default) 
                SELECT $1, 'Uncategorized', 1
                WHERE NOT EXISTS (SELECT 1 FROM categories WHERE is_default = 1)
                "#
            )
            .bind(default_uuid)
            .execute(pool)
            .await?;
        }

        Ok(())
    }

    /// Returns a new database with in-memory storage (for testing)
    pub async fn in_memory() -> Result<Self, DatabaseError> {
        let connect_options = SqliteConnectOptions::new()
            .filename(":memory:")
            .create_if_missing(true)
            .foreign_keys(true)
            .shared_cache(true);

        let pool = SqlitePoolOptions::new()
            .max_connections(1)  // Use single connection for in-memory to avoid pooling issues
            .connect_with(connect_options)
            .await?;

        // Run migrations on in-memory database
        Self::run_migrations(&pool, ":memory:").await?;

        Ok(Self { pool })
    }
}

/// Database-related errors
#[derive(Error, Debug)]
pub enum DatabaseError {
    #[error("SQLx error: {0}")]
    Sqlx(#[from] sqlx::Error),
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    #[error("Migration error: {0}")]
    Migration(String),
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn should_create_in_memory_database() {
        let db = Database::in_memory().await;
        assert!(db.is_ok());
    }

    #[tokio::test]
    async fn should_have_connection_string() {
        let config = DatabaseConfig::new("test.db");
        assert_eq!(config.connection_string(), "sqlite:test.db");
    }

    #[tokio::test]
    async fn should_use_default_path() {
        let config = DatabaseConfig::default();
        assert_eq!(config.path, "personal_finance_tracker.db");
    }
}
