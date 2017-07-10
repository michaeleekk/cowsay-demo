#!/bin/bash

docker run -it --rm -v $PWD:/go/src/app proto_build $@

