/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	leetcodeExe = "leetcode-cli.exe"
)

var (
	gopath          string
	leetcodeExePath string
	userJsonPath    string
)

func fileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func userJsonExist() bool {
	return true
}

func installLeetCodeCLI() error {
	installpath := os.Getenv("GOPATH") + "\\bin"
	log.Println(installpath)

	// 0. ready powershell
	posh, err := exec.LookPath("powershell.exe")
	if err != nil {
		return err
	}

	var (
		args []string
		cmd  *exec.Cmd
	)
	// 1. download and extract .zip on bin folder
	if !fileExist("bin/leetcode-cli.zip") {
		args = []string{"-c", "wget", "-outfile", "bin/leetcode-cli.zip", "-uri", "https://github.com/skygragon/leetcode-cli/releases/download/2.6.2/leetcode-cli.node10.win32.x64.zip"}
		cmd = exec.Command(posh, args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	args = []string{"-c", "expand-archive", "-path", "bin/leetcode-cli.zip", "-destinationpath", "bin"}
	cmd = exec.Command(posh, args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	// 2. copy all to the gopath
	files := []string{
		"leetcode-cli.exe",
		"binding.node",
		"ffi_bindings.node",
		"node_sqlite3.node",
	}

	for _, file := range files {
		if err := os.Rename("bin\\dist\\"+file, installpath+"\\"+file); err != nil {
			fmt.Println(err)
			// remove all the renamed of gopath
			return err
		}
	}
	return nil
}

// Not working yet
func getLeetcodeCookies() (string, string, string) {

	// first leetcode-cli.exe
	leetcodecli, err := exec.LookPath("leetcode-cli.exe")
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	fmt.Println(leetcodecli)

	cmd := exec.Command(leetcodecli, "user", "-L")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return "", "", ""
	}

	// second leetcode-cli.exe

	var username, leetcode_session, csrftoken string
	url := "https://leetcode.com/profile/"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Cookies())
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "LEETCODE_SESSION" {
			log.Println("Cookie LEETCODE_SESSION: ", cookie.Value)
		}
		if cookie.Name == "csrftoken" {
			log.Println("Cookie csrftoken: ", cookie.Value)
		}
		if cookie.Name == "username" {
			log.Println("Cookie username: ", cookie.Value)
		}
		fmt.Println("Found a cookie named:", cookie.Name)
	}
	return username, leetcode_session, csrftoken
}

/*
	try:
			userid, leetcode_session, crsftoken = get_leetcode_cookies()
	except ValueError as e:
			print(e.args)
	else:
			with open(os.path.join(home_folder, ".lc", "leetcode", "user.json"), "w") as f:
					f.write("{\n")
					f.write(f'    "login": "{userid}",\n')
					f.write('    "loginCSRF": "",\n')
					f.write(f'    "sessionCSRF": "{crsftoken}",\n')
					f.write(f'    "sessionId": "{leetcode_session}"\n')
					f.write("}")
			os.system(os.path.join("bin", "dist", "leetcode-cli") + " user -c")
			print(f"Logged in as {userid}")
*/

func makeUserJson() error {
	userid, leetcode_session, crsftoken := getLeetcodeCookies()

	_ = userid
	_ = leetcode_session
	_ = crsftoken
	return nil
}

func prepareLeetCodeCLI() error {
	if !fileExist(leetcodeExePath) {
		if err := installLeetCodeCLI(); err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("leetcode-cli.exe installed to GOPATH")
	}
	if !fileExist(userJsonPath) {
		if err := makeUserJson(); err != nil {
			return err
		}
	}
	fmt.Println("file exist: ", userJsonPath)
	return nil

}

func getProblems(cmd *cobra.Command, problems []string) {
	prepareLeetCodeCLI()

	fmt.Printf("getProblems of %v\n", problems)

	// log in: 	leetcode-cli.exe user -c
	leetcodecli, err := exec.LookPath("leetcode-cli.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	execCmd := exec.Command(leetcodecli, "user", "-c")
	if err := execCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	// get and save
	// leetcode-cli.exe show <num> -gx -l golang

	for _, pNum := range problems {
		// make dir
		os.Mkdir(pNum, 0600)
		os.Chdir(pNum)
		// download go
		cmd := exec.Command(leetcodecli, "show", pNum, "-gx", "-l", "golang")
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
		os.Chdir("..")

		// rename == move

		// do the job

	}
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a problem of the number from leetcode.com",
	Long: `Get a problem of the number from leetcode.com and
1) make folder of the problem number
2) make .go file with the main function which can test the code`,
	Run: getProblems,
}

func init() {
	rootCmd.AddCommand(getCmd)

	gopath = os.Getenv("GOPATH")
	leetcodeExePath = gopath + "\\bin\\" + leetcodeExe
	userJsonPath = os.Getenv("HOME") + "\\.lc\\leetcode\\user.json"

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
