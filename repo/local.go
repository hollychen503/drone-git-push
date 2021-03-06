package repo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

// LocalReplaceTag ...
func LocalReplaceTag(remote, localbranch string, branch string, force bool, followtags bool, delRemoteTag bool) *exec.Cmd {

	//#Replace the tag to reference the most recent commit
	//git tag -fa ${VERSION} -m "tag to $VERSION"

	// find version
	// read tag from .tag file
	b, err := ioutil.ReadFile("version") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return nil
	}
	//fmt.Println(b) // print the content as 'bytes'
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'

	re2 := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	str = re2.FindString(str)

	//fmt.Printf("%s\n")
	fa := "-a"
	if force {
		fa = "-fa"
	}

	cmd := exec.Command(
		"git",
		"tag",
		fa,
		str,
		"-m",
		"\"tag to "+str+" by drone CI\"")

	return cmd
}
