/*
 *
 * Copyright (C) 2021.  Zen.Liu
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io/ioutil"
	"os"
)

type Git int

func (s Git) UseTemp() bool {
	return s&(1<<1) == 0
}
func (s Git) Clone(url string) (repo *GitRepo, err error) {
	var re *git.Repository
	var dir string
	dir, err = ioutil.TempDir("", "git_clone*")
	if err != nil {
		return
	}
	re, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return
	}
	return &GitRepo{dir, re}, nil
}
func (s Git) CloneAuth(url string, auth transport.AuthMethod) (repo *GitRepo, err error) {
	var re *git.Repository
	var dir string
	dir, err = ioutil.TempDir("", "git_clone*")
	if err != nil {
		return
	}
	re, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: os.Stdout,
	})
	return &GitRepo{dir, re}, nil
}

type GitRepo struct {
	Directory string
	*git.Repository
}

func (s GitRepo) Branches() {

}
