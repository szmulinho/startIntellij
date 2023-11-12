package main

import (
	"fmt"
	"github.com/moutend/go-hook"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

const (
	VK_I = 0x49 // Virtual key code for 'I'
)

var (
	moduser32            *syscall.DLL
	procGetAsyncKeyState *syscall.Proc
)

func init() {
	moduser32 = syscall.MustLoadDLL("user32.dll")
	procGetAsyncKeyState = moduser32.MustFindProc("GetAsyncKeyState")
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("This example is intended for Windows only.")
		os.Exit(1)
	}

	// Start listening for the 'i' key
	err := hook.Register(hook.KeyDown, []string{"I"}, onKeyDown)
	if err != nil {
		log.Fatal(err)
	}
	defer hook.Unregister()

	// Run the program indefinitely
	select {}
}

func onKeyDown(e hook.Event) {
	if e.VkCode == VK_I && isKeyPressed(VK_I) {
		fmt.Println("Launching IntelliJ...")
		launchIntelliJ()
	}
}

func isKeyPressed(vkCode int) bool {
	var state uint16
	procGetAsyncKeyState.Call(uintptr(vkCode), uintptr(unsafe.Pointer(&state)))
	return state&0x8000 != 0
}

func launchIntelliJ() {
	// Adjust the command to the actual command to run IntelliJ
	cmd := exec.Command("cmd", "/c", "start", "intellij")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
