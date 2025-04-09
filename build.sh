#!/bin/bash

cd front
bun run build
cd ..
go build
