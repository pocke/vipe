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
	err := Main()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	var gui bool
	flag.BoolVar(&gui, "g", false, "Use GVim")
	flag.Parse()

	fname, err := InitTempFile()
	if err != nil {
		return err
	}
	defer os.Remove(fname)

	err = Vim(fname, gui)
	if err != nil {
		return err
	}

	return WriteResult(fname)
}

func Vim(fname string, gui bool) error {
	var cmd *exec.Cmd
	if gui {
		cmd = exec.Command("gvim", "--nofork", fname)
	} else {
		// Workaround
		cmd = exec.Command("bash", "-c", fmt.Sprintf("vim %s < /dev/tty > /dev/tty", fname))
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
