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
	if !checkInterface(interfaceName) {
		err = addInterface(interfaceName)
		if err != nil {
			return err
		}
	}

	err = checkPrivateKey(interfaceName, privateKey)
	if errors.Is(err, checkErrorInstance) {
		err = setPrivateKey(interfaceName, privateKey)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfPeerOk(interfaceName, publicKey, endpointIP, defaultWGPort)
	if errors.Is(err, checkErrorInstance) {
		err = setPeer(interfaceName, publicKey, endpointIP, defaultWGPort)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfAddressOk(interfaceName, interfaceIP)
	if errors.Is(err, checkErrorInstance) {
		err = setAddress(interfaceName, interfaceIP)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfLinkDown(interfaceName)
	if errors.Is(err, checkErrorInstance) {
		err = linkDown(interfaceName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	err = linkUp(interfaceName)
	if err != nil {
		return err
	}

	err = checkDefaultRoute(interfaceName, defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		err = addDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkServerRule(endpointIP)
	if errors.Is(err, checkErrorInstance) {
		err = deleteServerRule()
		if err != nil {
			return err
		}
		err = addServerRule(endpointIP)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	err = checkLookupRule(defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		return addLookupRule(defaultRouteTable)
	}
	return err
}
