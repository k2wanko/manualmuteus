package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/miyukki/manualmuteus/game"

	"github.com/bwmarrin/discordgo"
)

const (
	Name          = "ManualMuteUs Bot"
	CommandPrefix = ".mm "
)

const (
	JoinEmoji     = "aured"
	GameEmoji     = "aushhhhhhh"
	DiscussEmoji  = "audisscuss"
	LobbyEmoji    = "aucustomize"
	EndEmoji      = "❌"
	DeadEmoji     = "aureddead"
	ImposterEmoji = "aukill"

	ColorGreen = 3066993
)

type Config struct {
	Token                  string
	GuildID                string
	LobbyChannelID         string
	LobbyVoiceChannelID    string
	ImposterVoiceChannelID string
	BoothVoiceChannelIDs   []string
	LimboVoiceChannelID    string
}

type Bot interface {
	Start() error
	Stop() error
}

type bot struct {
	Token                  string
	GuildID                string
	LobbyChannelID         string
	LobbyVoiceChannelID    string
	ImposterVoiceChannelID string
	BoothVoiceChannelIDs   []string
	LimboVoiceChannelID    string

	discord        *discordgo.Session
	me             *discordgo.User
	host           *discordgo.User
	lobbyMessageID string
	state          game.State
	game           game.Session
	guildEmojiMap  map[string]*discordgo.Emoji // map[emojiName]*discordgo.Emoji
	userChannelMap map[string]string           // map[channelID]userID
}

func New(cfg *Config) (Bot, error) {
	b := &bot{
		Token:                  cfg.Token,
		GuildID:                cfg.GuildID,
		LobbyChannelID:         cfg.LobbyChannelID,
		LobbyVoiceChannelID:    cfg.LobbyVoiceChannelID,
		ImposterVoiceChannelID: cfg.ImposterVoiceChannelID,
		BoothVoiceChannelIDs:   cfg.BoothVoiceChannelIDs,
		LimboVoiceChannelID:    cfg.LimboVoiceChannelID,
	}
	err := b.init()
	return b, err
}

func (b *bot) init() (err error) {
	b.discord, err = discordgo.New("Bot " + b.Token)
	if err != nil {
		err = fmt.Errorf("error create discordgo instance: %w", err)
		return
	}

	err = b.discord.Open()
	if err != nil {
		err = fmt.Errorf("error open discord session: %w", err)
		return
	}

	b.me, err = b.discord.User("@me")
	if err != nil {
		err = fmt.Errorf("error get own user information: %w", err)
		return
	}

	emojis, err := b.discord.GuildEmojis(b.GuildID)
	if err != nil {
		err = fmt.Errorf("error get guild emojis: %w", err)
		return
	}
	guildEmojiMap := make(map[string]*discordgo.Emoji)
	for i := range emojis {
		guildEmojiMap[emojis[i].Name] = emojis[i]
	}
	b.guildEmojiMap = guildEmojiMap

	b.state = game.StateMenu
	b.game = game.NewSession()

	if len(b.BoothVoiceChannelIDs) != 10 {
		err = fmt.Errorf("error BoothVoiceChannelIDs must have 10 ids")
		return
	}

	return
}

func (b *bot) Start() (err error) {
	b.discord.AddHandler(b.handleMessageCreate)
	b.discord.AddHandler(b.handleMessageReactionAdd)
	b.discord.AddHandler(b.handleMessageReactionRemove)

	return nil
}

func (b *bot) Stop() error {
	if err := b.discord.Close(); err != nil {
		return fmt.Errorf("error closing discord session: %w", err)
	}

	return nil
}

func (b *bot) isControllableUser(userID string) bool {
	return userID == b.host.ID
}

func (b *bot) generateUserChannel() {
	m := make(map[string]string)
	for _, uid := range b.game.GetCrewmateUsers() {
		ch, err := b.discord.UserChannelCreate(uid)
		if err != nil {
			log.Printf("error create user channel: %v", err)
			return
		}
		m[ch.ID] = uid
	}

	b.userChannelMap = m
}

func (b *bot) handleMessageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != b.LobbyChannelID {
		return
	}
	if !strings.HasPrefix(m.Content, CommandPrefix) {
		return
	}

	command := m.Content[len(CommandPrefix):]
	switch command {
	case "new":
		if b.host != nil {
			_, err := b.discord.ChannelMessageSendReply(m.ChannelID, "すでに別のゲームが開始されています", m.Reference())
			if err != nil {
				log.Printf("error sending reply message: %b", err)
			}
			return
		}
		b.host = m.Author
		b.sendLobbyMessage()
	}
}

