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

import "github.com/ZenLiuCN/mpc"

// cache is a settable resolver
type Cache interface {
	SetVersions(module mpc.Module, versions mpc.Versions)
	Versions(module mpc.Module) mpc.Versions

	SetInfo(module mpc.Module, version mpc.Version, info *mpc.Info)
	Info(module mpc.Module, version mpc.Version) *mpc.Info

	SetMod(module mpc.Module, version mpc.Version, mod mpc.GoMod)
	Mod(module mpc.Module, version mpc.Version) mpc.GoMod

	SetZip(module mpc.Module, version mpc.Version, zip mpc.GoZip)
	Zip(module mpc.Module, version mpc.Version) mpc.GoZip
}
