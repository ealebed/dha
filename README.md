# dha

Docker Hub Administrator utility

## Get

```bash
git clone git@github.com:ealebed/dha.git
```

```bash
cd dha
```

```bash
go install github.com/ealebed/dha
```

## Set your dockerhub login/password as env-variables
```bash
export DOCKERHUB_USERNAME=
export DOCKERHUB_PASSWORD=
```

## Use

```bash
dha -h
```

---

## Syntax

Use the following syntax to run `dha` commands from your terminal window:

```bash
dha [command] [flags]
```

### Flags are

| flag | Description |
| ----------- | ------------ |
| `--dry-run` | bool; print output only (default true) |
| `--org` | string; source owner user/organization (default "DOCKERHUB_USERNAME") |
| `--version` | dha version |

### Commands are

| command | Description |
| ----------- | ------------ |
| `delete`, `del` | delete the specified dockerhub repository |
| `describe` | returns information about the specified dockerhub repository |
| `get` | returns list tags from the specified dockerhub repository |
| `list`, `ls` | returns list of all dockeruhub repositories |
| `truncate` | truncate tags in the specified docker image repository |
| `help` | help about any command |

### Manage Docker images

```bash
# Delete the specified docker image repository from DockerHub.
dha delete --image=airflow --dry-run=false

# List all image repositories (and count tags in square brackets) from DockerHub.
dha list

# Get detailed information about the specified docker image repository on DockerHub.
dha describe --image=airflow

# Get tags from the specified docker image repository on DockerHub.
dha get --image=airflow

# Truncate old tags (that are older than 30 days except latest 25 ones) in the specified docker image repository on DockerHub.
dha truncate --image=airflow --truncateOld=true --dry-run=false

# Truncate tags in the specified docker image repository on DockerHub by regEx.
dha truncate --image=airflow --regEx=dev --dry-run=false

# Truncate tags in all docker image repositories on DockerHub by regEx.
dha truncate --all --regEx=cron --dry-run=false

# Renew (pull/push) tags in the specified docker image repository on DockerHub.
dha renew --image=sentinel-dashboard --dry-run=false

# Renew (pull/push) tags in all organization repositories on DockerHub.
dha renew --all --dry-run=false
```
