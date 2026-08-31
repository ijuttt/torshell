package test

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"torshell/pkg/sysutil"
)

func TestGenerateID(t *testing.T) {
	t.Run("default length", func(t *testing.T) {
		id, err := sysutil.GenerateID(sysutil.DefaultIDLen)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(id) != sysutil.DefaultIDLen {
			t.Errorf("got len %d, want %d", len(id), sysutil.DefaultIDLen)
		}
		if _, err := hex.DecodeString(id); err != nil {
			t.Errorf("invalid hex: %v", err)
		}
	})

	t.Run("odd length", func(t *testing.T) {
		id, err := sysutil.GenerateID(7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(id) != 7 {
			t.Errorf("got len %d, want 7", len(id))
		}
	})

	t.Run("invalid length", func(t *testing.T) {
		for _, n := range []int{0, -1, -5} {
			if _, err := sysutil.GenerateID(n); err == nil {
				t.Errorf("expected error for len %d, got nil", n)
			}
		}
	})
}

func TestGenerateID_Uniqueness(t *testing.T) {
	const count = 10000
	seen := make(map[string]struct{}, count)

	for i := 0; i < count; i++ {
		id, err := sysutil.GenerateID(sysutil.DefaultIDLen)
		if err != nil {
			t.Fatalf("iter %d error: %v", i, err)
		}
		if _, exists := seen[id]; exists {
			t.Fatalf("collision at iter %d: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}

func TestVethPairNames(t *testing.T) {
	t.Run("valid 8-char id", func(t *testing.T) {
		host, ns, err := sysutil.VethPairNames("a1b2c3d4")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if host != "vh-a1b2c3d4" || ns != "vn-a1b2c3d4" {
			t.Errorf("got %s, %s; want vh-a1b2c3d4, vn-a1b2c3d4", host, ns)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		if _, _, err := sysutil.VethPairNames(""); !errors.Is(err, sysutil.ErrEmptyID) {
			t.Errorf("got %v, want %v", err, sysutil.ErrEmptyID)
		}
	})

	t.Run("boundary checks", func(t *testing.T) {
		// 12 chars + 3 ("vh-") = 15 (max allowed)
		if _, _, err := sysutil.VethPairNames(strings.Repeat("a", 12)); err != nil {
			t.Errorf("len 15 should be allowed, got error: %v", err)
		}

		// 13 chars + 3 ("vh-") = 16 (exceeds max 15)
		if _, _, err := sysutil.VethPairNames(strings.Repeat("a", 13)); !errors.Is(err, sysutil.ErrNameTooLong) {
			t.Errorf("len 16 should fail, got %v", err)
		}
	})
}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = sysutil.GenerateID(sysutil.DefaultIDLen)
	}
}

func BenchmarkVethPairNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = sysutil.VethPairNames("a1b2c3d4")
	}
}
