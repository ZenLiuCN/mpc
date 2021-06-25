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

package mpc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

type JustTestResolver int

func (j JustTestResolver) Versions(module Module) Versions {
	return "1\n2\n3"
}

func (j JustTestResolver) Info(module Module, version Version) *Info {
	return &Info{
		Version: "1",
		Time:    time.Now(),
	}
}

func (j JustTestResolver) Mod(module Module, version Version) GoMod {
	return "module github.com/ZenLiuCN/mpc\n\ngo 1.14\n"
}

func (j JustTestResolver) Zip(module Module, version Version) GoZip {
	return ioutil.NopCloser(bytes.NewBufferString("NOT A ZIP"))
}

func TestRegisterResolver(t *testing.T) {
	err := RegisterResolver("JustTestResolver1", 0, func(resolvers ...Resolver) Resolver {
		return JustTestResolver(0)
	})
	assert.Nil(t, err, "register resolver err: %+v", err)
	assert.Equal(t, []string{"JustTestResolver1"}, ResolverNames())
	err = RegisterResolver("JustTestResolver2", 0, func(resolvers ...Resolver) Resolver {
		return JustTestResolver(0)
	})
	assert.NotNil(t, err)
	err = RegisterResolver("JustTestResolver3", -1, func(resolvers ...Resolver) Resolver {
		return JustTestResolver(0)
	})
	assert.NotNil(t, ResolverNames())
	assert.Equal(t, []string{"JustTestResolver1", "JustTestResolver3"}, ResolverNames())
	Initial()
	assert.Equal(t, []string{"JustTestResolver3", "JustTestResolver1"}, ResolverNames())
}
