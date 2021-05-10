package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCheckBinFolderToInstall(t *testing.T) {
	installpath := os.Getenv("GOPATH") + "\\bin"
	fmt.Println(installpath)

	powershellpath, _ := exec.LookPath("powershell.exe")
	fmt.Println(powershellpath)

	fmt.Println(os.Getenv("HOME") + ".lc\\leetcode\\user.json")
}
