package cmd

import (
	"flag"

	"conductor/pkg/config"
	"conductor/pkg/manager"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ConductorCmd = &cobra.Command{
		Use:           "conductor",
		Short:         "semafour server",
		Long:          "api gateway for the semafour framework",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) error {
			return manager.Spawn()
		},
	}

	addr       string
	port       int
	kubeConfig string
)

func init() {
	viper.AutomaticEnv()

	ConductorCmd.PersistentFlags().StringVarP(&addr, config.Addr, "", "0.0.0.0", "ip address to listen on")
	ConductorCmd.PersistentFlags().IntVarP(&port, config.Port, "", 8080, "port to listen on")
	ConductorCmd.PersistentFlags().StringVarP(&kubeConfig, config.KubeConfig, "", "~/.kube/config", "path to kubeconfig file")

	// parse the go default flagset to get flags for glog and other packages in future
	ConductorCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	//suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})

	viper.BindPFlags(ConductorCmd.PersistentFlags())
	viper.BindPFlags(ConductorCmd.Flags())
}