func (b *bot) handleMessageReactionAdd(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == b.me.ID {
		return
	}

	var removeReaction bool
	defer func() {
		if !removeReaction {
			return
		}

		if err := b.discord.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID); err != nil {
			log.Printf("error remove message reaction: %v", err)
		}
	}()

	if r.MessageID == b.lobbyMessageID {
		switch r.Emoji.Name {
		case JoinEmoji:
			if b.state == game.StateMenu {
				log.Printf("new user joined userID=%s", r.UserID)
				b.game.AddCrewmateUser(r.UserID)
			}
			b.sendLobbyMessage()
		case GameEmoji:
			removeReaction = true
			if !b.isControllableUser(r.UserID) {
				break
			}

			if b.state == game.StateMenu {
				log.Printf("start new game")
				b.state = game.StateGame
				b.generateUserChannel()
				b.sendPrivateMessages()
				b.moveAllToBooth()
			} else if b.state == game.StateDiscuss {
				log.Printf("change state to game from discuss")
				b.state = game.StateGame
				b.moveAllToBooth()
			}
			b.sendLobbyMessage()
		case DiscussEmoji:
			removeReaction = true
			if !b.isControllableUser(r.UserID) {
				break
			}

			if b.state == game.StateGame {
				log.Printf("change state to discuss from game")
				b.state = game.StateDiscuss
				b.moveAllToLobby()
			}
			b.sendLobbyMessage()
		case LobbyEmoji:
			removeReaction = true
			if !b.isControllableUser(r.UserID) {
				break
			}

			log.Printf("end the game")
			b.state = game.StateMenu
			b.game.Reset()
			b.moveAllToLobbyAndUnmute()
			b.sendLobbyMessage()
		case EndEmoji:
			removeReaction = true
		}
	}
	if _, ok := b.userChannelMap[r.ChannelID]; ok {
		switch r.Emoji.Name {
		case DeadEmoji:
			if ok := b.game.DeleteCrewmateUser(r.UserID) || b.game.DeleteImposterUser(r.UserID); !ok {
				log.Printf("user is not belonging to crewmate or dead, ahhhhhhhh? userID=%s", r.UserID)
				return
			}
			b.game.AddDeadUser(r.UserID)
			if b.state == game.StateDiscuss {
				if err := b.discord.GuildMemberMute(b.GuildID, r.UserID, true); err != nil {
					log.Printf("error muting member userID=%s: %v", r.UserID, err)
				}
			} else {
				if err := b.discord.GuildMemberMove(b.GuildID, r.UserID, &b.LimboVoiceChannelID); err != nil {
					log.Printf("error moving member to booth userID=%s channelID=%s: %v", r.UserID, b.LimboVoiceChannelID, err)
				}
			}
		case ImposterEmoji:
			if ok := b.game.DeleteCrewmateUser(r.UserID); !ok {
				log.Printf("user is not belonging to crewmate, perhaps already dead? userID=%s", r.UserID)
				return
			}
			b.game.AddImposterUser(r.UserID)
			if b.state == game.StateGame {
				if err := b.discord.GuildMemberMove(b.GuildID, r.UserID, &b.ImposterVoiceChannelID); err != nil {
					log.Printf("error moving member to booth userID=%s channelID=%s: %v", r.UserID, b.ImposterVoiceChannelID, err)
				}
			}
		}
	}
}

func (b *bot) handleMessageReactionRemove(_ *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.MessageID == b.lobbyMessageID {
		switch r.Emoji.Name {
		case JoinEmoji:
			if b.state == game.StateMenu {
				log.Printf("user deleted userID=%s", r.UserID)
				b.game.DeleteCrewmateUser(r.UserID)
			}
			b.sendLobbyMessage()
		}
	}
}

func (b *bot) emojiAPIName(name string) string {
	if emoji, ok := b.guildEmojiMap[name]; ok {
		return emoji.APIName()
	}

	return name
}

func (b *bot) emojiMessageFormat(name string) string {
	if emoji, ok := b.guildEmojiMap[name]; ok {
		return emoji.MessageFormat()
	}

	return name
}

func (b *bot) printDebugLog() {
	log.Printf("[debug] state=%s", string(b.state))
	log.Printf("    crewmate(%d)=%v", len(b.game.GetCrewmateUsers()), b.game.GetCrewmateUsers())
	log.Printf("    imposter(%d)=%s", len(b.game.GetImposterUsers()), b.game.GetImposterUsers())
	log.Printf("    dead(%d)=%s", len(b.game.GetDeadUsers()), b.game.GetDeadUsers())
}
