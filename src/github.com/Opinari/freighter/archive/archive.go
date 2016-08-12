package archive

import (
	"fmt"
	_ "archive/tar"
)

func ArchiveFile(inputDirPath string, outputFilePath string) (archivedFilePath string) {

	fmt.Printf("Archiving Files from: %s to: %s", inputDirPath, outputFilePath);

	archivedFilePath = outputFilePath;
	return
}

func UnarchiveFile(inputFilePath string, outputDirPath string) (unarchiveDirPath string) {

	fmt.Printf("Unarchiving Files from: %s to: %s", inputFilePath, outputDirPath);

	unarchiveDirPath = outputDirPath;
	return
}
