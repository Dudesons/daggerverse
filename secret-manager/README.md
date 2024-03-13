# Secret manager

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A secret manager module which allow to work with different backends

## Features

* Gcp secret manager
  * Read secret
  * Create/update secret


## Examples

```shell
dagger call -m github.com/Dudesons/daggerverse/secret-manager \
  gcp \
  get-secret --project=my-gcp-project-id --name=MY_SECRET_KEY --gcloud-folder="$HOME/.config/gcloud/"
  plaintext
```

## To Do

- [ ] Add aws secret manager
- [ ] Add vault
- [ ] Add sops
- [ ] Improve documentation

