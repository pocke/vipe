package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	fname, err := InitTempFile()
	if err != nil {
		fail(err)
	}
	defer os.Remove(fname)

	err = Vim(fname)
	if err != nil {
		fail(err)
	}

	err = WriteResult(fname)
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func Vim(fname string) error {
	cmd := exec.Command("gvim", "--nofork", fname)
	return cmd.Run()
}

func InitTempFile() (string, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(os.TempDir(), "vipe-")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(b)
	return f.Name(), err
}

func WriteResult(fname string) error {
	res, err := os.Open(fname)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, res)
	return err
}
