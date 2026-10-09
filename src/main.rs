use tracing_subscriber::{fmt, EnvFilter};

mod app;
mod models;
mod storage;
mod tui;



#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logging
    fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    tracing::info!("Starting Personal Finance Tracker...");

    // Initialize database
    let db_config = storage::database::DatabaseConfig::default();
    let db = storage::Database::new(&db_config).await?;

    tracing::info!("Database initialized");

    // Create repositories
    let transaction_repo = storage::repository::SqliteTransactionRepository::new(db.pool.clone());
    let category_repo = storage::repository::SqliteCategoryRepository::new(db.pool.clone());
    let balance_repo = storage::repository::SqliteBalanceRepository::new(
        db.pool.clone(),
        Box::new(storage::repository::SqliteTransactionRepository::new(db.pool.clone()))
    );

    // Initialize application state
    let app_state = app::state::AppState::new(
        Box::new(transaction_repo),
        Box::new(category_repo),
        Box::new(balance_repo),
    );

    // Start TUI
    let mut tui_app = tui::app::TuiApp::new(app_state);
    tui_app.init().await?;
    tui_app.run().await?;

    Ok(())
}
