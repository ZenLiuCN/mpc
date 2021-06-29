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
	"archive/zip"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
	ssh2 "golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//region Git
type Git struct {
}

func (s *Git) Clone(url string, repoPath string, auth transport.AuthMethod) (repo *Repo, err error) {
	var dir string
	if repoPath == "" {
		dir, err = ioutil.TempDir("", "repo_*")
		if err != nil {
			return nil, err
		}
	} else {
		dir = repoPath
	}

	co := new(git.CloneOptions)
	co.URL = url
	co.RecurseSubmodules = git.DefaultSubmoduleRecursionDepth
	co.Progress = os.Stdout
	if auth != nil {
		co.Auth = auth
	}
	rootFs := osfs.New(dir)
	gitFs := osfs.New(path.Join(dir, ".git"))
	store := filesystem.NewStorage(gitFs, cache.NewObjectLRUDefault())
	var re *git.Repository
	re, err = git.Clone(store, rootFs, co)
	//re, err = git.PlainClone(dir, false, co)
	return &Repo{Raw: re, Path: dir}, err
}
func (s *Git) Open(path string) (repo *Repo, err error) {
	var re *git.Repository
	re, err = git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repo{
		Raw:  re,
		Path: path,
	}, nil
}

//endregion

type Signature struct {
	Name string
	Hash string
}

//region Repo
type Repo struct {
	Raw           *git.Repository
	Path          string
	br            []*Branch
	currentBranch *Branch
}

func (r *Repo) Pull(auth transport.AuthMethod) (err error) {
	w, err := r.Raw.Worktree()
	if err != nil {
		return err
	}
	opt := new(git.PullOptions)
	opt.Auth = auth
	opt.Progress = os.Stdout
	err = w.Pull(opt)
	return
}
func (r *Repo) CurrentBranch() (*Branch, error) {
	if r.currentBranch == nil {
		if v, err := r.Raw.Head(); err != nil {
			return nil, err
		} else {
			r.currentBranch = newBranch(v)
		}
	}
	return r.currentBranch, nil
}
func (r *Repo) Branches() ([]*Branch, error) {
	if r.br == nil {
		x, err := r.Raw.References()
		if err != nil {
			return nil, err
		}
		r.br = make([]*Branch, 0, 5)
		err = x.ForEach(func(ref *plumbing.Reference) error {
			bx := newBranch(ref)
			if bx != nil {
				r.br = append(r.br, bx)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return r.br, nil
}
func (r *Repo) Remotes() ([]*git.Remote, error) {
	return r.Raw.Remotes()
}
func (r *Repo) Tags() ([]*Signature, error) {
	t, er := r.Raw.Tags()
	if er != nil {
		return nil, er
	}
	tag := make([]*Signature, 0, 5)
	_ = t.ForEach(func(ref *plumbing.Reference) error {
		tag = append(tag, &Signature{
			ref.Name().String(),
			ref.Hash().String(),
		})
		return nil
	})
	return tag, er
}
func (r *Repo) Checkout(branch *Branch) (err error) {
	if branch.Remote {
		if _, err = r.Raw.Branch(branch.ShortName()); err != nil && err == git.ErrBranchNotFound {
			local := plumbing.NewBranchReferenceName(branch.ShortName())
			remote := plumbing.NewRemoteReferenceName(branch.Origin(), branch.ShortName())
			err = r.Raw.CreateBranch(&config.Branch{Name: branch.ShortName(), Remote: branch.Origin(), Merge: local})
			if err != nil {
				return err
			}
			newReference := plumbing.NewSymbolicReference(local, remote)
			if err = r.Raw.Storer.SetReference(newReference); err != nil {
				return
			}
		}

	}
	var w *git.Worktree
	if w, err = r.Raw.Worktree(); err != nil {
		return
	}
	opt := new(git.CheckoutOptions)
	opt.Branch = plumbing.NewBranchReferenceName(branch.ShortName())
	r.currentBranch = branch
	err = w.Checkout(opt)
	return
}

//zip as standard zip @see ZipWith
func (r *Repo) Zip() (path *os.File, err error) {
	path, err = ioutil.TempFile("", "temp_zip_*")
	w := zip.NewWriter(path)
	defer func() {
		err1 := w.Close()
		if err1 != nil {
			err = err1
		}
		return
	}()
	err = AddFilesToZip(w, r.Path, "")
	if err != nil {
		return nil, err
	}
	return path, nil
}

//create a self defined zip structure @see AddFilesToZip as example
func (r *Repo) ZipWith(zipperCreator func() (w *zip.Writer, biz string), addFiles func(w *zip.Writer, basePath, baseInZip string) error) (path *os.File, err error) {
	w, biz := zipperCreator()
	defer func() {
		err1 := w.Close()
		if err1 != nil {
			err = err1
		}
		return
	}()
	err = addFiles(w, r.Path, biz)
	if err != nil {
		return nil, err
	}
	return path, nil
}

//read go.mod in current repository under subPath
func (r *Repo) Mod(subPath string) (string, error) {
	m := path.Join(r.Path, subPath, "go.mod")
	b, err := ioutil.ReadFile(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//a standard Add Files to zip
func AddFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(path.Join(basePath, file.Name()))
			if err != nil {
				return err
			}
			f, err := w.Create(path.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			if file.Name() == ".git" {
				continue
			}
			// Recurse
			err = AddFilesToZip(w, path.Join(basePath, file.Name()), path.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//endregion

//region SSH AUTH
func SSHAuthOfFile(file string) (transport.AuthMethod, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", file, "")
	if err != nil {
		return nil, err
	}
	auth.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	return auth, nil
}
func SSHAuthOfPEM(pem string) (transport.AuthMethod, error) {
	auth, err := ssh.NewPublicKeys("git", []byte(pem), "")
	if err != nil {
		return nil, err
	}
	auth.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	return auth, nil
}

//endregion

//region Branch
type Branch struct {
	Remote bool
	Signature
	short  string
	origin string
}

func newBranch(b *plumbing.Reference) *Branch {
	if b.Name().IsTag() || b.Name().IsNote() || b.Hash().IsZero() {
		return nil
	}
	return &Branch{Remote: b.Name().IsRemote(), Signature: Signature{
		Name: b.Name().Short(),
		Hash: b.Hash().String(),
	}}
}
func (s *Branch) ShortName() string {
	if s.short == "" {
		if s.Remote {
			x := strings.SplitN(s.Name, "/", 2)
			s.origin = x[0]
			s.short = x[1]
		} else {
			s.short = s.Name
		}
	}
	return s.short
}
func (s *Branch) Origin() string {
	if s.short == "" {
		if s.Remote {
			x := strings.SplitN(s.Name, "/", 2)
			s.origin = x[0]
			s.short = x[1]
		} else {
			s.short = s.Name
		}
	}
	return s.short
}

//endregion
