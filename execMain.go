package main

import (
	"errors"
	libraryLogging "github.com/s-r-engineer/library/logging"
)

func execWGdown(interfaceName, interfaceIP, defaultRouteTable string) (err error) {
	libraryLogging.Debug("Checking default route")
	if err = checkDefaultRoute(interfaceName, defaultRouteTable); err == nil || errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Deleting default route")
		err = deleteDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			libraryLogging.Debug("Error deleting default route")
			return err
		}
	}
	libraryLogging.Debug("Default route clean")
	libraryLogging.Debug("Deleting server rule")
	err = deleteServerRule()
	if err != nil {
		libraryLogging.Debug("Error deleting server rule")
		return err
	}
	libraryLogging.Debug("Server rule clean")
	libraryLogging.Debug("Deleting lookup rule")
	err = deleteLookupRule(defaultRouteTable)
	if err != nil {
		libraryLogging.Debug("Error deleting lookup rule")
		return err
	}
	libraryLogging.Debug("Lookup rule clean")
	libraryLogging.Debug("Setting link down")
	err = linkDown(interfaceName)
	if err != nil {
		libraryLogging.Debug("Error setting link down")
		return err
	}
	libraryLogging.Debug("Link clean")
	libraryLogging.Debug("Deleting IP address")
	err = deleteIpAddress(interfaceIP, interfaceName)
	if err != nil {
		libraryLogging.Debug("Error deleting IP address")
		return err
	}
	libraryLogging.Debug("IP address clean")
	libraryLogging.Debug("Deleting interface")
	err = deleteInterface(interfaceName)
	if err != nil {
		libraryLogging.Debug("Error deleting interface")
		return err
	}
	libraryLogging.Debug("Interface clean")
	return nil
}

func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP, defaultWGPort, defaultRouteTable string) (err error) {
	libraryLogging.Debug("Checking the interface")
	if !checkInterface(interfaceName) {
		libraryLogging.Debug("Adding interface")
		err = addInterface(interfaceName)
		if err != nil {
			libraryLogging.Debug("Error adding interface")
			return err
		}
	}
	libraryLogging.Debug("Interface OK")
	//libraryLogging.Debug("Setting DNS")
	//err = setDNS(interfaceName, defaultNordVPNDNS)
	//if err != nil {
	//	libraryLogging.Debug("Error setting DNS")
	//	return err
	//}
	//libraryLogging.Debug("Setting DNS ok")
	libraryLogging.Debug("Checking private key")
	err = checkPrivateKey(interfaceName, privateKey)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Setting private key")
		err = setPrivateKey(interfaceName, privateKey)
		if err != nil {
			libraryLogging.Debug("Error setting private key")
			return err
		}
	} else if err != nil {
		libraryLogging.Debug("Error checking private key")
		return err
	}
	libraryLogging.Debug("Private key OK")
	libraryLogging.Debug("Checking peer")
	err = checkIfPeerOk(interfaceName, publicKey, endpointIP, defaultWGPort)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Deleting peers")
		err = deletePeers(interfaceName)
		if err != nil {
			libraryLogging.Debug("Error deleting peers")
			return err
		}
		libraryLogging.Debug("Deleting peers done")
		libraryLogging.Debug("Setting peer")
		err = setPeer(interfaceName, publicKey, endpointIP, defaultWGPort)
		if err != nil {
			libraryLogging.Debug("Error setting peer")
			return err
		}
		libraryLogging.Debug("Setting peer done")
	} else if err != nil {
		libraryLogging.Debug("Error checking peer")
		return err
	}
	libraryLogging.Debug("Checking peer ok")
	libraryLogging.Debug("Checking address")
	err = checkIfAddressOk(interfaceName, interfaceIP)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Setting address")
		err = setAddress(interfaceName, interfaceIP)
		if err != nil {
			libraryLogging.Debug("Error setting address")
			return err
		}
	} else if err != nil {
		libraryLogging.Debug("Error checking address")
		return err
	}
	libraryLogging.Debug("Checking address ok")
	libraryLogging.Debug("Checking if link is down")
	err = checkIfLinkDown(interfaceName)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Setting link down")
		err = linkDown(interfaceName)
		if err != nil {
			libraryLogging.Debug("Error setting link down")
			return err
		}
	} else if err != nil {
		return err
	}
	libraryLogging.Debug("Link is down")
	libraryLogging.Debug("Setting link up")
	err = linkUp(interfaceName)
	if err != nil {
		libraryLogging.Debug("Error setting link up")
		return err
	}
	libraryLogging.Debug("Link is up")
	libraryLogging.Debug("Checking default route")
	err = checkDefaultRoute(interfaceName, defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Adding default route")
		err = addDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			libraryLogging.Debug("Error adding default route")
			return err
		}
	} else if err != nil {
		libraryLogging.Debug("Error checking default route")
		return err
	}
	libraryLogging.Debug("Checking default route ok")
	libraryLogging.Debug("Checking server rule")
	err = checkServerRule(endpointIP)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Deleting server rule")
		err = deleteServerRule()
		if err != nil {
			libraryLogging.Debug("Error deleting server rule")
			return err
		}
		libraryLogging.Debug("Server rule deleted")
		libraryLogging.Debug("Adding server rule")
		err = addServerRule(endpointIP)
		if err != nil {
			libraryLogging.Debug("Error adding server rule")
			return err
		}
		libraryLogging.Debug("Server rule added")
	} else if err != nil {
		libraryLogging.Debug("Error checking server rule")
		return err
	}
	libraryLogging.Debug("Checking server rule ok")
	libraryLogging.Debug("Checking lookup rule")
	err = checkLookupRule(defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		libraryLogging.Debug("Adding lookup rule")
		err = addLookupRule(defaultRouteTable)
		if err != nil {
			libraryLogging.Debug("Error adding lookup rule")
			return err
		}
		libraryLogging.Debug("Lookup rule added")
	} else if err != nil {
		libraryLogging.Debug("Error checking lookup rule")
		return err
	}
	libraryLogging.Debug("Checking lookup rule ok")
	return nil
}
