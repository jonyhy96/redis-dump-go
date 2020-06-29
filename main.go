package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jonyhy96/redis-dump-go/redisdump"
)

func drawProgressBar(to io.Writer, currentPosition, nElements, widgetSize int) {
	if nElements == 0 {
		return
	}
	percent := currentPosition * 100 / nElements
	nBars := widgetSize * percent / 100

	bars := strings.Repeat("=", nBars)
	spaces := strings.Repeat(" ", widgetSize-nBars)
	fmt.Fprintf(to, "\r[%s%s] %3d%% [%d/%d]", bars, spaces, int(percent), currentPosition, nElements)
}

func realMain() int {
	var err error

	// TODO: Number of workers & TTL as parameters
	host := flag.String("host", "127.0.0.1", "Server host")
	port := flag.Int("port", 6379, "Server port")
	password := flag.String("password", "", "Password of redis server")
	keyRegex := flag.String("key", "*", "Keys regex to dump")
	nWorkers := flag.Int("n", 10, "Parallel workers")
	withTTL := flag.Bool("ttl", true, "Preserve Keys TTL")
	output := flag.String("output", "resp", "Output type - can be resp or commands")
	flag.Parse()

	var serializer func([]string) string
	switch *output {
	case "resp":
		serializer = redisdump.RESPSerializer

	case "commands":
		serializer = redisdump.RedisCmdSerializer

	default:
		log.Fatalf("Failed parsing parameter flag: can only be resp or json")
	}

	logger := log.New(os.Stdout, "", 0)
	if err = redisdump.DumpServer(*host, *port, *password, *keyRegex, *nWorkers, *withTTL, logger, serializer); err != nil {
		log.Fatal(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
