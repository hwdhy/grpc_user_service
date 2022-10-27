#!/bin/sh

/go/app -port 50051 &

/go/app -type rest -port 8080 -endpoint 0.0.0.0:50051