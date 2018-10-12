package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/hollychen503/drone-git-push/repo"
)

type (
	// Netrc structure
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	// Commit structure
	Commit struct {
		Author Author
	}

	// Author structure
	Author struct {
		Name  string
		Email string
	}

	// Config structure
	Config struct {
		Key           string
		Remote        string
		RemoteName    string
		Branch        string
		LocalBranch   string
		Path          string
		Force         bool
		FollowTags    bool
		TagRemote     bool
		SkipVerify    bool
		Commit        bool
		CommitMessage string
		EmptyCommit   bool
	}

	// Plugin Structure
	Plugin struct {
		Netrc  Netrc
		Commit Commit
		Config Config
	}
)

// Exec starts the plugin execution.
func (p Plugin) Exec() error {
	if err := p.HandlePath(); err != nil {
		return err
	}

	if err := p.WriteConfig(); err != nil {
		return err
	}

	if err := p.WriteKey(); err != nil {
		return err
	}

	if err := p.WriteNetrc(); err != nil {
		return err
	}

	if err := p.HandleCommit(); err != nil {
		return err
	}

	if err := p.HandleRemote(); err != nil {
		return err
	}

	//
	if p.Config.TagRemote {
		if err := p.HandlePushTagRemote(); err != nil {
			return err
		}
	} else {
		if err := p.HandlePush(); err != nil {
			return err
		}
	}

	return p.HandleCleanup()
}

// WriteConfig writes all required configurations.
func (p Plugin) WriteConfig() error {
	if err := repo.GlobalName(p.Commit.Author.Name).Run(); err != nil {
		return err
	}

	if err := repo.GlobalUser(p.Commit.Author.Email).Run(); err != nil {
		return err
	}

	if p.Config.SkipVerify {
		if err := repo.SkipVerify().Run(); err != nil {
			return err
		}
	}

	return nil
}

// WriteKey writes the private SSH key.
func (p Plugin) WriteKey() error {
	return repo.WriteKey(
		p.Config.Key,
	)
}

// WriteNetrc writes the netrc config.
func (p Plugin) WriteNetrc() error {
	return repo.WriteNetrc(
		p.Netrc.Machine,
		p.Netrc.Login,
		p.Netrc.Password,
	)
}

// HandleRemote adds the git remote if required.
func (p Plugin) HandleRemote() error {
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteAdd(p.Config.RemoteName, p.Config.Remote)); err != nil {
			return err
		}
	}

	return nil
}

// HandlePath changes to a different directory if required
func (p Plugin) HandlePath() error {
	if p.Config.Path != "" {
		if err := os.Chdir(p.Config.Path); err != nil {
			return err
		}
	}

	return nil
}

// HandleCommit commits dirty changes if required.
func (p Plugin) HandleCommit() error {
	if p.Config.Commit {
		if err := execute(repo.Add()); err != nil {
			return err
		}

		if err := execute(repo.TestCleanTree()); err != nil {
			// changes to commit
			if err := execute(repo.ForceCommit(p.Config.CommitMessage)); err != nil {
				return err
			}
		} else { // no changes
			if p.Config.EmptyCommit {
				// no changes but commit anyway
				if err := execute(repo.EmptyCommit(p.Config.CommitMessage)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// HandlePush pushs the changes to the remote repo.
func (p Plugin) HandlePush() error {
	var (
		name       = p.Config.RemoteName
		local      = p.Config.LocalBranch
		branch     = p.Config.Branch
		force      = p.Config.Force
		followtags = p.Config.FollowTags
	)

	return execute(repo.RemotePushNamedBranch(name, local, branch, force, followtags))
}

func tagExist(remote, local string) (bool, error) {
	// open local
	// find version
	// read tag from .tag file
	b, err := ioutil.ReadFile(local) // just pass the file name
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	//fmt.Println(b) // print the content as 'bytes'
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'

	re2 := regexp.MustCompile(`v\d+\.\d+\.\d+`) // v1.2.3
	str = re2.FindString(str)
	if len(str) < 6 {
		return false, errors.New("no version in file" + local)
	}

	// get data from remote tag info
	rb, err := ioutil.ReadFile(remote) // just pass the file name
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	//fmt.Println(b) // print the content as 'bytes'
	strRemote := string(rb) // convert content to a 'string'
	//fmt.Println(str) // print the content as a 'string'

	// find new version in remote tag info
	re := regexp.MustCompile("(?m)/" + str + `\b`)
	m := re.FindAllString(strRemote, -1)
	fmt.Printf("%d matches\n", len(m))
	if len(m) > 0 {
		for _, s := range m {
			fmt.Println(s)
		}
		return true, nil
	}

	return false, nil
}

// HandlePushTagRemote delete tag to the remote repo.
func (p Plugin) HandlePushTagRemote() error {
	var (
		name         = p.Config.RemoteName
		local        = p.Config.LocalBranch
		branch       = p.Config.Branch
		force        = p.Config.Force
		followtags   = p.Config.FollowTags
		delRemoteTag = p.Config.TagRemote
	)

	err := errors.New("")
	/*
		t := time.Now()
		fmt.Println(t.Format(time.RFC3339))

		err := executeToFile(repo.RemoteGetTags(name, local, branch, force, followtags, delRemoteTag), ".remoteTags")
		//err := execute(repo.RemoteGetTags(name, local, branch, force, followtags, delRemoteTag))
		if err != nil {
			return err
		}

		t = time.Now()
		fmt.Println(t.Format(time.RFC3339))
	*/

	if force {
		err = execute(repo.RemoteDeleteTag(name, local, branch, force, followtags, delRemoteTag))
		if err != nil {
			return err
		}
	}
	/*
		else {
			isIn, err := tagExist(".remoteTags", ".tags")
			if err != nil {
				return err
			}
			if isIn {
				return errors.New("tag is existed in repo")
			}
		}*/

	err = execute(repo.LocalReplaceTag(name, local, branch, force, followtags, delRemoteTag))
	if err != nil {
		return err
	}
	err = execute(repo.RemoteAddTag(name, local, branch, force, followtags, delRemoteTag))
	if err != nil {
		return err
	}

	return err
}

// HandleCleanup does eventually do some cleanup.
func (p Plugin) HandleCleanup() error {
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteRemove(p.Config.RemoteName)); err != nil {
			return err
		}
	}

	return nil
}
