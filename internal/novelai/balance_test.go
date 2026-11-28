package novelai

import (
	"context"
	"fmt"
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
				`{"subscription":{"tier":"opus","trainingStepsLeft":` +
					`{"fixedTrainingStepsLeft":9898,"purchasedTrainingSteps":32}}}`,
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

		if balance.Tier != "" {
			t.Errorf("Tier = %q, want it unset", balance.Tier)
		}
	})

	t.Run("maps the numeric tier", func(t *testing.T) {
		tests := map[int]string{
			0: "paper",
			1: "tablet",
			2: "scroll",
			3: "opus",
			7: "7",
		}

		for level, want := range tests {
			body := fmt.Sprintf(
				`{"subscription":{"tier":%d,"trainingStepsLeft":1200}}`,
				level,
			)

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte(body))
				}),
			)

			balance, err := New(server.URL, "key").CheckBalance(context.Background())
			server.Close()

			if err != nil {
				t.Fatalf("tier %d: CheckBalance() error: %v", level, err)
			}

			if balance.Tier != want {
				t.Errorf("tier %d = %q, want %q", level, balance.Tier, want)
			}
		}
	})

	t.Run("reports an invalid tier", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"subscription":{"tier":{},"trainingStepsLeft":1200}}`))
		}))
		t.Cleanup(server.Close)

		_, err := New(server.URL, "key").CheckBalance(context.Background())
		if err == nil || !strings.Contains(err.Error(), "subscription.tier") {
			t.Errorf("CheckBalance() error = %v, want the invalid tier named", err)
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
