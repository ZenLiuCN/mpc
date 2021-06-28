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
	"encoding/json"
	"io"
	"path"
	"strings"
	"time"
)

// a Module is a string present a module path
type Module string

const UndefinedModule Module = ""

//  Version
type Version string

const LatestVersion Version = "latest"
const UndefinedVersion Version = ""

// for $base/$module/@v/list
type Versions string

// for $base/$module/@v/$version.info
type Info struct {
	Version Version
	Time    time.Time //commit time
}

func (i Info) Marshal() []byte {
	b, _ := json.Marshal(i)
	return b
}
func (i *Info) UnMarshal(b []byte) error {
	return json.Unmarshal(b, i)
}

// for $base/$module/@v/$version.mod
type GoMod string

//for $base/$module/@v/$version.zip
type GoZip interface {
	io.ReadCloser
}

// Resolver is a GO PROXY command processor
// 1. can be a proxy to other GO Proxy server
// 2. can be a local cache
// 3. can be a customer processor to fetch from private locations (git repository for example).
type Resolver interface {
	//fetch the Versions of module, if UNKNOWN just return nil
	Versions(module Module) Versions
	// fetch the Info of a module with version, if UNKNOWN just return nil
	// version may is Latest
	Info(module Module, version Version) *Info
	// fetch the GoMod of a module with version, if UNKNOWN just return empty
	Mod(module Module, version Version) GoMod
	// fetch the GoZip of a module with version, if UNKNOWN just return nil
	Zip(module Module, version Version) GoZip
}

// ResolverFactory
type ResolverFactory func(resolvers ...Resolver) Resolver

type Cmd int

func (c Cmd) String() string {
	switch c {
	case CmdLatest:
		return "CmdLatest"
	case CmdList:
		return "CmdList"
	case CmdInfo:
		return "CmdInfo"
	case CmdMod:
		return "CmdMod"
	case CmdZip:
		return "CmdZip"
	default:
		return "CmdUndefined"
	}
}

const (
	CmdUndefined Cmd = iota
	CmdLatest        // $base/$module/@latest
	CmdList          // $base/$module/@v/list
	CmdInfo          //$base/$module/@v/$version.info
	CmdMod           // $base/$module/@v/$version.mod
	CmdZip           // $base/$module/@v/$version.zip
)

type SumCmd int

func (s SumCmd) String() string {
	switch s {
	case SumSupported:
		return "SumSupported"
	case SumLatest:
		return "SumLatest"
	case SumLookup:
		return "SumLookup"
	case SumTile:
		return "SumTile"
	default:
		return "SumUndefined"
	}
}

const (
	sumPrefix          = "sumdb/"
	SumSupportedSuffix = "supported"
	SumLatestSuffix    = "latest"
	SumLookupPrefix    = "lookup/"
	SumTilePrefix      = "tile/"

	SumUndefined SumCmd = iota
	SumSupported
	SumLatest
	SumLookup
	SumTile
)

// CheckSumResolver for a Local SumDB or Proxy SumDB
type CheckSumResolver interface {
	Supported() bool
	//$base/latest
	Latest() []byte
	//$base/lookup/$module@$version
	Lookup(module Module, version Version) []byte
	//$base/tile/$H/$L/$K[.p/$W]  also process tile data $base/tile/$H/data/$K[.p/$W]
	Tile(path string) []byte
}

