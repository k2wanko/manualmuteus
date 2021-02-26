package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/miyukki/manualmuteus/bot"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	var b bot.Bot
	b, err = bot.New(&bot.Config{
		Token:                  os.Getenv("TOKEN"),
		GuildID:                os.Getenv("GUILD_ID"),
		LobbyChannelID:         os.Getenv("LOBBY_CHANNEL_ID"),
		LobbyVoiceChannelID:    os.Getenv("LOBBY_VOICE_CHANNEL_ID"),
		BoothVoiceChannelIDs:   strings.Split(os.Getenv("BOOTH_VOICE_CHANNEL_IDS"), ","),
		ImposterVoiceChannelID: os.Getenv("IMPOSTER_VOICE_CHANNEL_ID"),
		LimboVoiceChannelID:    os.Getenv("LIMBO_VOICE_CHANNEL_ID"),
	})
	if err != nil {
		log.Fatalf("failed to initialize bot")
	}

	err = b.Start()
	if err != nil {
		log.Fatalf("failed to start bot")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	s := <-sig
	log.Printf("singal received: %s", s)

	err = b.Stop()
	if err != nil {
		log.Fatalf("failed to stop bot")
	}
}
