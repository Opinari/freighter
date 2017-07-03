package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup FROM_LOCAL_PATH TO_REMOTE_PATH",
	Short: "Backup a file or directory",
	Long: `Perform a backup of a directory as follows:
1. Archive directory as .tar
2. Compress directory as .tar.gz
3. If filename present in storage provider, rename with tagged date for archiving purposes
4. Upload compressed file to storage provider`,
	RunE: runBackup,
}

func init() {
	RootCmd.AddCommand(backupCmd)
}

func runBackup(cmd *cobra.Command, args []string) error {

	if len(args) != 2 {
		return errors.New("'freighter backup' requires two positional arguments: FROM_LOCAL_PATH TO_REMOTE_PATH")
	}

	if err := storageClient.BackupDirectory(args[0], args[1]); err != nil {
		return err
	}

	return nil
}
