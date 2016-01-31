// Copyright (C) 2015 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package dialer

import (
	"net"
	"time"
)

// Dial tries dialing via proxy if a proxy is configured, and falls back to
// a direct connection if no proxy is defined, or connecting via proxy fails.
func Dial(network, addr string) (net.Conn, error) {
	if usingProxy {
		return dialWithFallback(proxyDialer.Dial, net.Dial, network, addr)
	}
	return net.Dial(network, addr)
}

// DialTimeout tries dialing via proxy with a timeout if a proxy is configured,
// and falls back to a direct connection if no proxy is defined, or connecting
// via proxy fails. The timeout can potentially be applied twice, once trying
// to connect via the proxy connection, and second time trying to connect
// directly.
func DialTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	if usingProxy {
		// Because the proxy package is poorly structured, we have to
		// construct a struct that matches proxy.Dialer but has a timeout
		// and reconstrcut the proxy dialer using that, in order to be able to
		// set a timeout.
		dd := &timeoutDirectDialer{
			timeout: timeout,
		}
		// Check if the dialer we are getting is not timeoutDirectDialer we just
		// created. It could happen that usingProxy is true, but getDialer
		// returns timeoutDirectDialer due to env vars changing.
		if timeoutProxyDialer := getDialer(dd); timeoutProxyDialer != dd {
			directDialFunc := func(inetwork, iaddr string) (net.Conn, error) {
				return net.DialTimeout(inetwork, iaddr, timeout)
			}
			return dialWithFallback(timeoutProxyDialer.Dial, directDialFunc, network, addr)
		}
	}
	return net.DialTimeout(network, addr, timeout)
}

// SetTCPOptions sets syncthings default TCP options on a TCP connection
func SetTCPOptions(conn *net.TCPConn) error {
	var err error
	if err = conn.SetLinger(0); err != nil {
		return err
	}
	if err = conn.SetNoDelay(false); err != nil {
		return err
	}
	if err = conn.SetKeepAlivePeriod(60 * time.Second); err != nil {
		return err
	}
	if err = conn.SetKeepAlive(true); err != nil {
		return err
	}
	return nil
}
