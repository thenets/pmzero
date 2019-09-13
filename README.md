# pmzero

[![GoDoc](https://godoc.org/github.com/thenets/pmzero?status.svg)](https://godoc.org/github.com/thenets/pmzero)
[![Go Report Card](https://goreportcard.com/badge/github.com/thenets/pmzero)](https://goreportcard.com/report/github.com/thenets/pmzero)
[![GitHub Actions](https://github.com/thenets/pmzero/workflows/build/badge.svg)](https://github.com/thenets/pmzero/actions)

Inspired by PM2 and a hard Amazon interview process.

## Requirements

- Linux x64
- Windows x64
- MacOS X x64

## Installation

### Linux

You need to add the binary file to the `/usr/bin` and make it executable. You can use the following script to deploy.

```bash
curl -sSL https://raw.githubusercontent.com/thenets/pmzero/master/docs/install-linux.sh | sudo sh
```

## Docs

```
NAME:
   pmzero - Easiest multi-platform process manager

USAGE:
   pmzero.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.2-alpha

COMMANDS:
   load, l     load a config file
   start, c    start a deployment
   list        list all deployments
   delete      delete a deployment
   stop        stop a deployment
   logs        follow the logs files from a deployment
   foreground  keep running and respawn all deployments
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```