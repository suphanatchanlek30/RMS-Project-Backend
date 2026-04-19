package services

import "testing"

func TestParseTopMenusReportQuery(t *testing.T) {
	t.Run("valid query with explicit limit", func(t *testing.T) {
		query, err := ParseTopMenusReportQuery("2025-08-01", "2025-08-31", 5)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if query.Limit != 5 {
			t.Fatalf("expected limit 5, got %d", query.Limit)
		}
	})

	t.Run("default limit when omitted", func(t *testing.T) {
		query, err := ParseTopMenusReportQuery("2025-08-01", "2025-08-31", 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if query.Limit != 10 {
			t.Fatalf("expected default limit 10, got %d", query.Limit)
		}
	})

	t.Run("invalid when missing dateFrom", func(t *testing.T) {
		if _, err := ParseTopMenusReportQuery("", "2025-08-31", 10); err == nil {
			t.Fatalf("expected BAD_REQUEST for missing dateFrom")
		}
	})

	t.Run("invalid when date range reversed", func(t *testing.T) {
		if _, err := ParseTopMenusReportQuery("2025-09-01", "2025-08-31", 10); err == nil {
			t.Fatalf("expected BAD_REQUEST for reversed date range")
		}
	})

	t.Run("invalid when limit is negative", func(t *testing.T) {
		if _, err := ParseTopMenusReportQuery("2025-08-01", "2025-08-31", -1); err == nil {
			t.Fatalf("expected BAD_REQUEST for negative limit")
		}
	})

	t.Run("invalid when limit too large", func(t *testing.T) {
		if _, err := ParseTopMenusReportQuery("2025-08-01", "2025-08-31", 101); err == nil {
			t.Fatalf("expected BAD_REQUEST for limit > 100")
		}
	})
}
