package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete REMOTE_FILE_PATH [REMOTE_FILE_PATH...]",
	Short: "Delete file(s) hosted in storage provider",
	Long: `Delete one or more files hosted in the chosen
storage provider.`,
	RunE: runDelete,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return errors.New("'freighter delete' requires at least one positional argument: REMOTE_FILE_PATH [REMOTE_FILE_PATH...]")
	}

	for _, arg := range args {
		if err := storageClient.DeleteRemoteFile(arg); err != nil {
			return err
		}
	}

	return nil
}
