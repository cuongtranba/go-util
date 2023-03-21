package interval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	next3Min := time.Now().Add(10 * time.Second)
	result := Retry(1*time.Second, 3, func() (bool, error) {
		n := time.Now()
		if n.After(next3Min) {
			return true, nil
		}
		return false, nil
	})
	a := <-result
	require.Nil(t, a.Err)
	require.True(t, a.IsStop)
	require.True(t, a.IsMaxRetryReached)
}
