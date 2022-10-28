package serverindex

import (
	"errors"
	"sort"

	"github.com/vikpe/qw-hub-api/pkg/qnet"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
)

type ServerIndex map[string]qserver.GenericServer

func New(servers []qserver.GenericServer) ServerIndex {
	index := make(ServerIndex, 0)

	for _, server := range servers {
		index[server.Address] = server
	}

	return index
}

func (i ServerIndex) Servers() []qserver.GenericServer {
	servers := make([]qserver.GenericServer, 0)

	for _, server := range i {
		servers = append(servers, server)
	}

	if len(servers) > 0 {
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].Address < servers[j].Address
		})
	}

	return servers
}

func (i ServerIndex) Update(servers []qserver.GenericServer) {
	for _, server := range servers {
		i[server.Address] = server
	}
}

func (i ServerIndex) ActiveAddresses() []string {
	activeAddresses := make([]string, 0)

	for _, server := range i.Servers() {
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

func (i ServerIndex) Get(address string) (qserver.GenericServer, error) {
	ipHostPort, err := qnet.ToIpHostPort(address)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	if server, ok := i[ipHostPort]; ok {
		return server, nil
	}

	return qserver.GenericServer{}, errors.New("server not found")
}
