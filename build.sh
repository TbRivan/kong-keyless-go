docker run -it -v "/Users/inovasisolusi/go/src/kong-keyless:/go/src/kong-keyless"  golang:1.20-alpine3.18 sh -c " \
cd src/kong-keyless && \
go build -o ./build/kong-keyless main.go"