/*
Sniperkit-Bot
- Date: 2018-08-12 12:10:57.421562389 +0200 CEST m=+0.041271547
- Status: analyzed
*/

// +build windows

package osfs

import (
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

type fileInfo struct {
	os.FileInfo
	name string
}

func (fi *fileInfo) Name() string {
	return fi.name
}

var (
	kernel32DLL    = windows.NewLazySystemDLL("kernel32.dll")
	lockFileExProc = kernel32DLL.NewProc("LockFileEx")
	unlockFileProc = kernel32DLL.NewProc("UnlockFile")
)

const (
	lockfileExclusiveLock = 0x2
)

func (f *file) Lock() error {
	f.m.Lock()
	defer f.m.Unlock()

	var overlapped windows.Overlapped
	// err is always non-nil as per sys/windows semantics.
	ret, _, err := lockFileExProc.Call(f.File.Fd(), lockfileExclusiveLock, 0, 0xFFFFFFFF, 0,
		uintptr(unsafe.Pointer(&overlapped)))
	runtime.KeepAlive(&overlapped)
	if ret == 0 {
		return err
	}
	return nil
}

func (f *file) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()

	// err is always non-nil as per sys/windows semantics.
	ret, _, err := unlockFileProc.Call(f.File.Fd(), 0, 0, 0xFFFFFFFF, 0)
	if ret == 0 {
		return err
	}
	return nil
}
