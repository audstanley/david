#!/bin/bash
# copy the config file to the root directory of the project (only need to do while testing things out)
# cp ./example/config.yaml ./config.yaml
rm david; cd cmd/david && go build . && mv ./david ../../david; cd ../.. && ./david -config=config.yaml