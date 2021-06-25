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
	"errors"
	"sort"
)

var (
	names         = map[int]string{}
	factories     = map[int]ResolverFactory{}
	resolverIndex = make([]int, 0, 5)
)

func RegisterResolver(name string, order int, factory ResolverFactory) error {
	if _, ok := names[order]; ok {
		return errors.New("order is already exists")
	}
	names[order] = name
	factories[order] = factory
	resolverIndex = append(resolverIndex, order)
	return nil
}

var (
	checksum      = map[int]CheckSumResolver{}
	checkSumIndex = make([]int, 0, 5)
)

func RegisterCheckSumResolver(order int, resolver CheckSumResolver) error {
	if _, ok := checksum[order]; ok {
		return errors.New("order is already exists")
	}
	checksum[order] = resolver
	checkSumIndex = append(checkSumIndex, order)
	return nil
}

// prepare resolvers
func Initial() {
	sort.Ints(resolverIndex)
	resolvers = make([]Resolver, 0, len(resolverIndex))
	for _, index := range resolverIndex {
		resolvers = append(resolvers, factories[index](resolvers...))
	}
	sort.Ints(checkSumIndex)
}

// sorted if after Initial.
func ResolverNames() (r []string) {
	r = make([]string, 0, len(resolverIndex))
	for _, index := range resolverIndex {
		r = append(r, names[index])
	}
	return
}

// sorted if after Initial.
func ResolverFactories() (r []ResolverFactory) {
	sort.Ints(resolverIndex)
	r = make([]ResolverFactory, 0, len(resolverIndex))
	for _, index := range resolverIndex {
		r = append(r, factories[index])
	}
	return
}

var (
	resolvers []Resolver
)

// nil if before Initial.
func Resolvers() (r []Resolver) {
	return resolvers
}

//fetch the Versions of module, if UNKNOWN just return nil
func ResolveVersions(module Module) Versions {
	for _, resolver := range resolvers {
		if v := resolver.Versions(module); v != "" {
			return v
		}
	}
	return ""
}

// fetch the Info of a module with version, if UNKNOWN just return nil
// version may is Latest
func ResolveInfo(module Module, version Version) *Info {
	for _, resolver := range resolvers {
		if v := resolver.Info(module, version); v != nil {
			return v
		}
	}
	return nil
}

// fetch the GoMod of a module with version, if UNKNOWN just return empty
func ResolveMod(module Module, version Version) GoMod {
	for _, resolver := range resolvers {
		if v := resolver.Mod(module, version); v != "" {
			return v
		}
	}
	return ""
}

// fetch the GoZip of a module with version, if UNKNOWN just return nil
func ResolveZip(module Module, version Version) GoZip {
	for _, resolver := range resolvers {
		if v := resolver.Zip(module, version); v != nil {
			return v
		}
	}
	return nil
}

func SumResolveSupported() bool {
	for _, r := range checkSumIndex {
		if checksum[r].Supported() {
			return true
		}
	}
	return false
}

//$base/latest
func SumResolveLatest() []byte {
	for _, r := range checkSumIndex {
		if m := checksum[r].Latest(); m != nil {
			return m
		}
	}
	return nil
}

//$base/lookup/$module@$version
func SumResolveLookup(module Module, version Version) []byte {
	for _, r := range checkSumIndex {
		if m := checksum[r].Lookup(module, version); m != nil {
			return m
		}
	}
	return nil
}

//$base/tile/$H/$L/$K[.p/$W]  also process tile data $base/tile/$H/data/$K[.p/$W]
func SumResolveTile(path string) []byte {
	for _, r := range checkSumIndex {
		if m := checksum[r].Tile(path); m != nil {
			return m
		}
	}
	return nil
}
