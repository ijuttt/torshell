package sysutil

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	// max ifname length in linux (IFNAMSIZ - 1)
	MaxIfNameLen = 15
	DefaultIDLen = 8
)

var (
	ErrEmptyID     = errors.New("id cannot be empty")
	ErrNameTooLong = errors.New("interface name exceeds linux limit of 15 chars")
)

// GenerateID produces a random hex string of length n.
func GenerateID(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("id length must be > 0")
	}

	b := make([]byte, (n+1)/2)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto rand failed: %w", err)
	}

	return hex.EncodeToString(b)[:n], nil
}

// VethPairNames returns "vh-<id>" (host) and "vn-<id>" (namespace),
// checking that names fit within linux IFNAMSIZ limit (<= 15 chars).
func VethPairNames(id string) (hostVeth, nsVeth string, error) {
	if id == "" {
		return "", "", ErrEmptyID
	}

	hostVeth = "vh-" + id
	nsVeth = "vn-" + id

	if len(hostVeth) > MaxIfNameLen || len(nsVeth) > MaxIfNameLen {
		return "", "", fmt.Errorf("%w (len=%d): %s", ErrNameTooLong, len(hostVeth), hostVeth)
	}

	return hostVeth, nsVeth, nil
}
