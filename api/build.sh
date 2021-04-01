#!/usr/bin/env bash

if [ ! -d go ]; then
  echo "[!] Creating go output directory"
  mkdir go;
fi

echo "[+] Building docker container"
docker image build -t go-iperf-builder:1.0 .
docker container run --detach --name builder go-iperf-builder:1.0
docker cp grpc:/go/src/github.com/BGrewell/go-iperf/api/go/api.pb.go go/.
echo "[+] Updating of go library complete"

echo "[+] Removing docker container"
docker rm builder

echo "[+] Adding new files to source control"
git add go/api.pb.go
git commit -m "regenerated grpc libraries"
git push

echo "[+] Done. Everything has been rebuilt and the repository has been updated and pushed"