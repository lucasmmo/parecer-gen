package parecer

import (
	"testing"
	"time"
)

func TestNewData(t *testing.T) {
	t.Run("should return error if user is empty", func(t *testing.T) {
		_, err := NewData("", "creci", "content", time.Now())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should return error if creci is empty", func(t *testing.T) {
		_, err := NewData("user", "", "content", time.Now())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should return error if content is empty", func(t *testing.T) {
		_, err := NewData("user", "creci", "", time.Now())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should set date to now if date is zero value", func(t *testing.T) {
		data, _ := NewData("user", "creci", "content", time.Time{})
		if data.Date.IsZero() {
			t.Error("expected date to be set, got zero value")
		}
	})

}
