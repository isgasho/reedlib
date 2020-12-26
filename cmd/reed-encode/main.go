package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	dataShards      int
	parityShards    int
	outputDirectory string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVarP(&debug, "debug", "D", false, "enable debug logging")
	flag.BoolVarP(&version, "version", "v", false, "display version information")

	// Basic options
	flag.IntVarP(&dataShards, "data", "d", 3, "no. of data shards")
	flag.IntVarP(&parityShards, "parity", "p", 1, "no. of parity shards")
	flag.StringVarP(&outputDirectory, "output", "o", ".", "output directory")
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

	if (dataShards + parityShards) > 256 {
		fmt.Fprintf(os.Stderr, "Error: sum of data and parity shards cannot exceed 256\n")
		os.Exit(1)
	}

	fn := flag.Arg(0)

	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	checkErr(err)

	log.Debugf("Opening %s", fn)
	f, err := os.Open(fn)
	checkErr(err)

	instat, err := f.Stat()
	checkErr(err)

	shards := dataShards + parityShards
	out := make([]*os.File, shards)

	// Create the resulting files.
	dir, file := filepath.Split(fn)
	dir = filepath.Join(outputDirectory, dir)
	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		log.Debugf("Creating %s", outfn)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		checkErr(err)
	}

	// Split into files.
	data := make([]io.Writer, dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, instat.Size())
	checkErr(err)

	// Close and re-open the files.
	input := make([]io.Reader, dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		checkErr(err)
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, parityShards)
	for i := range parity {
		parity[i] = out[dataShards+i]
		defer out[dataShards+i].Close()
	}

	// Encode parity
	err = enc.Encode(input, parity)
	checkErr(err)
	log.Debugf("File split into %d data + %d parity shards.\n", dataShards, parityShards)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
