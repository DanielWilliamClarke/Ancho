# Ancho

- [Ancho](#ancho)
  - [Build](#build)
  - [Test](#test)
  - [Coverage](#coverage)
  - [Cli Registration](#cli-registration)
  - [API Registration](#api-registration)

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

## Cli Registration

> Run using `./cli.exe`, The user can then interrupt the registration by hitting `ctrl+c`. The user will see that the registration will end gracefully and will recieve all DevEUIs that have been registered
>
> Users can also allow the registration to complete fully, without intervention

```bash
#...

Registration batch 2 complete!
Total DevEUIs registered: 54
Total DevEUIs registered: 56
Total DevEUIs registered: 61
Registration batch 0 complete!
Total DevEUIs registered: 65
Total DevEUIs registered: 72
Total DevEUIs registered: 76
Registration batch 7 complete!
Quit signal received, gracefully shutdown registration...
Total DevEUIs registered: 77
Shutdown obervation

2020/08/11 13:46:19 Shutdown registration batch: 9
2020/08/11 13:46:19 Shutdown registration batch: 3
2020/08/11 13:46:19 Shutdown registration batch: 8
2020/08/11 13:46:19 Shutdown registration batch: 1
2020/08/11 13:46:20 Shutdown registration batch: 4
2020/08/11 13:46:20 Shutdown registration batch: 5
Registration Complete!
Cleaning Up
9 DevEUIs failed ----------------
2020/08/11 13:46:20 DevEUI b9c428f66d63f317 already Registered: 422
2020/08/11 13:46:20 DevEUI 7d7a31914ff4d5df already Registered: 422
2020/08/11 13:46:20 DevEUI 9633a9d72220c9f9 already Registered: 422
2020/08/11 13:46:20 DevEUI 40aa03d8ff5fc29d already Registered: 422
2020/08/11 13:46:20 DevEUI 3bf6fbcb5ca4a5cb already Registered: 422
2020/08/11 13:46:20 DevEUI 27e51621b6e1ce1e already Registered: 422
2020/08/11 13:46:20 DevEUI d0561e44e6de6132 already Registered: 422
2020/08/11 13:46:20 DevEUI 1b8b9b3eb775405f already Registered: 422
2020/08/11 13:46:20 DevEUI 46c3403611c4de81 already Registered: 422
---------------------------------------------------
82 DevEUIs registered successfully ----------------
149ce23fdc6e93ab
da9d91ac9e199c6c
36f8da15c513fb80
7065eb2851f6bf04
2e4fb72a79f59746
5ecd1db467a50f1d
517f54a57b89c27c
72829540f1900994
968d00fdaa50d5e2
9c5ef8e46bf41eab
a71bba035baa994a
39baa4bc028be7e0

#...
```

## API Registration

> Run using `./api.exe`, The user will be shown the Fiber header on server start

```bash
        _______ __
  ____ / ____(_) /_  ___  _____   HOST     0.0.0.0  OS      WINDOWS
_____ / /_  / / __ \/ _ \/ ___/   PORT     3000     THREADS 4
  __ / __/ / / /_/ /  __/ /       TLS      FALSE    MEM     15.9G
    /_/   /_/_.___/\___/_/1.13.3  HANDLERS 2        PID     16016
```

> Run a request in a seperate terminal by running the following command

```bash
# A Idempotency-key header is required to facilitate Idempotency within the the API

curl -X PUT localhost:3000/v1/api/register -H 'Idempotency-key: your-key' &
#...
{"deveuis":["5d9f84fe12504efb","3264c0c85616186d","02799aedd50d203c", ... (up to 100 DevEUIs)]}

# Trying to run registration again while it is current running with the same Idempotency-key header will temporarily yield an empty payload
{"deveuis":[]}

# Once complete running the request with the same Idempotency-key header will always yield the same DevEUIs
curl -X PUT localhost:3000/v1/api/register -H 'Idempotency-key: your-key' &
#...
{"deveuis":["5d9f84fe12504efb","3264c0c85616186d","02799aedd50d203c", ... (up to 100 DevEUIs)]}

# Running a request without the Idempotency-key header will result in a 403 error
curl -X PUT localhost:3000/v1/api/register
#...
Idempotency-key invalid
```

`Why Ancho? Because it's been 34 degrees at 9pm`
