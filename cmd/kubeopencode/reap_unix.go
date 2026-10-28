// Copyright Contributors to the KubeOpenCode project

//go:build unix

package main

import "syscall"

// reapOrphans collects every exited child of this process without blocking
// and returns how many it reaped.
//
// git-sync runs as PID 1 of its sidecar container, so any process orphaned
// inside the container (for example a helper git detaches during fetch) is
// re-parented to it. os/exec only waits for the children it started; nothing
// waits for these orphans, so each one stays a zombie and keeps holding a
// PID-controller slot. One leaked slot per sync, per sidecar, eventually
// exhausts the pod's pids.max or the node's threads-max.
//
// Call it only when no exec.Command started by this process is still running:
// wait4(-1) would otherwise steal that command's exit status from cmd.Wait.
func reapOrphans() int {
	n := 0
	for {
		var ws syscall.WaitStatus
		pid, err := syscall.Wait4(-1, &ws, syscall.WNOHANG, nil)
		if err == syscall.EINTR {
			continue
		}
		if pid <= 0 || err != nil {
			return n
		}
		n++
	}
}
