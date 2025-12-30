package libvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func ListDomains(user, keyPath, ip string) error {
	// e.g. qemu+ssh://alan@10.0.129.11/system?keyfile=/home/alan/.ssh/id_rsa&no_verify=1
	uri := fmt.Sprintf("qemu+ssh://%s@%s/system?keyfile=%s&no_verify=1", user, ip, keyPath)

	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return fmt.Errorf("failed to list domains: %v", err)
	}

	fmt.Printf("%d VMs found:\n", len(domains))
	for _, d := range domains {
		name, _ := d.GetName()
		id, _ := d.GetID()
		state, _, _ := d.GetState()
		stateStr := "Unknown"
		switch state {
		case libvirt.DOMAIN_RUNNING:
			stateStr = "Running"
		case libvirt.DOMAIN_SHUTOFF:
			stateStr = "Shutoff"
		case libvirt.DOMAIN_PAUSED:
			stateStr = "Paused"
		}
		fmt.Printf(" - %s (ID: %d, State: %s)\n", name, id, stateStr)
		d.Free()
	}

	return nil
}
