# kubenv

[![Actions Status](https://github.com/MqllR/kubenv/workflows/Build%20and%20release/badge.svg)](https://github.com/MqllR/kubenv/actions)

WIP WIP

## Install

```
go build
mv kubenv /usr/local/bin/
```

Define your configuration file as kubenv-example.yaml and export the environment variable:

```
export KUBENV_CONFIG=/path/to/my/config.yaml
```

## Tests

```
go test -v ./...
```
