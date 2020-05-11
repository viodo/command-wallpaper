package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var fpath string

func main() {
	fpath = os.Args[1]
	if runtime.GOOS == "windows" {
		SetWindows()
	}
	if runtime.GOOS == "linux" {
		SetLinux()
	}
	if runtime.GOOS == "darwin" {
		SetDarwin()
	}
	log.Fatalln("Unknown platform")
}

// SetWindows set windows wallpaper
func SetWindows() {
	dll := syscall.MustLoadDLL("user32.dll")
	proc := dll.MustFindProc("SystemParametersInfoW")
	defer syscall.FreeLibrary(dll.Handle)
	uiAction := 0x0014
	uiParam := 0
	pvParam := syscall.StringToUTF16Ptr(fpath)
	fWinIni := 1
	r2, _, err := proc.Call(uintptr(uiAction),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(pvParam)),
		uintptr(fWinIni))
	if r2 != 0 {
		log.Fatalln(r2, err, fpath)
	}
}

// SetLinux set linux wallpaper
func SetLinux() {
	shell := fmt.Sprintf("/usr/bin/xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitor%d/workspace0/last-image -s %s/", 0, fpath)
	arr := strings.Split(shell, " ")
	command := exec.Command(arr[0], arr[1:len(arr)]...)
	out, err := command.Output()
	if err != nil {
		log.Fatalln(out, err, fpath)
	}
}

// SetDarwin set darwin wallpaper
func SetDarwin() {
	shell := fmt.Sprintf(`osascript -e "tell application \"Finder\" to set desktop picture to POSIX file \"%s\""`, fpath)
	arr := strings.Split(shell, " ")
	command := exec.Command(arr[0], arr[1:len(arr)]...)
	out, err := command.Output()
	if err != nil {
		log.Fatalln(out, err, fpath)
	}
}
