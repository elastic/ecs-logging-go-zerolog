# Elastic Common Schema (ECS) support for zerolog

Use this library for automatically adding a minimal set of ECS fields to your logs, when using [zerolog](https://github.com/rs/zerolog).

## Usage

```go
package main

import (
	"os"

	"go.elastic.co/ecszerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := ecszerolog.New(os.Stdout)
	log.Logger = logger

	log.Info().Msg("Hello, World!")
}
```

## Test

```sh
go test ./...
```

## Contribute

Create a Pull Request from your own fork.

Run [mage](https://magefile.org/) to update and format you changes before submitting.

Add new dependencies to the NOTICE.txt.

## License

This software is licensed under the [Apache 2 license](./LICENSE).
