# go-composite-logger

# Description

**go-composite-logger** is a simple library for creating a composite logger.

# Installation

```
go get github.com/spatialcurrent/go-composite-logger
```

# Usage

**Import**

```
import (
  "github.com/spatialcurrent/go-composite-logger/compositelogger"
)
```

**Config**

```
logs {

  log {
    location = "stdout"
    level = "info"
    format = "text"
  }

  log {
    location = "~/myapplication/debug.jsonl"
    level = "warn"
    format = "json"
  }

}

```

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-composite-logger/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
