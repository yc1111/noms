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
  db, ds, err := cfg.GetDataset("mem::benchmark")
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not create dataset: %s\n", err)
    return
  }
  defer db.Close()

  // read init value
  _, init_keys, init_vals := readfile("/data/yc/USTORE/data/siri/init")

  /*********************
   * struct
   *********************/
//   init_struct := types.NewStruct("S1", types.StructData{})
//   for i,_ := range init_keys {
//     init_struct.Set("s" + init_keys[i], types.String(init_vals[i]))
//   }
//   ds, _ = db.CommitValue(ds, init_struct)

  /*********************
   * map
   *********************/
  init_map := types.NewMap(db);
  for i,_ := range init_keys {
    init_map.Edit().Set(types.String(init_keys[i]), types.String(init_vals[i])).Map()
  }
  ds, _ = db.CommitValue(ds, init_map)

  // perform operations
  ops, keys, vals := readfile("/data/yc/USTORE/data/siri/input")
  start := time.Now().UnixNano()
  for i, _ := range ops {
    /*********************
     * struct
     *********************/
//     ds, _ = db.CommitValue(ds, types.NewStruct("S1", types.StructData{"s" + keys[i]: types.String(vals[i])}))
    /*********************
     * map
     *********************/
    ds, _ = db.CommitValue(ds, types.NewMap(db, types.String(keys[i]), types.String(vals[i])))
  }
  end := time.Now().UnixNano()
  fmt.Println(float64(len(ops)) / (float64(end - start)/1000000))
}
