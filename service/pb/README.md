# Reference

http://www.grpc.io/docs/quickstart/go.html

http://www.grpc.io/docs/guides

http://www.grpc.io/docs/guides/concepts.html

http://www.grpc.io/docs/tutorials/basic/go.html

https://developers.google.com/protocol-buffers/docs/overview

https://github.com/google/protobuf/releases

https://segmentfault.com/a/1190000008672912 [http://www.tui8.com/articles/popular/90457.html ]

https://blog.1024coder.com/14766109567964.html

https://github.com/grpc/grpc/tree/master/doc

https://coreos.com/etcd/docs/latest/dev-guide/grpc_naming.html

# Environment
	localhost: OS X 10.10.3

# Installation

## gRPC

	$ git clone https://github.com/golang/text.git
	$ go get google.golang.org/grpc

## protocol buffers v3

	download protoc-3.3.0-osx-x86_64.zip from https://github.com/google/protobuf/releases
	unzip it into /usr/local/protoc

	$ echo  "PATH=$PATH:/usr/local/protoc/bin"  >>  ~/.bash_profile
	$ source  ~/.bash_profile

## protoc plugin for Go

	$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
