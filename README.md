# freighter


Frieghter is a utility to help move files to / from storage backends in a convenient and efficient manner.
 
Use cases may include:

- Restoring some database backup files to a cloud compute instance
- Uploading local files to Dropbox for safe storage
- Downloading configuration files to a host from a VCS such as Github


<h2>Currently Supported Backends</h2>

- Dropbox (v2 API)
- Github (v3 API)



<h2>Restore Files</h2>


With Docker:


```sh
docker run -u [USER] -v [LOCAL_VOLUME_MAPPING]:/restore -e BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] \
  opinari/freighter restore --remoteFilePath [REMOTE_FILE_PATH] --restoreFilePath /restore/[LOCAL_FILE_PATH] 
```

With native binary


```sh
env BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] ./freighter restore --remoteFilePath [REMOTE_FILE_PATH] --restoreFilePath [LOCAL_FILE_PATH] 
```


<h2>Backup Files</h2>


With Docker:


```sh
docker run -u [USER] -v [LOCAL_VOLUME_MAPPING]:/backup -e BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] \
  opinari/freighter backup --remoteFilePath [REMOTE_FILE_PATH] --backupFilePath /backup/[LOCAL_FILE_PATH] 
```

With native binary


```sh
env BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] ./freighter backup --remoteFilePath [REMOTE_FILE_PATH] --backupFilePath [LOCAL_FILE_PATH] 
```


<h2>Determine Age of File</h2>

To determine the age of a file (in days):


With Docker:


```sh
docker run -u [USER] -v [LOCAL_VOLUME_MAPPING]:/age -e BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] \
  opinari/freighter age --remoteFilePath [REMOTE_FILE_PATH] --ageOutputFilePath /age/[LOCAL_FILE_PATH] 
```


With native binary

```sh
env BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] /otf/freighter age --remoteFilePath /"$backupArchiveName".tar.gz --ageOutputFilePath /tmp/age
```

<h2>Delete File(s)</h2>


With Docker:


```sh
docker run -u [USER] -v [LOCAL_VOLUME_MAPPING]:/age -e BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] \
  opinari/freighter delete --remoteFilePath [REMOTE_FILE_PATH]
```


With native binary

```sh
env BACKUP_PROVIDER_TOKEN=[BACKUP_PROVIDER_TOKEN] /otf/freighter delete --remoteFilePath /"$backupArchiveName".tar.gz
```



<hr>
