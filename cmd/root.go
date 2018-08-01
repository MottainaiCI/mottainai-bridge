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

package cmd

import (
	"fmt"
	"os"
	//	"reflect"
	event "github.com/MottainaiCI/mottainai-bridge/cmd/event"

	utils "github.com/MottainaiCI/mottainai-server/pkg/utils"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"

	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
)

const (
	cliName = `Mottainai Bridge
Copyright (c) 2017-2018 Mottainai

Command line interface for Mottainai bridges`

	cliExamples = `$> mottainai-bridge -m http://127.0.0.1:8080 -k token run

$> mottainai-bridge -m http://127.0.0.1:8080 -k token run
`
)

var rootCmd = &cobra.Command{
	Short:        cliName,
	Version:      setting.MOTTAINAI_VERSION,
	Example:      cliExamples,
	Args:         cobra.OnlyValidArgs,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var v *viper.Viper = setting.Configuration.Viper

		v.SetConfigFile(v.Get("config").(string))
		// Parse configuration file
		err = setting.Configuration.Unmarshal()
		utils.CheckError(err)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {

	var pflags = rootCmd.PersistentFlags()
	v := setting.Configuration.Viper

	pflags.StringP("master", "m", "http://localhost:8080", "MottainaiCI webUI URL")
	pflags.StringP("apikey", "k", "fb4h3bhgv4421355", "Mottainai API key")

	v.BindPFlag("master", rootCmd.PersistentFlags().Lookup("master"))
	v.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))

	rootCmd.AddCommand(
		event.NewEventCommand(),
	)
}

func Execute() {

	// Start command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
