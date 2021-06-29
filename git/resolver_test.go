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
	"github.com/go-git/go-git/v5/plumbing/transport"
	"testing"
)

func TestResolver_resolve(t *testing.T) {
	type fields struct {
		git         *Git
		Auth        transport.AuthMethod
		Mapping     map[string]string
		mappingKeys []string
	}
	type args struct {
		module mpc.Module
	}
	var tests = []struct {
		name     string
		fields   fields
		args     args
		wantUri  string
		wantName string
		wantPath string
	}{
		{
			name: "resolve test",
			fields: fields{Mapping: map[string]string{
				"git.x/": "git@git.x.com/",
			}},
			args: args{
				"git.x/some/path",
			},
			wantUri:  "git@git.x.com/some.git",
			wantName: "some",
			wantPath: "path",
		},
		{
			name: "resolve test without path",
			fields: fields{Mapping: map[string]string{
				"git.x/": "git@git.x.com/",
			}},
			args: args{
				"git.x/some",
			},
			wantUri:  "git@git.x.com/some.git",
			wantName: "some",
			wantPath: "",
		}, {
			name: "resolve test long path",
			fields: fields{Mapping: map[string]string{
				"git.x/": "git@git.x.com/",
			}},
			args: args{
				"git.x/some/p1/p2/p3",
			},
			wantUri:  "git@git.x.com/some.git",
			wantName: "some",
			wantPath: "p1/p2/p3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Resolver{
				git:     tt.fields.git,
				Auth:    tt.fields.Auth,
				Mapping: tt.fields.Mapping,
			}
			gotUri, gotName, gotPath := s.resolve(tt.args.module)
			if gotUri != tt.wantUri {
				t.Errorf("resolve() gotUri = %v, want %v", gotUri, tt.wantUri)
			}
			if gotName != tt.wantName {
				t.Errorf("resolve() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotPath != tt.wantPath {
				t.Errorf("resolve() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
