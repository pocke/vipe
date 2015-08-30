package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	var gui bool
	flag.BoolVar(&gui, "g", false, "Use GVim")
	flag.Parse()

	fname, err := InitTempFile()
	if err != nil {
		fail(err)
	}
	defer os.Remove(fname)

	err = Vim(fname, gui)
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

func Vim(fname string, gui bool) error {
	var cmd *exec.Cmd
	if gui {
		cmd = exec.Command("gvim", "--nofork", fname)
	} else {
		in, err := os.Open("/dev/tty") // XXX: What's best way...?
		if err != nil {
			return err
		}
		defer in.Close()

		cmd = exec.Command("vim", fname)
		cmd.Stdin = in
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
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
