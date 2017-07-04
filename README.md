# freighter


Frieghter is a utility to help move files to / from storage backends in a convenient and efficient manner.
 
Use cases may include:

- Restoring some database backup files to a cloud compute instance
- Uploading local files to Dropbox for safe storage
- Downloading configuration files to a host from a VCS such as Github


<h2>Currently Supported Backends</h2>

- Dropbox (v2 API)
- Github (v3 API)


<h2>Configuration File</h2>
Freighter requires the presence of a .freighter file with any of the following supported extensions:
json, toml, yaml, yml, properties, props, prop, hcl

This configuration file currently supports two keys:

- Provider: Type of storage provider, i.e., dropbox or github.
- Token: Access token to storage provider's API.

Example file:

```text
provider: dropbox
token: mydropboxtoken
```

By default, Freighter assumes the presence of the configuration file in the /tmp directory; if it's located elsewhere it must be specified via the --config flag.

For docker, this file can be mounted as a secret _(only compatible with swarm)_.

<h2>Restore Files</h2>


With Docker:


```sh
docker service create -u [USER] -v [LOCAL_VOLUME_MAPPING]:/restore \
  --secret source=.freighter.yml,target=/tmp/.freighter.yml \
  opinari/freighter restore REMOTE_FILE_PATH /restore/LOCAL_FILE_PATH 
```

With native binary


```sh
./freighter restore FROM_REMOTE_PATH TO_LOCAL_PATH
```


<h2>Backup Files</h2>


With Docker:


```sh
docker service create -u [USER] -v [LOCAL_VOLUME_MAPPING]:/backup \
  --secret source=.freighter.yml,target=/tmp/.freighter.yml \
  opinari/freighter backup /backup/LOCAL_FILE_PATH REMOTE_FILE_PATH 
```

With native binary


```sh
./freighter backup FROM_LOCAL_PATH TO_REMOTE_PATH
```


<h2>Determine Age of File</h2>

To determine the age of a file (in days):


With Docker:


```sh
docker service create -u [USER] -v [LOCAL_VOLUME_MAPPING]:/age \
  --secret source=.freighter.yml,target=/tmp/.freighter.yml \
  opinari/freighter age REMOTE_FILE_PATH --ouput /age/OUTPUT_FILE_PATH 
```


With native binary

```sh
./freighter age REMOTE_FILE_PATH --output /tmp/age
```

<h2>Delete File(s)</h2>


With Docker:


```sh
docker service create -u [USER] \
  --secret source=.freighter.yml,target=/tmp/.freighter.yml \
  opinari/freighter delete REMOTE_FILE_PATH
```


With native binary

```sh
./freighter delete REMOTE_FILE_PATH [REMOTE_FILE_PATH...]
```



<hr>
