package repo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

// RemoteRemove drops the defined remote from a git repo.
func RemoteRemove(name string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"remote",
		"rm",
		name)

	return cmd
}

// RemoteAdd adds an additional remote to a git repo.
func RemoteAdd(name, url string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"remote",
		"add",
		name,
		url)

	return cmd
}

// RemotePush pushs the changes from the local head to a remote branch..
func RemotePush(remote, branch string, force bool, followtags bool) *exec.Cmd {
	return RemotePushNamedBranch(remote, "HEAD", branch, force, followtags)
}

// RemotePushNamedBranch puchs changes from a local to a remote branch.
func RemotePushNamedBranch(remote, localbranch string, branch string, force bool, followtags bool) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"push",
		remote,
		localbranch+":"+branch)

	if force {
		cmd.Args = append(
			cmd.Args,
			"--force")
	}

	if followtags {
		cmd.Args = append(
			cmd.Args,
			"--follow-tags")
	}

	return cmd
}

// RemoteGetTags ...
func RemoteGetTags(remote, localbranch string, branch string, force bool, followtags bool, delRemoteTag bool) *exec.Cmd {
	fmt.Println("local branch is ", localbranch)
	// git ls-remote --tags origin  // now work.
	// now i clone all tags from remote by plugin/git with tags=true

	// git show-ref --tags
	cmd := exec.Command(
		//"/bin/sh", //alpine
		//"-c",
		//`"/usr/bin/git tag > .remoteTags"`,
		"git",
		"tag",
	//">",
	//".remoteTags",
	//"ls-remote",
	//"--tags",
	//"show-ref",
	//"--tags",
	//"origin",
	//" | tee .remoteTags",
	)

	return cmd

}

// RemoteDeleteTag ...
func RemoteDeleteTag(remote, localbranch string, branch string, force bool, followtags bool, delRemoteTag bool) *exec.Cmd {
	fmt.Println("local branch is ", localbranch)

	// find version
	// read tag from .tag file
	b, err := ioutil.ReadFile(".tags") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return nil
	}
	//fmt.Println(b) // print the content as 'bytes'
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'

	re2 := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	str = re2.FindString(str)

	cmd := exec.Command(
		"git",
		"push",
		remote,
		":refs/tags/"+str)

	return cmd
}

// RemoteAddTag ...
func RemoteAddTag(remote, localbranch string, branch string, force bool, followtags bool, delRemoteTag bool) *exec.Cmd {

	// find version
	// read tag from .tag file
	b, err := ioutil.ReadFile(".tags") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return nil
	}
	//fmt.Println(b) // print the content as 'bytes'
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'

	re2 := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	str = re2.FindString(str)

	// #Push the tag to the remote origin
	// git push --tags

	cmd := exec.Command(
		"git",
		"push",
		"--tags")

	return cmd
}
