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

import "testing"

func TestCommandParser(t *testing.T) {
	type args struct {
		requestPath string
	}
	tests := []struct {
		name  string
		args  args
		wantM Module
		wantV Version
		wantC Cmd
		wantS SumCmd
		wantP string
	}{
		{
			name:  "latestCmd",
			args:  args{"github.com/ZenLiuCn/mpc/@latest"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: LatestVersion,
			wantC: CmdLatest,
			wantS: 0,
			wantP: "",
		},
		{
			name:  "listCmd",
			args:  args{"github.com/ZenLiuCn/mpc/@v/list"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: UndefinedVersion,
			wantC: CmdList,
			wantS: 0,
			wantP: "",
		}, {
			name:  "infoCmd",
			args:  args{"github.com/ZenLiuCn/mpc/@v/1.1.2.info"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: Version("1.1.2"),
			wantC: CmdInfo,
			wantS: 0,
			wantP: "",
		}, {
			name:  "modCmd",
			args:  args{"github.com/ZenLiuCn/mpc/@v/1.1.2.mod"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: Version("1.1.2"),
			wantC: CmdMod,
			wantS: 0,
			wantP: "",
		}, {
			name:  "zipCmd",
			args:  args{"github.com/ZenLiuCn/mpc/@v/1.1.2.zip"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: Version("1.1.2"),
			wantC: CmdZip,
			wantS: 0,
			wantP: "",
		}, {
			name:  "latestSum",
			args:  args{"sumdb/latest"},
			wantM: "",
			wantV: "",
			wantC: 0,
			wantS: SumLatest,
			wantP: "",
		}, {
			name:  "lookupSum",
			args:  args{"sumdb/lookup/github.com/ZenLiuCn/mpc@1.1.2"},
			wantM: "github.com/ZenLiuCn/mpc",
			wantV: Version("1.1.2"),
			wantC: 0,
			wantS: SumLookup,
			wantP: "",
		}, {
			name:  "tileSum1",
			args:  args{"sumdb/tile/1/2/3"},
			wantM: "",
			wantV: "",
			wantC: 0,
			wantS: SumTile,
			wantP: "1/2/3",
		}, {
			name:  "tileSum2",
			args:  args{"sumdb/tile/1/2/3.p/4"},
			wantM: "",
			wantV: "",
			wantC: 0,
			wantS: SumTile,
			wantP: "1/2/3.p/4",
		}, {
			name:  "tileDataSum1",
			args:  args{"sumdb/tile/1/0/3.p/4"},
			wantM: "",
			wantV: "",
			wantC: 0,
			wantS: SumTile,
			wantP: "1/0/3.p/4",
		}, {
			name:  "tileDataSum2",
			args:  args{"sumdb/tile/1/0/3"},
			wantM: "",
			wantV: "",
			wantC: 0,
			wantS: SumTile,
			wantP: "1/0/3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, gotV, gotC, gotS, gotP := CommandParser(tt.args.requestPath)
			if gotM != tt.wantM {
				t.Errorf("CommandParser() gotM = %v, want %v", gotM, tt.wantM)
			}
			if gotV != tt.wantV {
				t.Errorf("CommandParser() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotC != tt.wantC {
				t.Errorf("CommandParser() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotS != tt.wantS {
				t.Errorf("CommandParser() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotP != tt.wantP {
				t.Errorf("CommandParser() gotP = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
