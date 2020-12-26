package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauspost/reedsolomon"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/prologic/reedlib"
)

var (
	debug   bool
	version bool

	// Basic options
	dataShards     int
	parityShards   int
	outputFilename string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] basefile.ext\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNOTE: Do not add the number to the filename.\n")
		flag.PrintDefaults()
	}

	flag.BoolVarP(&debug, "debug", "D", false, "enable debug logging")
	flag.BoolVarP(&version, "version", "v", false, "display version information")

	// Basic options
	flag.IntVarP(&dataShards, "data", "d", 3, "no. of data shards")
	flag.IntVarP(&parityShards, "parity", "p", 1, "no. of parity shards")
	flag.StringVarP(&outputFilename, "output", "o", ".", "output filename")
}

func flagNameFromEnvironmentName(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "_", "-", -1)
	return s
}

func parseArgs() error {
	for _, v := range os.Environ() {
		vals := strings.SplitN(v, "=", 2)
		flagName := flagNameFromEnvironmentName(vals[0])
		fn := flag.CommandLine.Lookup(flagName)
		if fn == nil || fn.Changed {
			continue
		}
		if err := fn.Value.Set(vals[1]); err != nil {
			return err
		}
	}
	flag.Parse()
	return nil
}

func main() {
	parseArgs()

	if version {
		fmt.Printf("reedsolomon v%s", reedlib.FullVersion())
		os.Exit(0)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: No filename given\n")
		flag.Usage()
		os.Exit(1)
	}
	fn := flag.Arg(0)

	// Create matrix
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	checkErr(err)

	// Open the inputs
	shards, size, err := openInput(dataShards, parityShards, fn)
	checkErr(err)

	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		log.Debugf("No reconstruction needed")
	} else {
		log.Warnf("Verification failed. Reconstructing data")
		shards, size, err = openInput(dataShards, parityShards, fn)
		checkErr(err)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				outfn := fmt.Sprintf("%s.%d", fn, i)
				log.Debugf("Creating %s", outfn)
				out[i], err = os.Create(outfn)
				checkErr(err)
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			checkErr(fmt.Errorf("reconstruction failed: %w", err))
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				checkErr(err)
			}
		}
		shards, size, err = openInput(dataShards, parityShards, fn)
		ok, err = enc.Verify(shards)
		if !ok {
			checkErr(fmt.Errorf("verification failed after reconstruction, data likely corrupt: %w", err))
		}
		checkErr(err)
	}

	// Join the shards and write them
	outfn := outputFilename
	if outfn == "" {
		outfn = fn
	}

	log.Debugf("Writing data to %s", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	shards, size, err = openInput(dataShards, parityShards, fn)
	checkErr(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, int64(dataShards)*size)
	checkErr(err)
}

func openInput(dataShards, parityShards int, fn string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parityShards)
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fn, i)
		log.Debugf("Opening %s", infn)
		f, err := os.Open(infn)
		if err != nil {
			log.WithError(err).Warnf("Error reading file %s", infn)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		checkErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
