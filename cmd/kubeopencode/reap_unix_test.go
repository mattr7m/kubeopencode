// Copyright Contributors to the KubeOpenCode project

//go:build unix

package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestReapOrphansCollectsExitedChildren(t *testing.T) {
	// Start children without waiting on them: they become zombies of this
	// process, the same state an orphan reaches under a PID-1 git-sync.
	const children = 3
	for i := 0; i < children; i++ {
		if err := exec.Command("true").Start(); err != nil {
			t.Fatalf("start child: %v", err)
		}
	}

	reaped := 0
	deadline := time.Now().Add(5 * time.Second)
	for reaped < children && time.Now().Before(deadline) {
		reaped += reapOrphans()
		time.Sleep(10 * time.Millisecond)
	}
	if reaped != children {
		t.Fatalf("reaped %d children, want %d", reaped, children)
	}
	if n := reapOrphans(); n != 0 {
		t.Fatalf("second pass reaped %d, want 0", n)
	}
}
