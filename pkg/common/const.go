/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package common

const (
	MBRIDGE_ENV_PREFIX  = "MOTTAINAI_BRIDGE"
	MBRIDGE_CONFIG_NAME = "mbridge-profiles"
	// NOTE: doesn't use $HOME because os.Mkdir doesn't resolve it.
	MBRIDGE_HOME_PATH  = ".config/mottainai"
	MBRIDGE_LOCAL_PATH = ".mottainai"
)
