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
	"fmt"
	"github.com/ZenLiuCN/mpc"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io/ioutil"
	"os"
	"strings"
)

type Resolver struct {
	git *Git
	//current only on auth method
	Auth transport.AuthMethod
	//Mapping for resolve a Module to git repository
	// key and value must end with slash '/'
	Mapping map[string]string
	//for cache resolved data
	cache Cache
}

/**
resolve a gaven module to it's clone uri and repo name
*/
func (s *Resolver) resolve(module mpc.Module) (uri string, name string, path string) {

	for key, value := range s.Mapping {
		if strings.HasPrefix(string(module), key) {
			/**
			a module from git is prefix with a marker, and remain part is a git uri or uri with paths
			eg: git.pkg/abc/sl => git@ssh.some.com/abc.git and 'sl' consider as a path
			*/
			unPrefix := strings.TrimPrefix(string(module), key)
			idx := strings.Index(unPrefix, "/")
			if idx > 0 {
				x := strings.SplitN(strings.TrimPrefix(string(module), key), "/", 2)
				name = x[0]
				path = x[1]
			} else {
				name = unPrefix
			}
			uri = fmt.Sprintf("%s%s.git", value, name)
			return
		}
	}
	return
}
func (s *Resolver) cloneOrOpen(uri string, name string) (repo *Repo, err error) {
	dir, err := ioutil.TempDir("", "repo_"+name)
	if os.IsExist(err) {
		repo, err = s.git.Open(dir)
		if err != nil {
			return
		}
		err = repo.Pull(s.Auth)
		return
	} else if err != nil {
		return
	}
	repo, err = s.git.Clone(uri, dir, s.Auth) //todo may resolved not a repo
	return
}
func (s *Resolver) Versions(module mpc.Module) mpc.Versions {
	panic("implement me")
}

func (s *Resolver) Info(module mpc.Module, version mpc.Version) *mpc.Info {
	panic("implement me")
}

func (s *Resolver) Mod(module mpc.Module, version mpc.Version) mpc.GoMod {
	panic("implement me")
}

func (s *Resolver) Zip(module mpc.Module, version mpc.Version) mpc.GoZip {
	panic("implement me")
}
