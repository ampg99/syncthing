package fswatcher

import (
	"errors"
	"syscall"
	"testing"
)

func TestErrorInotifyInterpretation(t *testing.T) {
	msg := "Failed to install inotify handler for test-folder." +
		" Please increase inotify limits," +
		" see http://bit.ly/1PxkdUC for more information."
	var errTooManyFiles syscall.Errno = 24
	var errNoSpace syscall.Errno = 28
	err := interpretNotifyWatchError(errTooManyFiles, "test-folder")
	if err.Error() != msg {
		t.Errorf("Expected error about inotify limits, but got: %#v",
			err)
	}
	err = interpretNotifyWatchError(errNoSpace, "test-folder")
	if err.Error() != msg {
		t.Errorf("Expected error about inotify limits, but got: %#v",
			err)
	}
	err = interpretNotifyWatchError(
		errors.New("Another error"), "test-folder")
	if err.Error() != "Another error" {
		t.Errorf("Unexpected error: %#v", err)
	}
}
