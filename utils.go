package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execute(cmd *exec.Cmd) error {
	fmt.Println("+", strings.Join(cmd.Args, " "))

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// stdout
	//cmd.Stdout = os.Stdout

	return cmd.Run()
}

func executeToFile(cmd *exec.Cmd, f string) error {
	fmt.Println("+", strings.Join(cmd.Args, " "))

	/////////////////
	// open the out file for writing
	outfile, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = outfile

	return cmd.Run()
	/*
		err = cmd.Start()
		if err != nil {
			return err
		}
		return cmd.Wait()
	*/
}

//cmd := exec.Command("echo", "'WHAT THE HECK IS UP'")
