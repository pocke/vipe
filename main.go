package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fail(err)
	}

	f, err := ioutil.TempFile(os.TempDir(), "vipe-")
	if err != nil {
		fail(err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	_, err = f.Write(b)
	if err != nil {
		fail(err)
	}
	f.Close()

	err = Vim(f.Name())
	if err != nil {
		fail(err)
	}

	res, err := os.Open(f.Name())
	if err != nil {
		fail(err)
	}
	_, err = io.Copy(os.Stdout, res)
	if err != nil {
		fail(err)
	}
}

func Vim(fname string) error {
	cmd := exec.Command("gvim", "--nofork", fname)
	return cmd.Run()
}
