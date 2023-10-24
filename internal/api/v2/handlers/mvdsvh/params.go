package mvdsvh

import (
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"strings"
)

type MvdsvParams struct {
	Empty     string `query:"empty" validate:"omitempty,oneof=include exclude only"`
	HasClient string `query:"has_client" validate:"omitempty,min=2"`
	HasPlayer string `query:"has_player" validate:"omitempty,min=2"`
	Hostname  string `query:"hostname" validate:"omitempty,min=2"`
	Limit     int    `query:"limit" validate:"omitempty"`
}

func NewMvdsvParams() *MvdsvParams {
	return &MvdsvParams{
		Empty:     "exclude",
		HasClient: "",
		HasPlayer: "",
		Hostname:  "",
		Limit:     100,
	}
}

func FilterByHostname(servers []mvdsv.Mvdsv, hostname string) []mvdsv.Mvdsv {
	if 0 == len(hostname) {
		return servers
	}

	result := make([]mvdsv.Mvdsv, 0)
	for _, server := range servers {
		needle := server.Settings.Get("hostname_parsed", server.Address)
		if strings.Contains(needle, hostname) {
			result = append(result, server)
		}
	}

	return result
}

func FilterByHasPlayer(servers []mvdsv.Mvdsv, hasPlayer string) []mvdsv.Mvdsv {
	if 0 == len(hasPlayer) {
		return servers
	}

	result := make([]mvdsv.Mvdsv, 0)
	for _, server := range servers {
		if analyze.HasPlayer(server, hasPlayer) {
			result = append(result, server)
		}
	}

	return result
}

func FilterByHasClient(servers []mvdsv.Mvdsv, hasClient string) []mvdsv.Mvdsv {
	if 0 == len(hasClient) {
		return servers
	}

	result := make([]mvdsv.Mvdsv, 0)
	for _, server := range servers {
		if analyze.HasClient(server, hasClient) {
			result = append(result, server)
		}
	}

	return result
}

func FilterByEmpty(servers []mvdsv.Mvdsv, empty string) []mvdsv.Mvdsv {
	if empty == "" || empty == "include" {
		return servers
	}

	result := make([]mvdsv.Mvdsv, 0)

	if empty == "exclude" {
		for _, server := range servers {
			if server.PlayerSlots.Used > 0 {
				result = append(result, server)
			}
		}
	} else if empty == "only" {
		for _, server := range servers {
			if server.PlayerSlots.Used == 0 {
				result = append(result, server)
			}
		}
	}

	return result
}

func FilterByParams(servers []mvdsv.Mvdsv, params *MvdsvParams) []mvdsv.Mvdsv {
	result := FilterByHostname(servers, params.Hostname)

	if params.HasPlayer != "" {
		result = FilterByHasPlayer(result, params.HasPlayer)
	} else if params.HasClient != "" {
		result = FilterByHasClient(result, params.HasClient)
	} else if params.Empty != "" {
		result = FilterByEmpty(result, params.Empty)
	}

	if params.Limit > 0 && len(result) > params.Limit {
		result = result[0:params.Limit]
	}

	return result
}
