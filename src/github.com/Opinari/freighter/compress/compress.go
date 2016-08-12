package compress

import (
	"fmt"
	_ "compress/gzip"
)

func CompressFile(inputDirPath string, outputFilePath string) (archivedFilePath string) {

	fmt.Printf("Compressing File from: %s to: %s", inputDirPath, outputFilePath);

	archivedFilePath = outputFilePath;
	return
}

func UncompressFile(inputFilePath string, outputDirPath string) (uncompressDirPath string) {

	fmt.Printf("Uncompressing File from: %s to: %s", inputFilePath, outputDirPath);

	uncompressDirPath = outputDirPath;
	return
}
