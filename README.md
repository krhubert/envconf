# Envconf

[![Build Status](https://travis-ci.org/krhubert/envconf.png)](https://travis-ci.org/krhubert/envconf)
[![Coverage Status](https://coveralls.io/repos/github/krhubert/envconf/badge.svg)](https://coveralls.io/github/krhubert/envconf)
[![Go Report Card](http://goreportcard.com/badge/krhubert/envconf)](http://goreportcard.com/report/krhubert/envconf)
[![GoDoc](https://godoc.org/github.com/krhubert/envconf?status.svg)](https://godoc.org/github.com/krhubert/envconf)


Package envconf implements reading config from enviroment variables.

```Go
import "github.com/krhubert/envconfig"
```

## Usage

```Go
package main

type Config struct {
    Host string `envconf:"name,default,required"`
}
```

Tag:
- name - the name of env variable (if prefix is set, it will be added to name -> prefix_name)
- default - if no env variable found, the default value will be used
- required - can be set to true, then user must set this variable (required can't be combined with default)

## Example

```Go
package main

import (
    "log"

    "github.com/krhubert/envconf"
)

type Config struct {
    Host    string `envconf:"host,,true"`
    Port    int    `envconf:"port,80"`
    Timeout int    `envconf:"timeout"`
}
```

Basic usage

```Go
func main() {
    var conf Config
    if err := envconf.SetValues(&conf); err != nil {
        log.Fatal(err)
    }
    log.Printf("host=%s, port=%d, timeout=%d\n", conf.Host, conf.Port, conf.Timeout)
}
```

```Bash
$ HOST=127.0.0.1 PORT=8080 TIMEOUT=60 ./main
host=127.0.0.1, port=8080, timeout=60
$
```

Application prefix can be set by creating config

```Go
func main() {
    var conf Config
    if err := envconf.NewConfig("app").SetValues(&conf); err != nil {
        log.Fatal(err)
    }
    log.Printf("host=%s, port=%d, timeout=%d\n", conf.Host, conf.Port, conf.Timeout)
}
```

```Bash
$ APP_HOST=127.0.0.1 APP_PORT=8080 APP_TIMEOUT=60 ./main
host=127.0.0.1, port=8080, timeout=60
$
```
