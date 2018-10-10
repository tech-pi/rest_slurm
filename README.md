# rest_slurm
Restful API for Slurm

## compile for linux
1. `set GOOS=linux`
2. `set GOARCH=amd64`
3. `go build -o rest-slurm .\main.go .\handler.go`