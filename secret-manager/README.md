# Secret manager

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)

A secret manager module which allow to work with different backends

## Features

* Gcp secret manager
  * Read secret
  * Create/update secret
* Aws secret manager
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

- [ ] Add vault
- [ ] Add sops
- [ ] Improve documentation

