use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};
use std::str::FromStr;
use uuid::Uuid;

/// Represents a snapshot of the financial balance at a point in time
/// Used for historical tracking and performance optimization
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct BalanceSnapshot {
    /// Unique identifier for the snapshot
    pub id: Uuid,
    /// The net balance (credits - debits)
    pub net_balance: Decimal,
    /// Total of all credit transactions
    pub credit_total: Decimal,
    /// Total of all debit transactions
    pub debit_total: Decimal,
    /// When the snapshot was taken
    pub timestamp: DateTime<Utc>,
}

impl BalanceSnapshot {
    /// Creates a new balance snapshot
    pub fn new(net_balance: Decimal, credit_total: Decimal, debit_total: Decimal) -> Self {
        Self {
            id: Uuid::new_v4(),
            net_balance,
            credit_total,
            debit_total,
            timestamp: Utc::now(),
        }
    }

    /// Returns true if the net balance is positive (more credits than debits)
    pub fn is_positive(&self) -> bool {
        self.net_balance > Decimal::ZERO
    }

    /// Returns true if the net balance is negative (more debits than credits)
    pub fn is_negative(&self) -> bool {
        self.net_balance < Decimal::ZERO
    }

    /// Returns true if the net balance is zero
    pub fn is_zero(&self) -> bool {
        self.net_balance == Decimal::ZERO
    }

    /// Calculates the credit-to-debit ratio
    pub fn credit_debit_ratio(&self) -> Option<Decimal> {
        if self.debit_total == Decimal::ZERO {
            if self.credit_total == Decimal::ZERO {
                Some(Decimal::ONE)
            } else {
                None // Division by zero
            }
        } else {
            Some(self.credit_total / self.debit_total)
        }
    }

    /// Converts decimal fields to strings for database storage
    pub fn to_strings(&self) -> (String, String, String) {
        (
            self.net_balance.to_string(),
            self.credit_total.to_string(),
            self.debit_total.to_string(),
        )
    }

    /// Creates a BalanceSnapshot from database row data
    pub fn from_strings(
        id: Uuid,
        net_balance: String,
        credit_total: String,
        debit_total: String,
        timestamp: DateTime<Utc>,
    ) -> Result<Self, rust_decimal::Error> {
        let net_balance = Decimal::from_str(&net_balance)?;
        let credit_total = Decimal::from_str(&credit_total)?;
        let debit_total = Decimal::from_str(&debit_total)?;

        Ok(Self {
            id,
            net_balance,
            credit_total,
            debit_total,
            timestamp,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rust_decimal_macros::dec;

    #[test]
    fn should_create_balance_snapshot() {
        let snapshot = BalanceSnapshot::new(dec!(1000.00), dec!(1500.00), dec!(500.00));

        assert!(snapshot.id != Uuid::nil());
        assert_eq!(snapshot.net_balance, dec!(1000.00));
        assert_eq!(snapshot.credit_total, dec!(1500.00));
        assert_eq!(snapshot.debit_total, dec!(500.00));
    }

    #[test]
    fn should_detect_positive_balance() {
        let positive = BalanceSnapshot::new(dec!(100.00), dec!(200.00), dec!(100.00));
        let negative = BalanceSnapshot::new(dec!(-100.00), dec!(100.00), dec!(200.00));
        let zero = BalanceSnapshot::new(Decimal::ZERO, dec!(100.00), dec!(100.00));

        assert!(positive.is_positive());
        assert!(!positive.is_negative());
        assert!(!positive.is_zero());

        assert!(!negative.is_positive());
        assert!(negative.is_negative());
        assert!(!negative.is_zero());

        assert!(!zero.is_positive());
        assert!(!zero.is_negative());
        assert!(zero.is_zero());
    }

    #[test]
    fn should_calculate_credit_debit_ratio() {
        // Credit: 200, Debit: 100, Ratio: 2.0
        let snapshot1 = BalanceSnapshot::new(dec!(100.00), dec!(200.00), dec!(100.00));
        assert_eq!(snapshot1.credit_debit_ratio(), Some(dec!(2.0)));

        // Credit: 0, Debit: 100, Ratio: 0 (zero divided by non-zero)
        let snapshot2 = BalanceSnapshot::new(dec!(-100.00), Decimal::ZERO, dec!(100.00));
        assert_eq!(snapshot2.credit_debit_ratio(), Some(dec!(0.0)));

        // Credit: 0, Debit: 0, Ratio: Some(1.0) (special case)
        let snapshot3 = BalanceSnapshot::new(Decimal::ZERO, Decimal::ZERO, Decimal::ZERO);
        assert_eq!(snapshot3.credit_debit_ratio(), Some(Decimal::ONE));

        // Credit: 150, Debit: 100, Ratio: 1.5
        let snapshot4 = BalanceSnapshot::new(dec!(50.00), dec!(150.00), dec!(100.00));
        assert_eq!(snapshot4.credit_debit_ratio(), Some(dec!(1.5)));
    }

    #[test]
    fn should_convert_to_and_from_strings() {
        let original = BalanceSnapshot::new(dec!(100.00), dec!(200.00), dec!(50.00));
        let (net_str, credit_str, debit_str) = original.to_strings();
        
        let reconstructed = BalanceSnapshot::from_strings(
            original.id,
            net_str,
            credit_str,
            debit_str,
            original.timestamp,
        ).unwrap();

        assert_eq!(reconstructed.net_balance, original.net_balance);
        assert_eq!(reconstructed.credit_total, original.credit_total);
        assert_eq!(reconstructed.debit_total, original.debit_total);
    }
}
