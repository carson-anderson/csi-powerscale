package csi_utils

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

// GetNFSClientIP is used to fetch IP address from networks on which NFS traffic is allowed
func GetNFSClientIP(allowedNetworks []string) (string, error) {
	var nodeIP string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("Encountered error while fetching system IP addresses: %+v\n", err.Error())
		return "", err
	}

	// Populate map to optimize the algorithm for O(n)
	networks := make(map[string]bool)
	for _, cnet := range allowedNetworks {
		networks[cnet] = false
	}

	for _, a := range addrs {
		switch v := a.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil {
				ip, cnet, err := net.ParseCIDR(a.String())
				log.Debugf("IP address: %s and Network: %s", ip, cnet)
				if err != nil {
					log.Errorf("Encountered error while parsing IP address %v", a)
					continue
				}

				if _, ok := networks[cnet.String()]; ok {
					log.Infof("Found IP address: %s", ip)
					nodeIP = ip.String()
					return nodeIP, nil
				}
			}
		}
	}

	// If a valid IP address matching allowedNetworks is not found return error
	if nodeIP == "" {
		return "", errors.New(fmt.Sprintf("No valid IP address found matching against allowedNetworks %v", allowedNetworks))
	}

	return nodeIP, nil
}
