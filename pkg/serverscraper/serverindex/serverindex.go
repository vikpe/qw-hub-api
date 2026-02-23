package serverindex

import (
	"errors"
	"sort"
	"sync"

	"github.com/vikpe/qw-hub-api/pkg/qnet"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
)

type ServerIndex struct {
	mu               sync.RWMutex
	serversByAddress map[string]qserver.GenericServer
}

func New(servers []qserver.GenericServer) *ServerIndex {
	serversByAddress := make(map[string]qserver.GenericServer, len(servers))

	for _, server := range servers {
		serversByAddress[server.Address] = server
	}

	return &ServerIndex{
		serversByAddress: serversByAddress,
	}
}

func (idx *ServerIndex) Servers() []qserver.GenericServer {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	servers := make([]qserver.GenericServer, 0, len(idx.serversByAddress))
	for _, server := range idx.serversByAddress {
		servers = append(servers, server)
	}

	sort.Slice(servers, func(a, b int) bool {
		return servers[a].Address < servers[b].Address
	})

	return servers
}

func (idx *ServerIndex) Upsert(servers []qserver.GenericServer) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.serversByAddress == nil {
		idx.serversByAddress = make(map[string]qserver.GenericServer, len(servers))
	}

	for _, server := range servers {
		idx.serversByAddress[server.Address] = server
	}
}

func (idx *ServerIndex) ActiveAddresses() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	activeAddresses := make([]string, 0, len(idx.serversByAddress))
	for _, server := range idx.serversByAddress {
		if len(server.Clients) > 0 && hasHumanPlayers(server.Clients) {
			activeAddresses = append(activeAddresses, server.Address)
		}
	}
	return activeAddresses
}

func hasHumanPlayers(clients []qclient.Client) bool {
	for _, c := range clients {
		if !c.IsSpectator() && !c.IsBot() {
			return true
		}
	}

	return false
}

func (idx *ServerIndex) Get(address string) (qserver.GenericServer, error) {
	ipHostPort, err := qnet.ToIpHostPort(address)
	if err != nil {
		return qserver.GenericServer{}, err
	}

	idx.mu.RLock()
	server, ok := idx.serversByAddress[ipHostPort]
	idx.mu.RUnlock()

	if ok {
		return server, nil
	}

	return qserver.GenericServer{}, errors.New("server not found")
}
