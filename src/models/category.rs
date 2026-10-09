use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Represents a category for organizing transactions
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Category {
    /// Unique identifier for the category
    pub id: Uuid,
    /// Display name of the category
    pub name: String,
    /// Optional color code for TUI display (e.g., "#ff0000" or "red")
    pub color: Option<String>,
    /// Whether this is the default category for new transactions
    pub is_default: bool,
}

impl Category {
    /// Creates a new category
    pub fn new(name: String, color: Option<String>) -> Self {
        Self {
            id: Uuid::new_v4(),
            name,
            color,
            is_default: false,
        }
    }

    /// Creates a new default category
    pub fn new_default(name: String, color: Option<String>) -> Self {
        Self {
            id: Uuid::new_v4(),
            name,
            color,
            is_default: true,
        }
    }

    /// Updates the category name
    pub fn set_name(&mut self, name: String) {
        self.name = name;
    }

    /// Updates the category color
    pub fn set_color(&mut self, color: Option<String>) {
        self.color = color;
    }

    /// Sets whether this is the default category
    pub fn set_default(&mut self, is_default: bool) {
        self.is_default = is_default;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn should_create_category() {
        let category = Category::new("Food".to_string(), None);

        assert!(category.id != Uuid::nil());
        assert_eq!(category.name, "Food");
        assert_eq!(category.color, None);
        assert!(!category.is_default);
    }

    #[test]
    fn should_create_default_category() {
        let category = Category::new_default("Uncategorized".to_string(), None);

        assert!(category.id != Uuid::nil());
        assert_eq!(category.name, "Uncategorized");
        assert!(category.is_default);
    }

    #[test]
    fn should_create_category_with_color() {
        let category = Category::new("Food".to_string(), Some("#ff0000".to_string()));

        assert_eq!(category.color, Some("#ff0000".to_string()));
    }

    #[test]
    fn should_update_category_name() {
        let mut category = Category::new("Food".to_string(), None);
        category.set_name("Groceries".to_string());

        assert_eq!(category.name, "Groceries");
    }

    #[test]
    fn should_update_category_color() {
        let mut category = Category::new("Food".to_string(), None);
        category.set_color(Some("#00ff00".to_string()));

        assert_eq!(category.color, Some("#00ff00".to_string()));
    }

    #[test]
    fn should_set_default_category() {
        let mut category = Category::new("Food".to_string(), None);
        assert!(!category.is_default);

        category.set_default(true);
        assert!(category.is_default);

        category.set_default(false);
        assert!(!category.is_default);
    }
}
