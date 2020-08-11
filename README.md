# Ancho

- [Ancho](#ancho)
  - [Build](#build)
  - [Test](#test)
  - [Coverage](#coverage)

## Build

``` bash
  cd src

  # Build CLI
  go build ./cli

  ./cli.exe #windows
  ./cli #linux

  # Build API
  go build ./api

  ./api.exe #windows
  ./api #linu
```

## Test

```bash
cd src
go test -v ./test/... # -v for verbose
```

## Coverage

```bash
mkdir test_results
cd src
../scripts/generate_code_coverage.sh

# Upon completion you can then access the coverage artifacts in test_results
# if on windows
start chrome ../test_results/index.html
```

`Why Ancho? Because it's 34 degrees at 9pm`
