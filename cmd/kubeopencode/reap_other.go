// Copyright Contributors to the KubeOpenCode project

//go:build !unix

package main

// reapOrphans is a no-op where wait4 is unavailable.
func reapOrphans() int { return 0 }
