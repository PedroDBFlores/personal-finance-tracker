use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};
use std::str::FromStr;
use uuid::Uuid;

/// Enum representing the type of a financial transaction
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum TransactionType {
    /// Money received (income, refunds, etc.)
    Credit,
    /// Money spent (expenses, payments, etc.)
    Debit,
}

impl std::fmt::Display for TransactionType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            TransactionType::Credit => write!(f, "Credit"),
            TransactionType::Debit => write!(f, "Debit"),
        }
    }
}

impl FromStr for TransactionType {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "credit" => Ok(TransactionType::Credit),
            "debit" => Ok(TransactionType::Debit),
            _ => Err(format!("Unknown transaction type: {}", s)),
        }
    }
}

impl TransactionType {
    /// Returns true if this is a Credit transaction
    pub fn is_credit(&self) -> bool {
        matches!(self, TransactionType::Credit)
    }

    /// Returns true if this is a Debit transaction
    pub fn is_debit(&self) -> bool {
        matches!(self, TransactionType::Debit)
    }

    /// Returns the sign multiplier for calculations
    /// Credit = +1, Debit = -1
    pub fn sign(&self) -> i8 {
        match self {
            TransactionType::Credit => 1,
            TransactionType::Debit => -1,
        }
    }
}

/// Represents a single financial transaction
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Transaction {
    /// Unique identifier for the transaction
    pub id: Uuid,
    /// The monetary amount (stored as positive value)
    pub amount: Decimal,
    /// Whether this is a credit or debit
    pub transaction_type: TransactionType,
    /// Optional category for organization
    pub category_id: Option<Uuid>,
    /// User-provided description/notes
    pub description: String,
    /// When the transaction occurred
    pub date: DateTime<Utc>,
    /// When the transaction was recorded in the system
    pub created_at: DateTime<Utc>,
    /// Whether the transaction has been marked as deleted
    pub is_deleted: bool,
}

impl Transaction {
    /// Creates a new transaction
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        amount: Decimal,
        transaction_type: TransactionType,
        category_id: Option<Uuid>,
        description: String,
        date: DateTime<Utc>,
    ) -> Self {
        Self {
            id: Uuid::new_v4(),
            amount,
            transaction_type,
            category_id,
            description,
            date,
            created_at: Utc::now(),
            is_deleted: false,
        }
    }

    /// Returns the effective amount considering the transaction type
    /// Credit amounts are positive, Debit amounts are negative
    pub fn effective_amount(&self) -> Decimal {
        if self.transaction_type.is_credit() {
            self.amount
        } else {
            -self.amount
        }
    }

    /// Marks the transaction as deleted (soft delete)
    pub fn delete(&mut self) {
        self.is_deleted = true;
    }

    /// Returns true if the transaction has been deleted
    pub fn is_deleted(&self) -> bool {
        self.is_deleted
    }

    /// Returns the absolute amount (always positive)
    pub fn absolute_amount(&self) -> Decimal {
        self.amount.abs()
    }

    /// Converts amount to string for database storage
    pub fn amount_to_string(&self) -> String {
        self.amount.to_string()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rust_decimal_macros::dec;

    #[test]
    fn should_create_credit_transaction() {
        let transaction = Transaction::new(
            dec!(100.00),
            TransactionType::Credit,
            None,
            "Salary".to_string(),
            Utc::now(),
        );

        assert!(transaction.id != Uuid::nil());
        assert_eq!(transaction.amount, dec!(100.00));
        assert!(transaction.transaction_type.is_credit());
        assert_eq!(transaction.effective_amount(), dec!(100.00));
        assert!(!transaction.is_deleted());
    }

    #[test]
    fn should_create_debit_transaction() {
        let transaction = Transaction::new(
            dec!(50.00),
            TransactionType::Debit,
            None,
            "Coffee".to_string(),
            Utc::now(),
        );

        assert!(transaction.id != Uuid::nil());
        assert_eq!(transaction.amount, dec!(50.00));
        assert!(transaction.transaction_type.is_debit());
        assert_eq!(transaction.effective_amount(), dec!(-50.00));
        assert!(!transaction.is_deleted());
    }

    #[test]
    fn should_calculate_effective_amount_correctly() {
        let credit = Transaction::new(
            dec!(200.00),
            TransactionType::Credit,
            None,
            "".to_string(),
            Utc::now(),
        );
        let debit = Transaction::new(
            dec!(75.50),
            TransactionType::Debit,
            None,
            "".to_string(),
            Utc::now(),
        );

        assert_eq!(credit.effective_amount(), dec!(200.00));
        assert_eq!(debit.effective_amount(), dec!(-75.50));
    }

    #[test]
    fn should_soft_delete_transaction() {
        let mut transaction = Transaction::new(
            dec!(10.00),
            TransactionType::Debit,
            None,
            "".to_string(),
            Utc::now(),
        );

        assert!(!transaction.is_deleted());
        transaction.delete();
        assert!(transaction.is_deleted());
    }

    #[test]
    fn should_return_absolute_amount() {
        let credit = Transaction::new(
            dec!(100.00),
            TransactionType::Credit,
            None,
            "".to_string(),
            Utc::now(),
        );
        let debit = Transaction::new(
            dec!(50.00),
            TransactionType::Debit,
            None,
            "".to_string(),
            Utc::now(),
        );

        assert_eq!(credit.absolute_amount(), dec!(100.00));
        assert_eq!(debit.absolute_amount(), dec!(50.00));
    }

    #[test]
    fn should_have_automatic_creation_timestamp() {
        let before = Utc::now();
        let transaction = Transaction::new(
            dec!(1.00),
            TransactionType::Credit,
            None,
            "".to_string(),
            Utc::now(),
        );
        let after = Utc::now();

        assert!(transaction.created_at >= before);
        assert!(transaction.created_at <= after);
    }

    #[test]
    fn should_parse_transaction_type() {
        assert_eq!("Credit".parse::<TransactionType>().unwrap(), TransactionType::Credit);
        assert_eq!("credit".parse::<TransactionType>().unwrap(), TransactionType::Credit);
        assert_eq!("Debit".parse::<TransactionType>().unwrap(), TransactionType::Debit);
        assert_eq!("debit".parse::<TransactionType>().unwrap(), TransactionType::Debit);
        assert!("Invalid".parse::<TransactionType>().is_err());
    }

    #[test]
    fn should_display_transaction_type() {
        assert_eq!(format!("{}", TransactionType::Credit), "Credit");
        assert_eq!(format!("{}", TransactionType::Debit), "Debit");
    }
}
