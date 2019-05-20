#!/bin/bash
SHARE_HOME=/users/yc

SIZE=(1280000 640000 320000 160000 80000 40000 20000 10000)
THETA=(0)
W_RATIO=(0)
BENCHMARK_DIR=$SHARE_HOME/USTORE/test/siri
DATA_DIR=$SHARE_HOME/go/src/github.com/attic-labs/noms/samples/go/benchmark/data

go build main.go
mkdir -p $DATA_DIR
rm -f $DATA_DIR/*

for j in "${THETA[@]}"
do
  for k in "${W_RATIO[@]}"
  do
    echo
    echo "$j $k"
    echo "----------------------------------"
    echo -e "\"#Records\"\t\"Throughput\"" > $DATA_DIR/tps_${j}_${k}
    for i in "${SIZE[@]}"
    do
      echo $i

      let init_size=$i/2
      cp $BENCHMARK_DIR/dataset/input_${init_size}_${j}_1 $DATA_DIR/init
      let workload=$i
      cp $BENCHMARK_DIR/dataset/input_${workload}_${j}_${k} $DATA_DIR/input

      noms serve mem &
      echo -en "$i\t" >> ./data/tps_${j}_${k}
      ./main >> ./data/tps_${j}_${k}
      for i in `echo $!`
      do
        kill -INT $i
      done
    done
  done
done

