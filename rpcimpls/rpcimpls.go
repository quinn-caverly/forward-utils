package rpcimpls

import (
	"net"
	"net/rpc"
)

// this enables the ports which pods are listening on to be visible from one lcoation
func CreatePodListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func ConnectToGoWritePod() (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", "go-write-service:8080")
	if err != nil {
        return nil, err
	}

    return client, err
}
