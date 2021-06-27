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
	"github.com/ZenLiuCN/mpc"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io/ioutil"
	"os"
)

type Git struct {
	bit mpc.Bit8
}

func (s *Git) SetUseTemp() {
	s.bit = s.bit.Set(1)
}
func (s Git) UseTemp() bool {
	return s.bit.At(1)
}
func (s Git) Clone(url string) (repo *Repo, err error) {
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
	return &Repo{Directory: dir, repo: re}, nil
}
func (s Git) CloneAuth(url string, auth transport.AuthMethod) (repo *Repo, err error) {
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
	return &Repo{Directory: dir, repo: re}, nil
}

type Signature struct {
	Name string
	Hash string
}
type Repo struct {
	Directory string
	repo      *git.Repository
	branch    []Signature
}

func (s *Repo) Branches() ([]Signature, error) {
	if s.branch == nil {
		x, er := s.repo.Branches()
		if er != nil {
			return nil, er
		}
		s.branch = make([]Signature, 0, 5)
		_ = x.ForEach(func(ref *plumbing.Reference) error {
			s.branch = append(s.branch, Signature{Name: ref.Name().String(), Hash: ref.Hash().String()})
			return nil
		})
	}
	return s.branch, nil
}
