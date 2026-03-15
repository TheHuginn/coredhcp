package etcd

import (
	"context"
	"errors"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/insomniacslk/dhcp/dhcpv4"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var Plugin = plugins.Plugin{
	Name:   "etcd",
	Setup4: setup4,
}

var log = logger.GetLogger("plugins/etcd")

// Etcd connection
var etcd *clientv3.Client = nil

func createClient(endpoints []string) *clientv3.Client {
	if len(endpoints) < 3 {
		log.Error("Invalid etcd endpoints")
		return nil
	}

	log.Println("Creating client with provided endpoints...")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})

	if err != nil {
		log.Error("Error creating etcd client:", err)
		return nil
	}

	healthy := 0

	for endpoint := range cli.Endpoints() {

		//We need at least 2 endpoints

		log.Println("Checking status for", cli.Endpoints()[endpoint])
		_, err = cli.Status(context.Background(), cli.Endpoints()[endpoint])
		if err != nil {
			log.Error("Error getting etcd status:", err)
			break
		}

		log.Println("Endpoint is healthy.")
		healthy++
	}

	if healthy < 2 {
		log.Error("Etcd is in unhealthy state...")
		return nil
	}

	return cli
}

func Handler4(req *dhcpv4.DHCPv4, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool) {

	log.Println("Req on ", req.ServerIPAddr, "\n", req.Summary())

	return resp, false
}

func setup4(args ...string) (handler.Handler4, error) {

	//Create a client
	etcd = createClient(args[:3])
	if etcd == nil {
		return nil, errors.New("Error creating etcd client")
	}

	//Elect master, or become passive

	return Handler4, nil
}
