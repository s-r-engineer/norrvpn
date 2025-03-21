package main

import "errors"

func execWGdown(interfaceName, interfaceIP, defaultRouteTable string) (err error) {
	if err = checkDefaultRoute(interfaceName, defaultRouteTable); err == nil || errors.Is(err, checkErrorInstance) {
		err = deleteDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			return err
		}
	}
	err = deleteServerRule()
	if err != nil {
		return err
	}
	err = deleteLookupRule(defaultRouteTable)
	if err != nil {
		return err
	}
	err = linkDown(interfaceName)
	if err != nil {
		return err
	}
	err = deleteIpAddress(interfaceIP, interfaceName)
	if err != nil {
		return err
	}
	return deleteInterface(interfaceName)
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
