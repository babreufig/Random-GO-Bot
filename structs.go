package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

//Config Holds general configuration of the bot
type Config struct {
	Token         string                  `json:"token"`
	GameRoleEmbed *discordgo.MessageEmbed `json:"gameRoleEmbed"`
	Guilds        map[string]*GuildConfig `json:"guildID:GuildConfig"`
}

//GuildConfig Holds Guild specific configuration
type GuildConfig struct {
	TalkRoleID               string      `json:"talkRoleID"`
	TalkChannelID            string      `json:"talkChannelID"`
	MusicBotID               string      `json:"musicBotID"`
	BotChannelID             string      `json:"botChannelID"`
	GameRoles                []*Gamerole `json:"gameRoles"`
	ControlMessageID         string
	lastUsedExpensiveCommand time.Time
}

//Gamerole ..
type Gamerole struct {
	Emoji  string `json:"emoji"`
	RoleID string `json:"roleID"`
}

func (conf *Config) populateFromReader(reader io.Reader) error {
	jsonBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, conf)
	if err != nil {
		return err
	}
	return nil
}
