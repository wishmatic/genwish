package novelai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckBalance(t *testing.T) {
	t.Run("reads the object shape", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != userDataPath {
				t.Errorf("path = %q, want %q", r.URL.Path, userDataPath)
			}
			if got := r.Header.Get("Authorization"); got != "Bearer key" {
				t.Errorf("Authorization = %q, want Bearer key", got)
			}

			_, _ = w.Write([]byte(
				`{"subscription":{"tier":"opus","trainingStepsLeft":{"fixedTrainingStepsLeft":9898,"purchasedTrainingSteps":32}}}`,
			))
		}))
		t.Cleanup(server.Close)

		balance, err := New(server.URL, "key").CheckBalance(context.Background())
		if err != nil {
			t.Fatalf("CheckBalance() error: %v", err)
		}

		if balance.Anlas != 9930 || balance.Tier != "opus" {
			t.Errorf("balance = %+v, want 9930 Anlas on opus", balance)
		}
	})

	t.Run("reads the integer shape", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"subscription":{"trainingStepsLeft":1200}}`))
		}))
		t.Cleanup(server.Close)

		balance, err := New(server.URL, "key").CheckBalance(context.Background())
		if err != nil {
			t.Fatalf("CheckBalance() error: %v", err)
		}

		if balance.Anlas != 1200 {
			t.Errorf("Anlas = %d, want 1200", balance.Anlas)
		}
	})

	t.Run("reports a missing field", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"subscription":{}}`))
		}))
		t.Cleanup(server.Close)

		_, err := New(server.URL, "key").CheckBalance(context.Background())
		if err == nil || !strings.Contains(err.Error(), "trainingStepsLeft") {
			t.Errorf("CheckBalance() error = %v, want the missing field named", err)
		}
	})

	t.Run("reports unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "bad key", http.StatusUnauthorized)
		}))
		t.Cleanup(server.Close)

		_, err := New(server.URL, "key").CheckBalance(context.Background())
		if err == nil || !strings.Contains(err.Error(), "401") {
			t.Errorf("CheckBalance() error = %v, want the status surfaced", err)
		}
	})
}
