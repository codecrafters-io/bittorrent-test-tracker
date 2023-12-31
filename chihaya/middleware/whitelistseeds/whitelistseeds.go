// Package clientapproval implements a Hook that whitelists seeds to return
// on an Announce call.
package whitelistseeds

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chihaya/chihaya/bittorrent"
	"github.com/chihaya/chihaya/middleware"
	"github.com/chihaya/chihaya/pkg/log"
)

// Name is the name by which this middleware is registered with Chihaya.
const Name = "whitelist seeds"

var _ middleware.Driver = driver{}

type driver struct{}

func init() {
	middleware.RegisterDriver(Name, driver{})
}

func (d driver) NewHook(optionBytes []byte) (middleware.Hook, error) {
	clientIpsStr := os.Getenv("EXPOSED_CLIENT_IPS")

	if clientIpsStr == "" {
		return nil, fmt.Errorf("expected EXPOSED_CLIENT_IPS to be present")
	} else {
		return NewHook(strings.Split(clientIpsStr, ","))
	}
}

type hook struct {
	whitelistedSeederIps map[string]bool
}

// NewHook returns an instance of the seed whitelisting middleware.
func NewHook(whitelistedSeederIps []string) (middleware.Hook, error) {
	h := &hook{
		whitelistedSeederIps: make(map[string]bool),
	}

	for _, ip := range whitelistedSeederIps {
		h.whitelistedSeederIps[ip] = true
	}

	if len(whitelistedSeederIps) == 0 {
		return nil, fmt.Errorf("expected whitelisted seeder ips to be present")
	}

	return h, nil
}

func (h *hook) HandleAnnounce(ctx context.Context, req *bittorrent.AnnounceRequest, resp *bittorrent.AnnounceResponse) (context.Context, error) {
	log.Debug(fmt.Sprintf("whitelistseeds: executed, peers length: %v", len(resp.IPv4Peers)))

	filteredIPv4Peers := make([]bittorrent.Peer, 0, len(resp.IPv4Peers))
	for _, peer := range resp.IPv4Peers {
		log.Debug(fmt.Sprintf("whielistseeds: iterating, peer: %v", peer.IP.String()))

		if _, whitelisted := h.whitelistedSeederIps[peer.IP.String()]; whitelisted {
			log.Debug(fmt.Sprintf("whitelistseeds: whitelisting peer %v", peer.IP.String()))
			filteredIPv4Peers = append(filteredIPv4Peers, peer)
		} else {
			log.Debug(fmt.Sprintf("whitelistseeds: ignoring peer %v", peer.IP.String()))
		}
	}

	resp.IPv4Peers = filteredIPv4Peers

	filteredIPv6Peers := make([]bittorrent.Peer, 0, len(resp.IPv6Peers))
	for _, peer := range resp.IPv6Peers {
		if _, whitelisted := h.whitelistedSeederIps[peer.IP.String()]; whitelisted {
			filteredIPv6Peers = append(filteredIPv6Peers, peer)
		}
	}

	resp.IPv6Peers = filteredIPv6Peers

	return ctx, nil
}

func (h *hook) HandleScrape(ctx context.Context, req *bittorrent.ScrapeRequest, resp *bittorrent.ScrapeResponse) (context.Context, error) {
	// Scrapes don't require any protection.
	return ctx, nil
}
