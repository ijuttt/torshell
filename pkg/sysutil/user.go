package sysutil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

var ErrNotRoot = errors.New("torshell needs root (run with sudo)")

// UserContext holds identity info of the user running torshell.
type UserContext struct {
	UID      int
	GID      int
	Username string
	HomeDir  string
	IsSudo   bool
}

// IsRoot returns true if running as root (euid 0).
func IsRoot() bool {
	return os.Geteuid() == 0
}

// RequireRoot returns an error if not running as root.
func RequireRoot() error {
	if !IsRoot() {
		return ErrNotRoot
	}
	return nil
}

// GetRealUser finds who is actually running torshell,
// handling sudo env vars (SUDO_UID/SUDO_GID) or falling back to current user.
func GetRealUser() (*UserContext, error) {
	sudoUID := os.Getenv("SUDO_UID")
	sudoGID := os.Getenv("SUDO_GID")
	sudoUser := os.Getenv("SUDO_USER")

	// if invoked via sudo, use original caller's UID/GID
	if sudoUID != "" && sudoGID != "" {
		uid, err := strconv.Atoi(sudoUID)
		if err != nil {
			return nil, fmt.Errorf("bad SUDO_UID %q: %w", sudoUID, err)
		}

		gid, err := strconv.Atoi(sudoGID)
		if err != nil {
			return nil, fmt.Errorf("bad SUDO_GID %q: %w", sudoGID, err)
		}

		username := sudoUser
		homeDir := ""

		if u, err := user.LookupId(sudoUID); err == nil {
			if username == "" {
				username = u.Username
			}
			homeDir = u.HomeDir
		} else if sudoUser != "" {
			if u, err := user.Lookup(sudoUser); err == nil {
				homeDir = u.HomeDir
			}
		}

		// fallback home if user lookup fails
		if homeDir == "" {
			homeDir = os.Getenv("HOME")
		}

		return &UserContext{
			UID:      uid,
			GID:      gid,
			Username: username,
			HomeDir:  homeDir,
			IsSudo:   true,
		}, nil
	}

	// normal run without sudo
	current, err := user.Current()
	if err != nil {
		return &UserContext{
			UID:      os.Getuid(),
			GID:      os.Getgid(),
			Username: os.Getenv("USER"),
			HomeDir:  os.Getenv("HOME"),
			IsSudo:   false,
		}, nil
	}

	uid, _ := strconv.Atoi(current.Uid)
	gid, _ := strconv.Atoi(current.Gid)

	return &UserContext{
		UID:      uid,
		GID:      gid,
		Username: current.Username,
		HomeDir:  current.HomeDir,
		IsSudo:   false,
	}, nil
}

// SysProcCredential returns credentials for dropping child process privileges (tor, shell)
// to the real user with no supplementary groups.
func (u *UserContext) SysProcCredential() *syscall.Credential {
	return &syscall.Credential{
		Uid:         uint32(u.UID),
		Gid:         uint32(u.GID),
		Groups:      []uint32{},
		NoSetGroups: true,
	}
}
