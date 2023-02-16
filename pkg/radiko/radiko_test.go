package radiko

import (
	"log"
	"os"
	"testing"

	"github.com/moutend/go-radiko/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetPartialKey(t *testing.T) {
	partialKey, err := GetPartialKey()

	require.NoError(t, err)
	require.Equal(t, 40, len(partialKey))
}

func TestLoginWithWebForm(t *testing.T) {
	t.Parallel()

	s := NewSession(os.Getenv("RADIKO_USERNAME"), os.Getenv("RADIKO_PASSWORD"))

	if s.username == "" || s.password == "" {
		t.Logf("Skip this test")

		return
	}

	if true || testutil.IsVerbose() {
		s.SetLogger(log.New(os.Stdout, "test: ", 0))
	}

	require.NoError(t, s.loginWithWebForm())
}
