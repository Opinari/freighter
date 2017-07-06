# freighter

Freighter is a command-line utility to move files to/from storage backends in a convenient and efficient manner.

Use cases include:

- Restoring database backup files to a cloud compute instance
- Uploading local files to Dropbox for safe storage
- Downloading configuration files to a host from a VCS such as Github

## Currently Supported Storage Providers

- Dropbox (v2 API)
- Github (v3 API)

## Storage Provider Configuration
Freighter will firstly attempt to read a .freighter configuration file with any of the following supported extensions:
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

For docker, this file can be mounted as a secret _(only compatible with swarm)_ or a regular volume file mapping _(both docker run and swarm compatible)_

If the configuration file is not present or is missing keys, Freighter will fall back to the following environment variables:

- FREIGHTER_PROVIDER
- FREIGHTER_TOKEN

## Restore Files

With Docker Swarm and secret-mapped YML configuration file in non-default location:

```shell
docker service create -u [USER] \
  --mount type=bind,source=[LOCAL_VOLUME_MAPPING],destination=/restore \
  --secret freighter.yml \
  opinari/freighter --config /run/secrets/freighter.yml restore REMOTE_FILE_PATH /restore/LOCAL_FILE_PATH 
```

With Docker and environment variables:

```shell
docker run -u [USER] \
  -v [LOCAL_VOLUME_MAPPING]:/restore \
  -e FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] \
  -e FREIGHTER_TOKEN=[FREIGHTER_TOKEN] \
  opinari/freighter restore REMOTE_FILE_PATH /restore/LOCAL_FILE_PATH 
```

With native binary and configuration file in default location:

```shell
./freighter restore FROM_REMOTE_PATH TO_LOCAL_PATH
```

With native binary and environment variables:

```shell
env FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] FREIGHTER_TOKEN=[FREIGHTER_TOKEN] ./freighter restore FROM_REMOTE_PATH TO_LOCAL_PATH
```

## Backup Files

With Docker Swarm and secret-mapped YML configuration file in default location:

```shell
docker service create -u [USER] \
  --mount type=bind,source=[LOCAL_VOLUME_MAPPING],destination=/backup \
  --secret source=.freighter.yml,target=/tmp/.freighter.yml \
  opinari/freighter backup /backup/LOCAL_FILE_PATH REMOTE_FILE_PATH 
```

With Docker and environment variables:

```shell
docker run -u [USER] \
  -v [LOCAL_VOLUME_MAPPING]:/backup \
  -e FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] \
  -e FREIGHTER_TOKEN=[FREIGHTER_TOKEN] \
  opinari/freighter backup /backup/LOCAL_FILE_PATH REMOTE_FILE_PATH 
```

With native binary and TOML configuration file in non-default location:

```shell
./freighter --config /home/core/freighter.toml backup FROM_LOCAL_PATH TO_REMOTE_PATH
```

With native binary and environment variables (Dropbox used as an example):

```shell
env FREIGHTER_PROVIDER=dropbox FREIGHTER_TOKEN=my_dropbox_access_token ./freighter backup FROM_LOCAL_PATH TO_REMOTE_PATH
```

## Determine Age of File

With Docker Swarm, environment variable for provider and secret configuration file for token in non-default location:

```shell
docker service create -u [USER] \
  --env FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] \
  --mount type=bind,source=[LOCAL_VOLUME_MAPPING],destination=/age \
  --secret token.yml \
  opinari/freighter --config /run/secrets/token.yml age REMOTE_FILE_PATH --ouput /age/OUTPUT_FILE_PATH 
```

With native binary and configuration file in default location writing to output file:

```shell
./freighter age REMOTE_FILE_PATH --output /tmp/age
```

With native binary and environment variables writing to stdout:

```shell
env FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] FREIGHTER_TOKEN=[FREIGHTER_TOKEN] ./freighter age REMOTE_FILE_PATH
```

## Delete Files

With Docker Swarm and environment variables:

```shell
docker service create -u [USER] \
  --env FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] \
  --env FREIGHTER_TOKEN=[FREIGHTER_TOKEN] \
  opinari/freighter delete REMOTE_FILE_PATH
```

With native binary and configuration file in default location, deleting a single file:

```shell
./freighter delete REMOTE_FILE_PATH
```

With native binary and environment variables, deleting multiple files:

```shell
env FREIGHTER_PROVIDER=[FREIGHTER_PROVIDER] FREIGHTER_TOKEN=[FREIGHTER_TOKEN] ./freighter delete /remote/file1.tar.gz /remote/other/file2.tar.gz /remote/oldfile.tar.gz
```
