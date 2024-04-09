# Drift

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.2-0f0f19.svg?style=flat-square)

A drift detection module for terraform / terragrunt

## Features

```shell
Flags:
      --focus                Only show output for focused commands (default true)
      --json                 Present result as JSON
  -m, --mod string           Path to dagger.json config file for the module or a directory containing that file. Either local path (e.g. "/path/to/some/dir") or a github repo (e.g. "github.com/dagger/dagger/path/to/some/subdir")
      --mount-point string   Define where the code is mounted, this could impact for absolute module path (default "/terraform")
  -o, --output string        Path in the host to save the result to

Function Commands:
  detection         Trigger the drift detection
  report-to-slack   Send the report formated to slack
```


## Examples

```shell
dagger call \
  detection \
  --src=../testdata/infrabox/terraform/ \
  --stack-root-path=stacks/dev/europe-west1/staging \
  --max-parallelization=0 \
  report-to-slack \
  --token=env:MY_SLACK_TOKEN \
  --channel-id=C06NZPY0FM2
```

## To Do

- [ ] Add support for terraform
- [ ] Add support for custom templates
