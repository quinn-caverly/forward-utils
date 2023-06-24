package rpcimpls

import (
    "log"
    "net"
    "net/rpc"

    "utils/endpointstructs"
)

func CallToGoReadPod() {

}

func CreateGoReadPodListener() (net.Listener, error) {
    listener, err := net.Listen("tcp", ":8080")
	if err != nil {
        return nil, err
	}

    return listener, nil
}
