package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kyle-aoki/uu"
)

func main() {
	defer uu.MainRecover()

	args := os.Args[1:]
	if len(args) < 1 {
		panic("incorrect number of arguments")
	}

	log.SetOutput(ioutil.Discard)
	if len(args) == 2 {
		if args[1] == "-v" {
			log.Println("verbose mode")
			log.SetOutput(os.Stdout)
		}
	}
	
	pwd, err := os.Getwd()
	uu.MustExec(err)
	log.Println("pwd:", pwd)

	baseDir := filepath.Base(pwd)
	log.Println("base dir:", baseDir)

	pubDir := join("tmp", "pub")
	log.Println("pub dir:", pubDir)

	workbenchDir := join(pubDir, baseDir)
	log.Println("workbench dir:", workbenchDir)

	CMD("rm", "-rf", pubDir)
	CMD("mkdir", pubDir)

	CMD("cp", "-r", pwd, workbenchDir)
	CMD("rm", "-rf", join(workbenchDir, "main.go"))

	fileName := args[0]                                    // program.go
	programName := strings.Replace(fileName, ".go", "", 1) // program

	contents, err := ioutil.ReadFile(join(workbenchDir, fileName))
	uu.MustExec(err)

	fileFunc := fmt.Sprintf("func %s()", programName)

	updatedContents := strings.Replace(string(contents), fileFunc, "func main()", 1)

	err = ioutil.WriteFile(join(workbenchDir, "main.go"), []byte(updatedContents), 0777)
	uu.MustExec(err)

	CMD("rm", "-rf", join(workbenchDir, fileName)) // delete program.go

	// build
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = workbenchDir
	err = cmd.Run()
	uu.MustExec(err)

	home, err := os.UserHomeDir()
	uu.MustExec(err)

	mod := getModuleName(workbenchDir)
	CMD("mv", join(workbenchDir, mod), join(home, "bin", programName))
	CMD("rm", "-rf", "/tmp/pub")
}

func CMD(name string, arg ...string) {
	log.Println(name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	err := cmd.Run()
	uu.MustExec(err)
}

func join(elem ...string) string {
	return filepath.Clean("/" + filepath.Join(elem...))
}

func getModuleName(pubDir string) string {
	moduleFile := join(pubDir, "go.mod")
	bs, err := ioutil.ReadFile(moduleFile)
	uu.MustExec(err)
	s := string(bs)
	lines := strings.Split(s, "\n")
	for i := range lines {
		if strings.Contains(lines[i], "module ") {
			trimmed := strings.Trim(lines[i], " \n\t\r")
			moduleName := strings.Replace(trimmed, "module ", "", -1)
			return moduleName
		}
	}
	panic("invalid go.mod file")
}
