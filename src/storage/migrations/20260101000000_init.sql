-- Initial database schema for Personal Finance Tracker
-- This migration creates all the necessary tables

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id BLOB PRIMARY KEY NOT NULL,
    amount TEXT NOT NULL,
    transaction_type TEXT NOT NULL CHECK(transaction_type IN ('Credit', 'Debit')),
    category_id BLOB,
    description TEXT NOT NULL,
    date TEXT NOT NULL,
    created_at TEXT NOT NULL,
    is_deleted INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id BLOB PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE,
    color TEXT,
    is_default INTEGER NOT NULL DEFAULT 0
);

-- Balance snapshots table
CREATE TABLE IF NOT EXISTS balance_snapshots (
    id BLOB PRIMARY KEY NOT NULL,
    net_balance TEXT NOT NULL,
    credit_total TEXT NOT NULL,
    debit_total TEXT NOT NULL,
    timestamp TEXT NOT NULL
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted ON transactions(is_deleted);
CREATE INDEX IF NOT EXISTS idx_balance_snapshots_timestamp ON balance_snapshots(timestamp);

-- Default category
INSERT INTO categories (id, name, is_default) 
SELECT '00000000-0000-0000-0000-000000000000', 'Uncategorized', 1
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE is_default = 1);
