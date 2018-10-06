package cmd

import (
	"fmt"
	"strings"

	"github.com/fvdveen/mu2-config"
	provider "github.com/fvdveen/mu2-config/consul"
	"github.com/fvdveen/mu2-config/events"
	searchpb "github.com/fvdveen/mu2-proto/go/proto/search"
	"github.com/fvdveen/mu2-search/watch"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLvl string
	conf   struct {
		Consul struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"consul"`
		Config struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		} `mapstructure:"config"`
	}
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mu2-search",
	Short: "Mu2 search service",
	Long: `Mu2 is a discord music bot. 

This is the search service for mu2.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cc := api.DefaultConfig()
		if conf.Consul.Address != "" {
			cc.Address = conf.Consul.Address
		}

		var p config.Provider

		srv := grpc.NewService(
			micro.Name("mu2.service.search"),
			micro.Version("latest"),
			micro.Registry(consul.NewRegistry(consul.Config(cc))),
			micro.AfterStop(func() error {
				return p.Close()
			}),
		)

		c, err := api.NewClient(cc)
		if err != nil {
			return fmt.Errorf("Create consul client: %v", err)
		}

		p, err = provider.NewProvider(c, conf.Config.Path, conf.Config.Type, nil)
		if err != nil {
			return fmt.Errorf("could not create provider: %v", err)
		}

		ych, rch := events.Youtube(events.Watch(p.Watch()))
		events.Null(rch)

		ss := watch.Youtube(ych)

		searchpb.RegisterSearchServiceHandler(srv.Server(), ss)

		return srv.Run()
	},
	SilenceUsage: true,
}

// Execute runs the cli
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MU2")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&conf.Consul.Address, "consul-addr", "", "consul address")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Path, "config-path", "search/config", "config path on the kv store")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Type, "config-type", "json", "config type on the kv store")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logrus.WithField("type", "main").Fatalf("Unmarshalling config: %v", err)
		return
	}

	var lvl logrus.Level

	logrus.SetLevel(lvl)
}
