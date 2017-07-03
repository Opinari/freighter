package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var outputAgePath string

var ageCmd = &cobra.Command{
	Use:   "age REMOTE_FILE_PATH",
	Short: "Determine age of file",
	Long: `Determine age of a file hosted in the specified storage provider
by calculating the number of days since its last modified time until today.
The result is then optionally written to a file.`,
	RunE: runAge,
}

func init() {

	ageCmd.PersistentFlags().StringVarP(&outputAgePath, "output", "o", "", "The local file path to write the age result")

	RootCmd.AddCommand(ageCmd)
}

func runAge(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return errors.New("'freighter age' requires one positional argument: REMOTE_FILE_PATH")
	}

	if err := storageClient.AgeRemoteFile(args[0], outputAgePath); err != nil {
		return err
	}

	return nil
}
