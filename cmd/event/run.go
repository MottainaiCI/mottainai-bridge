/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
                         Daniele Rondina <geaaru@sabayonlinux.org>

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

package event

import (
	"fmt"

	service "github.com/MottainaiCI/mottainai-bridge/pkg/service"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"

	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newEventRun() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Run event listener",
		Args:  cobra.OnlyValidArgs,
		// TODO: PreRun check of minimal args if --json is not present
		Run: func(cmd *cobra.Command, args []string) {
			var v *viper.Viper = setting.Configuration.Viper

			listener := service.NewClient(client.NewTokenClient(v.GetString("master"), v.GetString("apikey")))

			listener.Listen("task.created", func(c service.TaskMap) {
				fmt.Println("[Task][Create]: ", c)
			})

			listener.Listen("task.removed", func(c service.TaskMap) {
				fmt.Println("[Task][Remove]:", c)
			})

			listener.Listen("task.update", func(TaskUpdates *service.TaskUpdate) {
				fmt.Println("[Task][Update]:", TaskUpdates.Task, "\n[Diff ]", TaskUpdates.Diff)
			})
			listener.Run()
		},
	}

	return cmd
}
