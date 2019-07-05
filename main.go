package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var config Config

func main() {

	fileReader, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = config.populateFromReader(fileReader)

	if err != nil {
		fmt.Println(err)
		return
	}

	dg, err := discordgo.New("Bot " + config.Token)
	dg.State.MaxMessageCount = 1000
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	//data, _ := json.Marshal(config)
	//fmt.Printf(string(data))

	dg.AddHandler(messageCreate)
	dg.AddHandler(voiceStateUpdate)
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageReactionRemove)
	dg.AddHandler(ready)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	dm, _ := comesFromDM(s, m)
	if dm {
		return
	}

	isAdmin, err := memberHasPermission(s, m.GuildID, m.Message.Author.ID, 8)
	if err != nil {
		return
	}

	if m.Content == "!w2" {
		if m.ChannelID == config.Guilds[m.GuildID].TalkChannelID {
			ID, err := createNewWTRoom()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
			}
			s.ChannelMessageSend(m.ChannelID, "https://www.watch2gether.com/rooms/"+ID)
		}
		return
	}
	if isAdmin {
		if m.Content == "!clear" {
			clearChannel(s, m.ChannelID, 100)
			return
		}
	}

}

func voiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	if m.ChannelID == "" {
		err := s.GuildMemberRoleRemove(m.GuildID, m.UserID, config.Guilds[m.GuildID].TalkRoleID)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := s.GuildMemberRoleAdd(m.GuildID, m.UserID, config.Guilds[m.GuildID].TalkRoleID)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.MessageID == config.Guilds[m.GuildID].ControlMessageID {
		for _, role := range config.Guilds[m.GuildID].GameRoles {
			if m.Emoji.Name == role.Emoji {
				s.GuildMemberRoleAdd(m.GuildID, m.UserID, role.RoleID)
				break
			}
		}
	}
}

func messageReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.MessageID == config.Guilds[m.GuildID].ControlMessageID {
		for _, role := range config.Guilds[m.GuildID].GameRoles {
			if m.Emoji.Name == role.Emoji {
				s.GuildMemberRoleRemove(m.GuildID, m.UserID, role.RoleID)
				break
			}
		}
	}
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	for guild, conf := range config.Guilds {
		clearChannel(s, conf.BotChannelID, 10)
		message, err := s.ChannelMessageSendEmbed(conf.BotChannelID, config.GameRoleEmbed)
		if err != nil {
			fmt.Println(err)
			return
		}
		config.Guilds[guild].ControlMessageID = message.ID

		for _, role := range config.Guilds[guild].GameRoles {
			s.MessageReactionAdd(conf.BotChannelID, message.ID, role.Emoji)
		}
	}
}
