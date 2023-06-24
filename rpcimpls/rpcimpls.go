package rpcimpls

import (
    "net"
)

// this enables the ports which pods are listening on to be visible from one lcoation
func CreatePodListener() (net.Listener, error) {
    listener, err := net.Listen("tcp", ":8080")
	if err != nil {
        return nil, err
	}

    return listener, nil
}


