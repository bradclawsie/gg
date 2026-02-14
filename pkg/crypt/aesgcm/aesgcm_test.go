package aesgcm

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAESGCM(t *testing.T) {
	t.Run("Encrypt", func(t *testing.T) {
		t.Parallel()
		k := NewKey()
		s := uuid.NewString()
		e, err := Encrypt(s, k)
		require.NoError(t, err, "encrypt fail")
		digestBytes := sha256.Sum256([]byte(s))
		digest := hex.EncodeToString(digestBytes[:])
		d, err := Decrypt(e, digest, k)
		require.NoError(t, err, "decrypt fail")
		require.Equal(t, s, d, "round trip")

		// Bad key.
		_, err = Decrypt(e, digest, NewKey())
		require.Error(t, err, "bad key")

		// Bad digest.
		_, err = Decrypt(e, "abcd", k)
		require.Error(t, err, "bad digest")
		require.Equal(t, ErrDigest, err, "digest err")
	})
}
