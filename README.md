# reedlib

A high-level [Reed Solomon](https://en.wikipedia.org/wiki/Reed%E2%80%93Solomon_error_correction)
erasure encoding library and set of command-line tools written in [Go](https://golang.org).

## Quick Start

### Installing the library

Either run:
```#!console
go get github.com/prologic/reedlib
```

Or just import the library in your project:
```#!go
import "github.com/prologic/reedlib"
```

### Installing the command-line tools

Install `reed-encode`:
```#!console
go get github.com/prologic/reedlib/cmd/reed-encode/....

Install `reed-decode`:
```#!console
go get github.com/prologic/reedlib/cmd/reed-decode/....
```

## Usage

### Encoding a file (command-line)

To encode a file using Reed-Solomon erasure encoding using the default Data and
Parity Shards (_3 + 1 respectively_):

```#!console
reed-encode ./testdata/IMG_7895.JPG
```

This will result in a number of output files from the result of splitting up
and encoding the original input file using the specified number of data and
parity shards:

```#!console
$ ls -lah ./testdata/
total 6.4M
drwxr-xr-x  8 prologic staff  256 Dec 27 09:10 .
drwxr-xr-x 17 prologic staff  544 Dec 27 09:12 ..
-rw-r--r--  1 prologic staff 6.1K Dec 27 09:10 .DS_Store
-rw-r--r--  1 prologic staff 2.7M Dec 27 09:10 IMG_7895.JPG
-rw-r--r--  1 prologic staff 916K Dec 27 09:10 IMG_7895.JPG.0
-rw-r--r--  1 prologic staff 916K Dec 27 08:55 IMG_7895.JPG.1
-rw-r--r--  1 prologic staff 916K Dec 27 08:55 IMG_7895.JPG.2
-rw-r--r--  1 prologic staff 916K Dec 27 08:55 IMG_7895.JPG.3
```

Keep the pieces (_data and parity shards_) on different storage devices or
storage nodes which can then later be used to reconstruct the original file,
even if one of the shards is lost or corrupt (_parity of one_).

For increased data saftey and redundancy, you can increase the partity shards.

For increased network concurrnecy you can increase the data shards.

### Decoding shards (command-line)

To decode a number of shards from a previous encoding using the default Data
and Parity Shards (_3 + 1 respectively_):

First remove the original input file to demonstrate recovery:

```#!console
rm -f ./testdata/IMG_7895.JPG
```

Now reconstruct the original input file using the shards:

```#!console
reed-decode ./testdata/IMG_7895.JPG
```

You should now have the original file recovered and intact:

```#!console
$ ls -lah ./testdata/
total 6.3M
drwxr-xr-x  7 prologic staff  224 Dec 27 09:33 .
drwxr-xr-x 15 prologic staff  480 Dec 27 09:33 ..
-rw-r--r--  1 prologic staff 2.7M Dec 27 09:33 IMG_7895.JPG
-rw-r--r--  1 prologic staff 916K Dec 27 09:33 IMG_7895.JPG.0
-rw-r--r--  1 prologic staff 916K Dec 27 09:33 IMG_7895.JPG.1
-rw-r--r--  1 prologic staff 916K Dec 27 09:33 IMG_7895.JPG.2
-rw-r--r--  1 prologic staff 916K Dec 27 09:33 IMG_7895.JPG.3
```

You can even remove the original input file and either remove or corrupt one
of the shards and the original input file is still recoverable from the
remaining shards.

## License

`reedlib` is licensed under the terms of the [MIT License](/LICENSE)
