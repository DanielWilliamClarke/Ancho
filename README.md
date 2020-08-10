# Ancho

- [Ancho](#ancho)
  - [Build](#build)
  - [Test](#test)
  - [Coverage](#coverage)

## Build

``` bash
  cd src
  go build

  ./ancho.exe #windows
  ./ancho #linux
```

## Test

```bash
cd src
go test ./test/...
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
