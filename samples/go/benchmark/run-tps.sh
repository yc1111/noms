#!/bin/bash
SIZE=(10000 20000 40000 80000 160000 320000 640000 1280000)
THETA=(0)
W_RATIO=(1)
BENCHMARK_DIR=/data/yc/USTORE/test/siri
DATA_DIR=/data/yc/USTORE/data/siri

go build main.go
mkdir -p ./data
rm -f ./data/*

for j in "${THETA[@]}"
do
  for k in "${W_RATIO[@]}"
  do
    echo
    echo "$j $k"
    echo "----------------------------------"
    echo -e "\"#Records\"\t\"Throughput\"" > ./data/tps_${j}_${k}
    for i in "${SIZE[@]}"
    do
      echo $i

      let init_size=$i/2
      cp $BENCHMARK_DIR/dataset/input_${init_size}_${j}_1 $DATA_DIR/init
      let workload=$i
      cp $BENCHMARK_DIR/dataset/input_${workload}_${j}_${k} $DATA_DIR/input

      echo -en "$i\t" >> ./data/tps_${j}_${k}
      ./main >> ./data/tps_${j}_${k}
    done
  done
done

