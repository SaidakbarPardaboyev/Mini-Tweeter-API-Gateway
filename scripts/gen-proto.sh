#!/bin/bash

# Get the current directory
CURRENT_DIR=$(pwd)

# Iterate through all directories inside Mini-Tweeter-Protocol-Buffer
for x in $(find ${CURRENT_DIR}/Mini-Tweeter-Protocol-Buffer/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/Mini-Tweeter-Protocol-Buffer -I /usr/local/include \
    --go_out=${CURRENT_DIR} \
    --go-grpc_out=${CURRENT_DIR} \
    --experimental_allow_proto3_optional \
    ${x}/*.proto
done

echo "Proto generation completed successfully."
