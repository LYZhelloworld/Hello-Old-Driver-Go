# Hello Old Driver (Go)
A LLSS Reader

## Build
Output folder: `./bin`

### Linux
```
make

# To clean output folder
make clean
```
### Windows
```
build
```

## Test
### Linux
```
make test
```
### Windows
```
test
```

## Usage
```
Usage: llss [-h] [-p protocol] [-d domain] [page]
  page: the page number of result. Must be greater than 0. Default value is 1.

Options:
  -d string
        Domain used to get the feed (auto checks domain if not given)
  -h    Show help
  -p string
        Protocol used to get the feed, like "http" or "https" (default "https")
```
