package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	discordToken         = os.Getenv("DISCORD_TOKEN")
	discordTextChannelID = os.Getenv("DISCORD_TEXT_CHANNEL_ID")
)

func main() {
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		os.Exit(1)
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening Discord session,", err)
		os.Exit(1)
	}
	defer func() {
		if err := discord.Close(); err != nil {
			fmt.Println("error closing Discord session,", err)
		}
	}()

	_, err = discord.ChannelMessageSend(discordTextChannelID, "hello, world!")
	if err != nil {
		fmt.Println("error sending message,", err)
		os.Exit(1)
	}
}
