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
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	pathPrefix = "/"
	CacheAge   = 86400
)

// will call Initial
func InitialHandler(prefix string) {
	if prefix != "" {
		pathPrefix = prefix
	}
	Initial()
}
func GoProxyHandler(w http.ResponseWriter, r *http.Request) {
	re := res{w}
	if strings.HasPrefix(r.URL.Path, pathPrefix) {
		cmd := strings.TrimPrefix(r.URL.Path, pathPrefix)
		m, v, c, s, p := CommandParser(cmd)
		switch c {
		case CmdList:
			i := ResolveVersions(m)
			if i != "" {
				re.okCache([]byte(i))
				return
			}
		case CmdInfo, CmdLatest:
			i := ResolveInfo(m, v)
			if i != nil {
				re.okCache(i.Marshal())
				return
			}
		case CmdMod:
			i := ResolveMod(m, v)
			if i != "" {
				re.okCache([]byte(i))
				return
			}
		case CmdZip:
			i := ResolveZip(m, v)
			if i != nil {
				re.okCacheReader(i)
				return
			}
		case CmdUndefined:
			switch s {
			case SumSupported:
				if SumResolveSupported() {
					re.okCache(nil)
					return
				}
			case SumLatest:
				i := SumResolveLatest()
				if i != nil {
					re.contentText()
					re.okCache(i)
					return
				}
			case SumLookup:
				i := SumResolveLookup(m, v)
				if i != nil {
					re.contentText()
					re.okCache(i)
					return
				}
			case SumTile:
				i := SumResolveTile(p)
				if i != nil {
					re.contentStream()
					re.okCache(i)
					return
				}
			}

		}
	}
	re.notFoundCache()
}

type res struct {
	http.ResponseWriter
}

func (r res) notFoundCache() {
	r.writeCache(CacheAge)
	r.WriteHeader(404)
}
func (r res) contentText() {
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
func (r res) contentStream() {
	r.Header().Set("Content-Type", "application/octet-stream")
}
func (r res) writeCache(age int) {
	if age <= 0 {
		r.Header().Set("Cache-Control", "must-revalidate, no-cache, no-store")
	} else {
		r.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", age))
	}
}
func (r res) okCache(data []byte) {
	r.writeCache(CacheAge)
	_, _ = r.Write(data)
}
func (r res) okNoCache(data []byte) {
	r.writeCache(0)
	_, _ = r.Write(data)
}
func (r res) ok(data []byte) {
	_, _ = r.Write(data)
}
func (r res) okCacheReader(data io.ReadCloser) {
	r.writeCache(CacheAge)
	c := make([]byte, 0, 1024)
	defer data.Close()
	n, err := data.Read(c)
	if err != nil {
		r.WriteHeader(500)
	}
	for n > 0 {
		_, _ = r.Write(c[:n])
		n, err = data.Read(c)
	}
}
