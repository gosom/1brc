 #!/bin/sh

GOOS=linux GOARCH=amd64 go build -o 1brc -ldflags="-s -w"  -gcflags="-l=4" main.go

/usr/bin/time -p -o time.txt ./1brc $1 > results.txt
