// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix || wasip1 || windows

package adns

import (
	"context"
	"net"
	"sync"

	"github.com/qompassai/adns/internal/bytealg"
)

var onceReadProtocols sync.Once

// readProtocols loads contents of /etc/protocols into protocols map
// for quick access.
func readProtocols() {
	file, err := open("/etc/protocols")
	if err != nil {
		return
	}
	defer file.close()

	for line, ok := file.readLine(); ok; line, ok = file.readLine() {
		// tcp    6   TCP    # transmission control protocol
		if i := bytealg.IndexByteString(line, '#'); i >= 0 {
			line = line[0:i]
		}
		f := getFields(line)
		if len(f) < 2 {
			continue
		}
		if proto, _, ok := dtoi(f[1]); ok {
			if _, ok := protocols[f[0]]; !ok {
				protocols[f[0]] = proto
			}
			for _, alias := range f[2:] {
				if _, ok := protocols[alias]; !ok {
					protocols[alias] = proto
				}
			}
		}
	}
}

// lookupProtocol looks up IP protocol name in /etc/protocols and
// returns correspondent protocol number.
func lookupProtocol(_ context.Context, name string) (int, error) {
	onceReadProtocols.Do(readProtocols)
	return lookupProtocolMap(name)
}

func (r *Resolver) lookupHost(ctx context.Context, host string) (addrs []string, result Result, err error) {
	order, conf := systemConf().hostLookupOrder(r, host)
	return r.goLookupHostOrder(ctx, host, order, conf)
}

func (r *Resolver) lookupIP(ctx context.Context, network, host string) (addrs []net.IPAddr, result Result, err error) {
	if r.preferGo() {
		return r.goLookupIP(ctx, network, host)
	}
	order, conf := systemConf().hostLookupOrder(r, host)
	ips, _, result, err := r.goLookupIPCNAMEOrder(ctx, network, host, order, conf)
	return ips, result, err
}

func (r *Resolver) lookupPort(ctx context.Context, network, service string) (int, error) {
	// Port lookup is not a DNS operation.
	// Prefer the cgo resolver if possible.
	return goLookupPort(network, service)
}

func (r *Resolver) lookupCNAME(ctx context.Context, name string) (string, Result, error) {
	order, conf := systemConf().hostLookupOrder(r, name)
	return r.goLookupCNAME(ctx, name, order, conf)
}

func (r *Resolver) lookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, Result, error) {
	return r.goLookupSRV(ctx, service, proto, name)
}

func (r *Resolver) lookupMX(ctx context.Context, name string) ([]*net.MX, Result, error) {
	return r.goLookupMX(ctx, name)
}

func (r *Resolver) lookupNS(ctx context.Context, name string) ([]*net.NS, Result, error) {
	return r.goLookupNS(ctx, name)
}

func (r *Resolver) lookupTXT(ctx context.Context, name string) ([]string, Result, error) {
	return r.goLookupTXT(ctx, name)
}

func (r *Resolver) lookupAddr(ctx context.Context, addr string) ([]string, Result, error) {
	order, conf := systemConf().addrLookupOrder(r, addr)
	return r.goLookupPTR(ctx, addr, order, conf)
}

func (r *Resolver) lookupTLSA(ctx context.Context, port int, protocol, host string) ([]TLSA, Result, error) {
	return r.goLookupTLSA(ctx, port, protocol, host)
}
