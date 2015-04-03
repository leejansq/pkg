package syscallutils

import (
	"syscall"
	"unsafe"
)

func SysInfo() syscall.Sysinfo_t {
	sysinfo := syscall.Sysinfo_t{}
	syscall.RawSyscall(syscall.SYS_SYSINFO, uintptr(unsafe.Pointer(&sysinfo)), 0, 0)

	return sysinfo
}
