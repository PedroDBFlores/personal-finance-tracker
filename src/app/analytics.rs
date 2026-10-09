use crate::models::{Transaction, BalanceSnapshot};
use crate::models::transaction::TransactionType;
use chrono::{DateTime, Utc, Datelike};
use rust_decimal::Decimal;
use std::collections::HashMap;
use uuid::Uuid;

/// Service for calculating financial analytics
#[derive(Debug)]
pub struct AnalyticsService;

impl AnalyticsService {
    pub fn new() -> Self {
        Self
    }

    /// Calculates net balance from a list of transactions
    pub fn calculate_net_balance(transactions: &[Transaction]) -> Decimal {
        transactions.iter()
            .map(|t| t.effective_amount())
            .fold(Decimal::ZERO, |acc, amount| acc + amount)
    }

    /// Calculates total credits from a list of transactions
    pub fn calculate_total_credits(transactions: &[Transaction]) -> Decimal {
        transactions.iter()
            .filter(|t| t.transaction_type.is_credit())
            .map(|t| t.amount)
            .fold(Decimal::ZERO, |acc, amount| acc + amount)
    }

    /// Calculates total debits from a list of transactions
    pub fn calculate_total_debits(transactions: &[Transaction]) -> Decimal {
        transactions.iter()
            .filter(|t| t.transaction_type.is_debit())
            .map(|t| t.amount)
            .fold(Decimal::ZERO, |acc, amount| acc + amount)
    }

    /// Groups transactions by category and calculates totals
    pub fn by_category(transactions: &[Transaction]) -> HashMap<Option<Uuid>, Decimal> {
        let mut result = HashMap::new();
        
        for transaction in transactions {
            let entry = result.entry(transaction.category_id).or_insert(Decimal::ZERO);
            *entry += transaction.effective_amount();
        }
        
        result
    }

    /// Groups transactions by month and calculates totals
    pub fn by_month(transactions: &[Transaction]) -> HashMap<String, Decimal> {
        let mut result = HashMap::new();
        
        for transaction in transactions {
            let month_key = format!("{}-{:02}", transaction.date.year(), transaction.date.month());
            let entry = result.entry(month_key).or_insert(Decimal::ZERO);
            *entry += transaction.effective_amount();
        }
        
        result
    }

    /// Calculates the current balance trend based on recent snapshots
    pub fn balance_trend(snapshots: &[BalanceSnapshot]) -> Decimal {
        if snapshots.len() < 2 {
            return Decimal::ZERO;
        }
        
        let newest = &snapshots[0];
        let oldest = &snapshots[snapshots.len() - 1];
        
        newest.net_balance - oldest.net_balance
    }

    /// Filters transactions by date range
    pub fn filter_by_date_range(
        transactions: &[Transaction],
        start: DateTime<Utc>,
        end: DateTime<Utc>,
    ) -> Vec<Transaction> {
        transactions.iter()
            .filter(|t| t.date >= start && t.date <= end)
            .cloned()
            .collect()
    }

    /// Filters transactions by type
    pub fn filter_by_type(
        transactions: &[Transaction],
        transaction_type: &crate::models::transaction::TransactionType,
    ) -> Vec<Transaction> {
        transactions.iter()
            .filter(|t| &t.transaction_type == transaction_type)
            .cloned()
            .collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rust_decimal_macros::dec;
    use chrono::TimeZone;

    #[test]
    fn should_calculate_net_balance() {
        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(50.00), TransactionType::Debit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(200.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
        ];

        let net = AnalyticsService::calculate_net_balance(&transactions);
        assert_eq!(net, dec!(250.00));
    }

    #[test]
    fn should_calculate_total_credits() {
        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(50.00), TransactionType::Debit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(200.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
        ];

        let total = AnalyticsService::calculate_total_credits(&transactions);
        assert_eq!(total, dec!(300.00));
    }

    #[test]
    fn should_calculate_total_debits() {
        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(50.00), TransactionType::Debit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(200.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(75.00), TransactionType::Debit, None, "".to_string(), Utc::now()),
        ];

        let total = AnalyticsService::calculate_total_debits(&transactions);
        assert_eq!(total, dec!(125.00));
    }

    #[test]
    fn should_group_by_category() {
        let category1 = uuid::Uuid::new_v4();
        let category2 = uuid::Uuid::new_v4();

        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, Some(category1), "".to_string(), Utc::now()),
            Transaction::new(dec!(50.00), TransactionType::Debit, Some(category1), "".to_string(), Utc::now()),
            Transaction::new(dec!(200.00), TransactionType::Credit, Some(category2), "".to_string(), Utc::now()),
        ];

        let by_category = AnalyticsService::by_category(&transactions);
        
        assert_eq!(by_category.len(), 2);
        
        // Category 1: 100 - 50 = 50
        if let Some(cat1_total) = by_category.get(&Some(category1)) {
            assert_eq!(*cat1_total, dec!(50.00));
        } else {
            panic!("Category 1 not found");
        }
        
        // Category 2: 200
        if let Some(cat2_total) = by_category.get(&Some(category2)) {
            assert_eq!(*cat2_total, dec!(200.00));
        } else {
            panic!("Category 2 not found");
        }
    }

    #[test]
    fn should_filter_by_type() {
        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(50.00), TransactionType::Debit, None, "".to_string(), Utc::now()),
            Transaction::new(dec!(200.00), TransactionType::Credit, None, "".to_string(), Utc::now()),
        ];

        let credits = AnalyticsService::filter_by_type(&transactions, &TransactionType::Credit);
        let debits = AnalyticsService::filter_by_type(&transactions, &TransactionType::Debit);

        assert_eq!(credits.len(), 2);
        assert_eq!(debits.len(), 1);
    }

    #[test]
    fn should_filter_by_date_range() {
        let start = Utc.with_ymd_and_hms(2026, 1, 1, 0, 0, 0).unwrap();
        let mid = Utc.with_ymd_and_hms(2026, 1, 15, 0, 0, 0).unwrap();
        let end = Utc.with_ymd_and_hms(2026, 2, 1, 0, 0, 0).unwrap();
        let after = Utc.with_ymd_and_hms(2026, 2, 15, 0, 0, 0).unwrap();

        let transactions = vec![
            Transaction::new(dec!(100.00), TransactionType::Credit, None, "".to_string(), start),
            Transaction::new(dec!(50.00), TransactionType::Debit, None, "".to_string(), mid),
            Transaction::new(dec!(200.00), TransactionType::Credit, None, "".to_string(), end),
            Transaction::new(dec!(75.00), TransactionType::Debit, None, "".to_string(), after),
        ];

        let filtered = AnalyticsService::filter_by_date_range(&transactions, start, end);
        assert_eq!(filtered.len(), 3);
    }
}
