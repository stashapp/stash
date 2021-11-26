package cmd

import (
	"embed"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stashapp/stash/pkg/api"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
)

var (
	uiBox      embed.FS
	loginUIBox embed.FS
	flags      config.FlagStruct

	rootCmd = &cobra.Command{
		Use:   "stash",
		Short: "Stash - An organizer for your porn",
		Long:  "Stash organizes your porn collection, and lets you watch it",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(flags)
			manager.Initialize(flags)
			api.Start(uiBox, loginUIBox)

			// stop any profiling at exit
			defer pprof.StopCPUProfile()
			blockForever()

			err := manager.GetInstance().Shutdown()
			if err != nil {
				logger.Errorf("Error when closing: %s", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Tell Cobra about our flags
	rootCmd.PersistentFlags().IPVar(&flags.Host, "host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	rootCmd.PersistentFlags().IntVar(&flags.Port, "port", 9999, "port to serve from")
	rootCmd.PersistentFlags().StringVarP(&flags.ConfigFilePath, "config", "c", "", "config file to use")
	rootCmd.PersistentFlags().StringVar(&flags.CpuProfilePath, "cpuprofile", "", "write cpu profile to file")
	rootCmd.PersistentFlags().BoolVar(&flags.NoBrowser, "nobrowser", false, "don't open a browser window after launch")

	// Bind the flags on the command line to viper
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("cpuprofile", rootCmd.PersistentFlags().Lookup("cpuprofile"))
	viper.BindPFlag("nobrowser", rootCmd.PersistentFlags().Lookup("nobrowser"))
}

func initConfig() {
}

func Execute(ui, login embed.FS) error {
	uiBox, loginUIBox = ui, login
	return rootCmd.Execute()
}

func blockForever() {
	// handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
