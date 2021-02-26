package bot

import (
	"log"
)

func (b *bot) moveAllToBooth() {
	aliveUserMap := make(map[string]string) // map[userID]channelID
	aliveUsers := make([]string, 0)

	for i, userID := range b.game.GetCrewmateUsers() {
		aliveUserMap[userID] = b.BoothVoiceChannelIDs[i]
		aliveUsers = append(aliveUsers, userID)
	}
	for _, userID := range b.game.GetImposterUsers() {
		aliveUserMap[userID] = b.ImposterVoiceChannelID
		aliveUsers = append(aliveUsers, userID)
	}

	// Why shuffle? - To avoid detection by the order in which you move.
	shuffle(aliveUsers)

	for _, userID := range aliveUsers {
		channelID := aliveUserMap[userID]
		if err := b.discord.GuildMemberMove(b.GuildID, userID, &channelID); err != nil {
			log.Printf("error moving member to booth userID=%s channelID=%s: %v", userID, channelID, err)
		}
	}

	deadUsers := b.game.GetDeadUsers()
	for _, userID := range deadUsers {
		if err := b.discord.GuildMemberMove(b.GuildID, userID, &b.LimboVoiceChannelID); err != nil {
			log.Printf("error moving member to booth userID=%s channelID=%s: %v", userID, b.LimboVoiceChannelID, err)
		}
		if err := b.discord.GuildMemberMute(b.GuildID, userID, false); err != nil {
			log.Printf("error unmuting member userID=%s: %v", userID, err)
		}
	}
}

func (b *bot) moveAllToLobby() {
	for _, userID := range b.game.GetDeadUsers() {
		if err := b.discord.GuildMemberMove(b.GuildID, userID, &b.LobbyVoiceChannelID); err != nil {
			log.Printf("error moving member to booth userID=%s channelID=%s: %v", userID, b.LobbyVoiceChannelID, err)
		}
		if err := b.discord.GuildMemberMute(b.GuildID, userID, true); err != nil {
			log.Printf("error muting member userID=%s: %v", userID, err)
		}
	}

	aliveUsers := make([]string, 0)
	aliveUsers = append(aliveUsers, b.game.GetCrewmateUsers()...)
	aliveUsers = append(aliveUsers, b.game.GetImposterUsers()...)

	// Why shuffle? - To avoid detection by the order in which you move.
	shuffle(aliveUsers)

	for _, userID := range aliveUsers {
		if err := b.discord.GuildMemberMove(b.GuildID, userID, &b.LobbyVoiceChannelID); err != nil {
			log.Printf("error moving member to booth userID=%s channelID=%s: %v", userID, b.LobbyVoiceChannelID, err)
		}
	}
}

func (b *bot) moveAllToLobbyAndUnmute() {
	allUsers := make([]string, 0)
	allUsers = append(allUsers, b.game.GetCrewmateUsers()...)
	allUsers = append(allUsers, b.game.GetImposterUsers()...)
	allUsers = append(allUsers, b.game.GetDeadUsers()...)
	for _, userID := range allUsers {
		if err := b.discord.GuildMemberMove(b.GuildID, userID, &b.LobbyVoiceChannelID); err != nil {
			log.Printf("error moving member to booth userID=%s channelID=%s: %v", userID, b.LobbyVoiceChannelID, err)
		}
		if err := b.discord.GuildMemberMute(b.GuildID, userID, false); err != nil {
			log.Printf("error unmuting member userID=%s: %v", userID, err)
		}
	}
}
