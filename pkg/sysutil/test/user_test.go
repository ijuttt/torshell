package test

import (
	"os"
	"os/user"
	"strconv"
	"testing"

	"torshell/pkg/sysutil"
)

func TestIsRoot(t *testing.T) {
	expected := os.Geteuid() == 0
	if got := sysutil.IsRoot(); got != expected {
		t.Errorf("IsRoot() = %v, want %v", got, expected)
	}
}

func TestRequireRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		if err := sysutil.RequireRoot(); err != nil {
			t.Fatalf("unexpected error on root: %v", err)
		}
	} else {
		if err := sysutil.RequireRoot(); err != sysutil.ErrNotRoot {
			t.Errorf("got %v, want %v", err, sysutil.ErrNotRoot)
		}
	}
}

func TestGetRealUser_Direct(t *testing.T) {
	t.Setenv("SUDO_UID", "")
	t.Setenv("SUDO_GID", "")
	t.Setenv("SUDO_USER", "")

	ctx, err := sysutil.GetRealUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.IsSudo {
		t.Errorf("want IsSudo=false, got true")
	}
	if ctx.UID != os.Getuid() {
		t.Errorf("got UID %d, want %d", ctx.UID, os.Getuid())
	}
	if ctx.GID != os.Getgid() {
		t.Errorf("got GID %d, want %d", ctx.GID, os.Getgid())
	}
}

func TestGetRealUser_Sudo(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Skip("skip: cannot lookup current user")
	}

	t.Setenv("SUDO_UID", current.Uid)
	t.Setenv("SUDO_GID", current.Gid)
	t.Setenv("SUDO_USER", current.Username)

	ctx, err := sysutil.GetRealUser()
	if err != nil {
		t.Fatalf("GetRealUser failed in sudo mode: %v", err)
	}

	if !ctx.IsSudo {
		t.Errorf("want IsSudo=true, got false")
	}

	expectedUID, _ := strconv.Atoi(current.Uid)
	if ctx.UID != expectedUID {
		t.Errorf("got UID %d, want %d", ctx.UID, expectedUID)
	}

	expectedGID, _ := strconv.Atoi(current.Gid)
	if ctx.GID != expectedGID {
		t.Errorf("got GID %d, want %d", ctx.GID, expectedGID)
	}

	if ctx.Username != current.Username {
		t.Errorf("got username %s, want %s", ctx.Username, current.Username)
	}
}

func TestGetRealUser_BadSudoVars(t *testing.T) {
	t.Run("bad uid", func(t *testing.T) {
		t.Setenv("SUDO_UID", "not-a-number")
		t.Setenv("SUDO_GID", "1000")

		if _, err := sysutil.GetRealUser(); err == nil {
			t.Error("expected error with invalid SUDO_UID, got nil")
		}
	})

	t.Run("bad gid", func(t *testing.T) {
		t.Setenv("SUDO_UID", "1000")
		t.Setenv("SUDO_GID", "not-a-number")

		if _, err := sysutil.GetRealUser(); err == nil {
			t.Error("expected error with invalid SUDO_GID, got nil")
		}
	})
}

func TestSysProcCredential(t *testing.T) {
	ctx := &sysutil.UserContext{
		UID:      1000,
		GID:      1000,
		Username: "testuser",
		HomeDir:  "/home/testuser",
		IsSudo:   true,
	}

	cred := ctx.SysProcCredential()
	if cred.Uid != 1000 || cred.Gid != 1000 {
		t.Errorf("unexpected cred UID/GID: %d/%d", cred.Uid, cred.Gid)
	}
	if len(cred.Groups) != 0 {
		t.Errorf("want empty groups, got %v", cred.Groups)
	}
	if !cred.NoSetGroups {
		t.Error("want NoSetGroups=true")
	}
}
