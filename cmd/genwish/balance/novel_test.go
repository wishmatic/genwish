package balance

import (
	"bytes"
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/novelai"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestRunNovel(t *testing.T) {
	t.Run("prints the Anlas balance", func(t *testing.T) {
		server, _ := testutil.CaptureServer(t, http.StatusOK,
			`{"subscription":{"tier":"opus","trainingStepsLeft":{"fixedTrainingStepsLeft":9898,"purchasedTrainingSteps":32}}}`,
		)

		out := &bytes.Buffer{}
		deps := cli.Deps{Out: out, Novel: novelai.New(server.URL, "key")}

		if err := RunNovel(context.Background(), deps); err != nil {
			t.Fatalf("RunNovel() error: %v", err)
		}

		for _, want := range []string{"Anlas: 9930", "Tier: opus"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("reports a missing key", func(t *testing.T) {
		deps := cli.Deps{
			Out:         &bytes.Buffer{},
			Credentials: credential.NewFileStoreAt(filepath.Join(t.TempDir(), "key")),
		}

		err := RunNovel(context.Background(), deps)
		if err == nil || !strings.Contains(err.Error(), "genwish login novel") {
			t.Errorf("RunNovel() error = %v, want a login hint", err)
		}
	})
}
