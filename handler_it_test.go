// +build integration

package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

func TestE2EProxy(t *testing.T) {
	is := is.New(t)
	s := httptest.NewServer(newServer())

	repoDir, err := ioutil.TempDir("", "git-cache-it_******")
	is.NoErr(err)
	defer os.RemoveAll(repoDir)

	err = exec.Command(
		"git",
		"-c", "http.extraHeader='GIT-CACHE-UPSTREAM: https://github.com'",
		"-c", "protocol.version=2",
		"clone",
		"--single-branch",
		"--no-tags",
		"--no-checkout",
		s.URL+"/git/git.git",
		repoDir,
	).Run()
	is.NoErr(err)
}
