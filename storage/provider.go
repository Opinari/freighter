package storage

// A StorageProvider provides the ability to transit / manipulate files in the below fashion
//
// An implementation must provide the guarantee that the operation was atomically complete before returning or
// should return an error to indicate why this was not the case.
type StorageProvider interface {
	DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error)
	UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error)
	DeleteFile(remoteFilePath string) (deleteFilePath string, err error)
	AgeFile(remoteFilePath string) (ageOfFile int, err error)
}
