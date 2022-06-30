package handlers

import (
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv"
	"golang.org/x/exp/slices"
)

func ServerHasPlayerByName(server mvdsv.Mvdsv, playerName string) bool {
	if 0 == server.PlayerSlots.Used {
		return false
	}

	for _, c := range server.Players {
		normalizedName := strings.ToLower(c.Name.ToPlainString())

		if strings.Contains(normalizedName, strings.ToLower(playerName)) {
			return true
		}
	}

	return false
}

func ServerHasClientByName(server mvdsv.Mvdsv, clientName string) bool {
	clientCount := server.PlayerSlots.Used + server.SpectatorSlots.Used + server.QtvStream.SpectatorCount

	if 0 == clientCount {
		return false
	}

	if server.SpectatorSlots.Used > 0 && slices.Contains(server.SpectatorNames, clientName) {
		return true
	}

	if server.QtvStream.SpectatorCount > 0 && slices.Contains(server.QtvStream.SpectatorNames, clientName) {
		return true
	}

	if 0 == server.PlayerSlots.Used {
		return false
	}

	for _, c := range server.Players {
		if strings.ToLower(clientName) == strings.ToLower(c.Name.ToPlainString()) {
			return true
		}
	}

	return false
}
