package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore FROM_REMOTE_PATH TO_LOCAL_PATH",
	Short: "Restore a file or directory",
	Long: `Perform a restore of a directory as follows:
1. Download compressed file from storage provider
2. Uncompress file from .tar.gz to .tar
3. Unarchive directory from .tar`,
	RunE: runRestore,
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {

	if len(args) != 2 {
		return errors.New("'freighter restore' requires two positional arguments: FROM_REMOTE_PATH TO_LOCAL_PATH")
	}

	if err := storageClient.RestoreFile(args[0], args[1]); err != nil {
		return err
	}

	return nil
}
