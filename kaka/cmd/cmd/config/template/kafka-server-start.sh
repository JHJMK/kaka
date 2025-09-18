#!/bin/bash

if [ $# -lt 1 ]; then
  echo "USAGE: $0 [-daemon] server.properties [--override property=value]*"
  exit 1
fi

base_dir=$(dirname $0)

# 日志
if [ "x$KAFKA_LOG4J_OPTS" = "x" ]; then
  export KAFKA_LOG4J_OPTS="-Dlog4j.configuration=file:$base_dir/../config/log4j.properties"
fi

# 堆大小
if [ "x$KAFKA_HEAP_OPTS" = "x" ]; then
  export KAFKA_HEAP_OPTS="-Xmx{{.MaxHeapSize}}G -Xms{{.MinHeapSize}}G"
fi


# jmx 端口
export JMX_PORT="9999"

EXTRA_ARGS=${EXTRA_ARGS-'-name kafkaServer -loggc'}

COMMAND=$1
case $COMMAND in
  -daemon)
    EXTRA_ARGS="-daemon "$EXTRA_ARGS
    shift
    ;;
  *)
    ;;
esac


exec $base_dir/kafka-run-class.sh $EXTRA_ARGS kafka.Kafka "$@"
