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

type Bit64 uint64

func (s Bit64) At(pos int) bool {
	if pos < 0 || pos >= 64 {
		return false
	}
	return 1&(s>>pos) != 0
}
func (s Bit64) Set(pos int) Bit64 {
	if pos < 0 || pos >= 64 {
		return s
	}
	return s | (1 << pos)
}
func (s Bit64) Reset(pos int) Bit64 {
	if pos < 0 || pos >= 64 {
		return s
	}
	return s & (1 << pos)
}

type Bit32 uint32

func (s Bit32) At(pos int) bool {
	if pos < 0 || pos >= 32 {
		return false
	}
	return 1&(s>>pos) != 0
}
func (s Bit32) Set(pos int) Bit32 {
	if pos < 0 || pos >= 32 {
		return s
	}
	return s | (1 << pos)
}
func (s Bit32) Reset(pos int) Bit32 {
	if pos < 0 || pos >= 32 {
		return s
	}
	return s & (1 << pos)
}

type Bit16 uint16

func (s Bit16) At(pos int) bool {
	if pos < 0 || pos >= 16 {
		return false
	}
	return 1&(s>>pos) != 0
}
func (s Bit16) Set(pos int) Bit16 {
	if pos < 0 || pos >= 16 {
		return s
	}
	return s | (1 << pos)
}
func (s Bit16) Reset(pos int) Bit16 {
	if pos < 0 || pos >= 16 {
		return s
	}
	return s & (1 << pos)
}

type Bit8 uint8

func (s Bit8) At(pos int) bool {
	if pos < 0 || pos >= 8 {
		return false
	}
	return 1&(s>>pos) != 0
}
func (s Bit8) Set(pos int) Bit8 {
	if pos < 0 || pos >= 8 {
		return s
	}
	return s | (1 << pos)
}
func (s Bit8) Reset(pos int) Bit8 {
	if pos < 0 || pos >= 8 {
		return s
	}
	return s & (1 << pos)
}
