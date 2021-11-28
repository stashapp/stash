// package cmd is the CLI command handling for the stash application
package cmd

import (
	"context"
	"embed"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	host       net.IP
	port       int
	conf       *viper.Viper

	rootCmd = &cobra.Command{
		Use:   "stash",
		Short: "Stash - An organizer for your porn",
		Long:  "Stash is an organizer for your porn, with search and watch functionality",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cfg, err := config.Initialize(flags, conf)
			if err != nil {
				panic(fmt.Sprintf("error initializing configuration: %s", err))
			}
			manager.Initialize(ctx, cfg)
			api.Start(uiBox, loginUIBox)

			// stop any profiling at exit
			defer pprof.StopCPUProfile()
			blockForever()

			err = manager.GetInstance().Shutdown()
			if err != nil {
				logger.Errorf("Error when closing: %s", err)
			}
		},
	}
)

func init() {
	conf = viper.New()
	cobra.OnInitialize(initConfig)

	// Tell Cobra about our flags
	rootCmd.PersistentFlags().IPVar(&host, "host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	rootCmd.PersistentFlags().IntVar(&port, "port", 9999, "port to serve from")
	rootCmd.PersistentFlags().StringVarP(&flags.ConfigFilePath, "config", "c", "", "config file to use")
	rootCmd.PersistentFlags().StringVar(&flags.CpuProfilePath, "cpuprofile", "", "write cpu profile to file")
	rootCmd.PersistentFlags().BoolVar(&flags.NoBrowser, "nobrowser", false, "don't open a browser window after launch")

	// Bind the flags on the command line to viper
	bindPFlag(conf, "host", rootCmd.PersistentFlags().Lookup("host"))
	bindPFlag(conf, "port", rootCmd.PersistentFlags().Lookup("port"))
	bindPFlag(conf, "config", rootCmd.PersistentFlags().Lookup("config"))
	bindPFlag(conf, "cpuprofile", rootCmd.PersistentFlags().Lookup("cpuprofile"))
	bindPFlag(conf, "nobrowser", rootCmd.PersistentFlags().Lookup("nobrowser"))

	// Bind to the environment as well
	conf.SetEnvPrefix("stash")     // will be uppercased automatically
	bindEnv(conf, "host")          // STASH_HOST
	bindEnv(conf, "port")          // STASH_PORT
	bindEnv(conf, "external_host") // STASH_EXTERNAL_HOST
	bindEnv(conf, "generated")     // STASH_GENERATED
	bindEnv(conf, "metadata")      // STASH_METADATA
	bindEnv(conf, "cache")         // STASH_CACHE
	bindEnv(conf, "stash")         // STASH_STASH
}

func bindPFlag(viper *viper.Viper, key string, flag *pflag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		panic(fmt.Sprintf("unable to bind to pflag: %v", err))
	}
}

func bindEnv(viper *viper.Viper, key string) {
	if err := viper.BindEnv(key); err != nil {
		panic(fmt.Sprintf("unable to set environment key (%v): %v", key, err))
	}
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
