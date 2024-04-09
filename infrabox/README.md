# Terrabox

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.2-0f0f19.svg?style=flat-square)

A Terraform ecosystem module

## Features

###Terragrunt

```shell
Expose a terragrunt runtime

Usage:
  dagger call terragrunt [flags]
  dagger call terragrunt [flags] [command]

Flags:
      --image string     The image to use which contain terragrunt ecosystem (default "alpine/terragrunt")
      --version string   The version of the image to use (default "1.7.4")

Function Commands:
  apply                 Run an apply on a specific stack
  catalog               expose the module catalog (only available for terragrunt)
  container             Expose the container
  directory             Return the source directory
  disable-color         Indicate to disable the the color in the output
  do                    Execute the call chain
  format                Format the code
  output                Return the output of a specific stack
  plan                  Run a plan on a specific stack
  run-all               Execute the run-all command (only available for terragrunt)
  shell                 Open a shell
  with-cache-burster    Define the cache buster strategy
  with-container        Use a new container
  with-secret-dot-env   Convert a dotfile format to secret environment variables in the container (could be use to configure providers)
  with-source           Mount the source code at the given path
```

**more example in the `/ci/terrabox.go`**

## To Do

- [ ] Add support for terraform / opentofu
- [ ] Add support for terraform-docs
- [ ] Add support for [boilerplate](https://github.com/gruntwork-io/boilerplate)
- [ ] Add sops for secret management
