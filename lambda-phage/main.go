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
		fmt.Printf("Error building go executable, %s", err.Error())
		return
	}

	zFile, err := newZipFile(binName + ".zip")
	if err != nil {
		zipFileFail(err)
		return
	}

	_, err = zFile.AddFile("/tmp/" + binName)

	if err != nil {
		zipFileFail(err)
		return
	}

	_, err = zFile.AddString("index.js", jsloader)

	if err != nil {
		zipFileFail(err)
		return
	}

	err = zFile.Close()

}

type zipFile struct {
	f *os.File
	*zip.Writer
}

// creates a new zip file
func newZipFile(fName string) (*zipFile, error) {
	f, err := os.Create(fName)
	if err != nil {
		return nil, err
	}

	return &zipFile{
		f,
		zip.NewWriter(f),
	}, nil
}

// adds a file to this archive
func (z *zipFile) AddFile(fName string) (int64, error) {
	debug := debug.Debug("zipFile.AddFile")
	debug("opening source file")
	f, err := os.Open(fName)
	if err != nil {
		return 0, err
	}

	debug("source file opened")

	// want to make sure we get the base name for a file,
	// so fstat it, which is able to do that
	debug("fstat source file")
	s, err := f.Stat()
	if err != nil {
		return 0, err
	}
	bName := s.Name()

	// add file to archive
	debug("adding file to archive")
	wr, err := z.Create(bName)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(wr, f)
	if err != nil {
		return 0, err
	}

	err = f.Close()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// adds a file from string data to the archive
func (z *zipFile) AddString(fName string, str []byte) (int64, error) {
	debug := debug.Debug("zipFile.AddString")

	buf := bytes.NewBuffer(str)

	// add file to archive
	debug("adding file to archive")

	wr, err := z.Create(fName)
	if err != nil {
		return 0, err
	}

	debug("copying data for file")
	n, err := io.Copy(wr, buf)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// closes the writer and the file
func (z *zipFile) Close() error {
	err := z.Writer.Close()
	if err != nil {
		return err
	}

	return z.f.Close()
}

func zipFileFail(err error) {
	_, f, l, _ := runtime.Caller(1)
	fmt.Printf("[%s:%s]error creating zip file, %s\n", f, l, err.Error())
	return
}
