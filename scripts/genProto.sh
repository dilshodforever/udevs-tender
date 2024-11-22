#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/gen_proto
for x in $(find ${CURRENT_DIR}/protos* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR} -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done