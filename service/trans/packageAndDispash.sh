CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o transService main.go
scp transService vagrant@192.168.99.40:/opt/service/
scp transService vagrant@192.168.99.50:/opt/service/
scp transService vagrant@192.168.99.60:/opt/service/
rm ./transService
