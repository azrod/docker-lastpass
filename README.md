[![Go Report Card](https://goreportcard.com/badge/github.com/azrod/docker-lastpass)](https://goreportcard.com/report/github.com/azrod/docker-lastpass)

# Sync Lastpass Secret to Docker Secret

Docker lastpass is a binary written in golang to synchronize your lastpass secrets with docker secret.

docker-lastpass use gret [lastpass-go](https://github.com/ansd/lastpass-go) library.

## Requierement

* Lastpass Account (Free or premium)
* [Docker Swarm](https://docs.docker.com/engine/swarm/) (Docker Secret is unavailable in docker standalone)

## Usage

```bash

./docker-lastpass --config config.toml --username <lastpass-email> --password <lastpass-password>
```

### Options

| Parameter    | Description          | Required           | Default       |
| ------------ | -------------------- | ------------------ | ------------- |
| `--config`   | Set config file path | :x:                | `config.toml` |
| `--username` | Lastpass Username    | :heavy_check_mark: |               |
| `--password` | Lastpass Password    | :heavy_check_mark: |               |
| `--otp`      | One Time Password    | :x:                |               |

### Configuration File

**config.toml example**

```toml
[log]
level = "debug" # debug,info,warn or error 

[lastpass]
twofactor = "push" # disable,push or OTP 

[secrets]
groups = ["docker"]
lists = []

```

**Secrets**

In groups add one or more "Folder" name in your lastpass. All secrets of each group will be synchronized.

In list add one or more secret "Name".


## Docker Secret

For each lastpass secret 2 docker secrets are created (Username and Password).

For example if your secret name in lastpass is `test secret` docker-lastpass create secret `lastpass_test-secret_Username` and `lastpass_test-secret_Password`

```bash
docker secret ls
ID               NAME                            DRIVER        CREATED              UPDATED
bhu3uuyl9nuxxx   lastpass_test-secret_Password                 xx days ago          xx days ago
jn9rqksbf00xxx   lastpass_test-secret_Username                 xx days ago          xx days ago
```

## Limitation

Lastpass API not provide timestamp for edit secret. it is therefore impossible to modify an existing secret. It is therefore to delete the secret in docker so that it can be recreated.


