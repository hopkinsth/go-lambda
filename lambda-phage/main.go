package main

import "fmt"
import "github.com/spf13/cobra"
import "archive/zip"
import "os"
import "io"
import "os/exec"
import "github.com/lucsky/cuid"
import "runtime"
import "github.com/tj/go-debug"
import "bytes"

func main() {
	pkgCmd := &cobra.Command{
		Use:   "pkg [import-path]",
		Short: "packages your lambda function up into a zip file",
		Long:  "wtf?",
		Run:   pkg,
	}

	var root = &cobra.Command{Use: "lambda-phage"}
	root.AddCommand(pkgCmd)
	root.Execute()
}

// packages your package up into a zip file
func pkg(c *cobra.Command, _ []string) {
	var err error
	// cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("lambda-phage can't find your package!")
		return
	}

	binName := "lambda-phage-" + cuid.New()

	cmd := exec.Command("go", "build", "-o", "/tmp/"+binName)
	cmd.Stderr = os.Stderr

	// copy the environment from parent proc
	// and add flags for
	pEnv := os.Environ()
	env := make([]string, len(pEnv)+2)
	copy(env, pEnv)
	env[len(env)-2] = "GOOS=linux"
	env[len(env)-1] = "GOARCH=amd64"
	cmd.Env = env

	err = cmd.Run()
	if err != nil {
		fmt.Printf("error building go program, %s", err.Error())
		return
	}

	zFile, err := os.Create(binName + ".zip")
	if err != nil {
		fmt.Printf("error creating zip file, %s", err.Error())
		return
	}

	zWr := zip.NewWriter(zFile)

	binZip, err := zWr.Create(binName)
	if err != nil {
		fmt.Printf("error adding to zip file, %s", err.Error())
		return
	}

	bFile, err := os.Open("/tmp/" + binName)
	if err != nil {
		fmt.Printf("error opening compiled binary, %s", err.Error())
		return
	}

	_, err = io.Copy(binZip, bFile)
	if err != nil {
		fmt.Printf("error adding to zip file, %s", err.Error())
		return
	}

	err = zWr.Close()
	if err != nil {
		fmt.Printf("error creating zip file, %s", err.Error())
		return
	}

}