// parser a command
// @parameter requestPath request url path without base part.
// @return (Module , Version , Cmd , SumCmd,sum parameter)
// Use Cmd and SumCmd to validate result.
// sum parameter only for SumTile and SumTileData.
// if is a Cmd: SumCmd is SumUndefined.
// if is a SumCmd: Cmd is CmdUndefined.
// if both Cmd is CmdUndefined and SumCmd is SumUndefined means the path can not be parsed.
func CommandParser(requestPath string) (m Module, v Version, c Cmd, s SumCmd, p string) {
	req := path.Clean(requestPath)
	if strings.HasPrefix(req, sumPrefix) {
		return sumCmdProcessor(strings.TrimPrefix(req, sumPrefix))
	} else {
		return cmdProcessor(req)
	}
}
func cmdProcessor(req string) (m Module, v Version, c Cmd, s SumCmd, p string) {
	switch {
	case strings.HasSuffix(req, "latest"): // $base/$module/@latest
		v = LatestVersion
		c = CmdLatest
		m = Module(strings.TrimSuffix(req, "/@latest"))
	case strings.HasSuffix(req, "/@v/list"): // $base/$module/@v/list
		v = UndefinedVersion
		c = CmdList
		m = Module(strings.TrimSuffix(req, "/@v/list"))
	case strings.HasSuffix(req, ".info"): //$base/$module/@v/$version.info
		x := strings.TrimSuffix(req, ".info")
		ss := strings.Split(x, "/@v/")
		if len(ss) != 2 {
			return
		}
		c = CmdInfo
		v = Version(ss[1])
		m = Module(ss[0])
	case strings.HasSuffix(req, ".mod"): //$base/$module/@v/$version.mod
		x := strings.TrimSuffix(req, ".mod")
		ss := strings.Split(x, "/@v/")
		if len(ss) != 2 {
			return
		}
		c = CmdMod
		v = Version(ss[1])
		m = Module(ss[0])
	case strings.HasSuffix(req, ".zip"): //$base/$module/@v/$version.zip
		x := strings.TrimSuffix(req, ".zip")
		ss := strings.Split(x, "/@v/")
		if len(ss) != 2 {
			return
		}
		c = CmdZip
		v = Version(ss[1])
		m = Module(ss[0])
	}
	return
}
func sumCmdProcessor(req string) (m Module, v Version, c Cmd, s SumCmd, p string) {
	switch {
	case strings.HasSuffix(req, SumSupportedSuffix): //$base/supported
		s = SumSupported
		v = UndefinedVersion
		m = UndefinedModule
		p = ""
	case req == SumLatestSuffix: //$base/latest
		s = SumLatest
		v = UndefinedVersion
		m = UndefinedModule
		p = ""
	case strings.HasPrefix(req, SumLookupPrefix): //$base/lookup/$module@$version
		s = SumLookup
		ss := strings.Split(strings.TrimPrefix(req, "lookup/"), "@")
		if len(ss) != 2 {
			return
		}
		m = Module(ss[0])
		v = Version(ss[1])
	case strings.HasPrefix(req, SumTilePrefix): //$base/tile/$H/$L/$K[.p/$W] or $base/tile/$H/data/$K[.p/$W]
		s = SumTile
		v = UndefinedVersion
		m = UndefinedModule
		p = strings.TrimPrefix(req, SumTilePrefix)
	}
	return
}
func BuildCmd(proxy string, cmd Cmd, module Module, version Version) string {
	switch cmd {
	case CmdLatest:
		return fmt.Sprintf("%s/%s/@latest", proxy, module)
	case CmdList:
		return fmt.Sprintf("%s/%s/@v/list", proxy, module)
	case CmdInfo:
		return fmt.Sprintf("%s/%s/@v/%s.info", proxy, module, version)
	case CmdMod:
		return fmt.Sprintf("%s/%s/@v/%s.mod", proxy, module, version)
	case CmdZip:
		return fmt.Sprintf("%s/%s/@v/%s.zip", proxy, module, version)
	default:
		return ""
	}
}
func BuildSumCmd(proxy string, cmd SumCmd, module Module, version Version, param string) string {
	switch cmd {
	case SumSupported:
		return fmt.Sprintf("%s/supported", proxy)
	case SumLatest:
		return fmt.Sprintf("%s/latest", proxy)
	case SumLookup:
		return fmt.Sprintf("%s/lookup/%s@%s", proxy, module, version)
	case SumTile:
		return fmt.Sprintf("%s/tile/%s", proxy, param)
	default:
		return ""
	}
}
