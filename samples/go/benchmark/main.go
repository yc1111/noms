package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/types"
)

func readfile(file string) ([]int, []string, []string) {
	fd, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", file, err))
	}
	scanner := bufio.NewScanner(fd)
	ops := make([]int, 0)
	keys := make([]string, 0)
	vals := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		op, _ := strconv.Atoi(line[0:1])
		ops = append(ops, op)
		line = line[2:]
		idx := strings.Index(line, "\t")
		keys = append(keys, line[0:idx])
		vals = append(vals, line[idx+1:])
	}
	return ops, keys, vals
}

func main() {
	// connect database
	cfg := config.NewResolver()
	db, ds, err := cfg.GetDataset("http://localhost:8000::benchmark")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create dataset: %s\n", err)
		return
	}
	defer db.Close()

	// init map
	_, initKeys, initVals := readfile("./data/init")
	initMap := types.NewMap(db).Edit()
	for i := range initKeys {
		initMap = initMap.Set(types.String(initKeys[i]), types.String(initVals[i]))
	}
	ds, _ = db.CommitValue(ds, initMap.Map())

	// execute workload
	ops, keys, vals := readfile("./data/input")
	if ops[0] == 1 {
		start := time.Now().UnixNano()
		for i := range ops {
			hv, _ := ds.MaybeHeadValue()
			currMap := hv.(types.Map).Edit()
			currMap.Set(types.String(keys[i]), types.String(vals[i]))
			ds, _ = db.CommitValue(ds, currMap.Map())
		}
		end := time.Now().UnixNano()
		fmt.Println(float64(len(ops)) / (float64(end-start) / 1000000))
	} else {
		start := time.Now().UnixNano()
		for i := 0; i < 10000; i++ {
			hv := ds.HeadValue()
			currMap := hv.(types.Map)
			currMap.Get(types.String(initKeys[i%len(initKeys)]))
		}
		end := time.Now().UnixNano()
		fmt.Println(float64(10000) / (float64(end-start) / 1000000))
	}
}
