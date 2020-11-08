# Contributing

Issues and Pull Requests are more than welcome!


## Running the code locally

> Prerequisites:
> - Go 1.15 (required)
> - Mage (optional, see install instructions at https://magefile.org/)




```bash
git clone https://github.com/benjlevesque/ghr
cd ghr
mage build
```

## Building

- with Mage
```bash
mage build
```

- without Mage
```bash
go build -o ghr .
```

NB: The Mage scripts does more than this command, but you probably don't need this for local testing.

## Tests
```bash
go test ./...
```