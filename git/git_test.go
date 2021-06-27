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
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"testing"
)
import "github.com/stretchr/testify/assert"

const (
	PEM = ""
	SSH = ""
	HTT = ""
)

func TestGitHttpClone(t *testing.T) {
	g := Git(0)
	r, err := g.Clone(HTT)
	assert.Nil(t, err)
	t.Logf("%+v ", r)
}
func TestGitSSHClone(t *testing.T) {
	g := Git(0)
	auth, err := ssh.NewPublicKeys("git", []byte(PEM), "")
	assert.Nil(t, err)
	r, dir, err := g.CloneAuth(SSH, auth)
	assert.Nil(t, err)
	t.Logf("%+v %s", r, dir)
}
