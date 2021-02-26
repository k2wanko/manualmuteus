package bot

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *bot) sendLobbyMessage() {
	b.printDebugLog()

	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: Name,
		Color: ColorGreen,
		Description: fmt.Sprintf("**ゲームに参加するには** %s **をクリック**\n\nBOT操作(ホストのみリアクションで操作可能)\nゲーム・タスク中は %s ディスカッション中は %s\nロビー画面は %s ホスト終了は %s",
			b.emojiMessageFormat(JoinEmoji), b.emojiMessageFormat(GameEmoji), b.emojiMessageFormat(DiscussEmoji), b.emojiMessageFormat(LobbyEmoji), b.emojiMessageFormat(EndEmoji)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Host",
				Value:  fmt.Sprintf("<@%s>", b.host.ID),
				Inline: true,
			},
			{
				Name:   "State",
				Value:  string(b.state),
				Inline: true,
			},
			{
				Name:   "Players",
				Value:  strconv.Itoa(len(b.game.GetCrewmateUsers()) + len(b.game.GetImposterUsers()) + len(b.game.GetDeadUsers())),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	users := make([]string, 0)
	users = append(users, b.game.GetCrewmateUsers()...)
	users = append(users, b.game.GetImposterUsers()...)
	users = append(users, b.game.GetDeadUsers()...)
	sort.Strings(users)
	for i := range users {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Player %d", i+1),
			Value:  fmt.Sprintf("<@%s>", users[i]),
			Inline: true,
		})
	}

	if b.lobbyMessageID == "" {
		message, err := b.discord.ChannelMessageSendEmbed(b.LobbyChannelID, embed)
		if err != nil {
			log.Printf("error sending lobby message channelID=%s: %v", b.LobbyChannelID, err)
			return
		}
		b.lobbyMessageID = message.ID

		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(JoinEmoji))
		if err != nil {
			log.Printf("error add reaction to lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, message.ID, err)
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(GameEmoji))
		if err != nil {
			log.Printf("error add reaction to lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, message.ID, err)
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(DiscussEmoji))
		if err != nil {
			log.Printf("error add reaction to lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, message.ID, err)
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(LobbyEmoji))
		if err != nil {
			log.Printf("error add reaction to lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, message.ID, err)
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(EndEmoji))
		if err != nil {
			log.Printf("error add reaction to lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, message.ID, err)
		}

	} else {
		_, err := b.discord.ChannelMessageEditEmbed(b.LobbyChannelID, b.lobbyMessageID, embed)
		if err != nil {
			log.Printf("error editing lobby message channelID=%s messageID=%s: %v", b.LobbyChannelID, b.lobbyMessageID, err)
			return
		}
	}
}

func (b *bot) sendPrivateMessages() {
	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: Name,
		Description: fmt.Sprintf("**ゲーム開始**\nもし、あなたが殺された・追放されたなら %s、Imposterなら %s のリアクションをしてね",
			b.emojiMessageFormat(DeadEmoji), b.emojiMessageFormat(ImposterEmoji)),
	}

	for channelID, userID := range b.userChannelMap {
		message, err := b.discord.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Printf("error sending private message userID=%s: %v", userID, err)
			continue
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(DeadEmoji))
		if err != nil {
			log.Printf("error add reaction to private message userID=%s: %v", userID, err)
			continue
		}
		err = b.discord.MessageReactionAdd(message.ChannelID, message.ID, b.emojiAPIName(ImposterEmoji))
		if err != nil {
			log.Printf("error add reaction to private message userID=%s: %v", userID, err)
			continue
		}
	}
}
