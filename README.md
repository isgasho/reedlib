# reedlib

A high-level [Reed Solomon](https://en.wikipedia.org/wiki/Reed%E2%80%93Solomon_error_correction)
erasure encoding library and set of command-line tools written in [Go](https://golang.org).

## Quick Start

Installing the command-line tools:

```#!console
go get github.com/prologic/reedlib/cmd/reed-encode/....
go get github.com/prologic/reedlib/cmd/reed-decode/....
```

Usage:

```#!
$ ./reed-encode -h
Usage: ./reed-encode [options]
  -d, --data int        no. of data shards (default 3)
  -D, --debug           enable debug logging
  -o, --output string   output directory (default ".")
  -p, --parity int      no. of parity shards (default 1)
  -v, --version         display version information
pflag: help requested

$ ./reed-decode -h
Usage: ./reed-decode [options] basefile.ext

NOTE: Do not add the number to the filename.

  -d, --data int        no. of data shards (default 3)
  -D, --debug           enable debug logging
  -o, --output string   output filename
  -p, --parity int      no. of parity shards (default 1)
  -v, --version         display version information
pflag: help requested
```

## License

`reedlib` is licensed under the terms of the [MIT License](/LICENSE)
