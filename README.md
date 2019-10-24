# Hello Old Driver (Go)
A LLSS Reader

Inspired by [Chion82/hello-old-driver](https://github.com/Chion82/hello-old-driver).

## Test
* Windows: `test.cmd`
* Linux: `./test.sh`

## Build
* Windows: `build.cmd`
* Linux: `./build.sh`

The executable file will be in `bin` folder.

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