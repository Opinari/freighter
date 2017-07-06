package cmd

import (
	"fmt"
	"github.com/opinari/freighter/storage"
	"github.com/opinari/freighter/storage/dropbox"
	"github.com/opinari/freighter/storage/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const version = "0.5.0"

var (
	configFile    string
	storageClient storage.StorageClient
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "freighter",
	Short: "Freighter is a utility to move files from/to storage backends",
	Long: `A fast and minimalistic command-line program which can be used to:
  - Restore database backup files to a cloud compute instance.
  - Upload local files to Dropbox for safe storage.
  - Download configuration files to a host from a VCS such as Github.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Freighter v%s\n", version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "configuration file (default /tmp/.freighter.yml)")
	RootCmd.Flags().BoolP("version", "v", false, "print freighter's version")
}

// Read in configuration file containing provider type and access token.
// Fall back to using environment variables if the configuration file is
// not found.
func initConfig() {

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Search config in tmp directory with name ".freighter" (without extension).
		viper.AddConfigPath(os.TempDir())
		viper.SetConfigName(".freighter")
	}

	// Check for environment variables in all caps prefixed with freighter
	// as a fallback strategy if reading the configuration file fails.
	viper.AutomaticEnv()
	viper.SetEnvPrefix("freighter")
	viper.ReadInConfig()

	if viper.GetString("provider") == "" || viper.GetString("token") == "" {
		log.Fatalln("Provider and/or token not found in either configuration file or environment variables")
	}

	storageProvider := resolveStorageProvider(viper.GetString("provider"), viper.GetString("token"))
	storageClient = storage.NewStorageClient(storageProvider)
}

func resolveStorageProvider(storageProvider, backupProviderToken string) storage.StorageProvider {
	switch storageProvider {
	case "dropbox":
		return dropbox.NewDropboxStorageClient(backupProviderToken)
	case "github":
		return github.NewGithubStorageProvider(backupProviderToken)
	default:
		log.Fatalln("Error: Invalid Provider Type given: " + storageProvider)
		return nil
	}
}
