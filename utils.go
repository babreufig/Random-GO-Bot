package main

import (
	"github.com/bwmarrin/discordgo"
)

func comesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}
	return channel.Type == discordgo.ChannelTypeDM, nil
}

func clearChannel(s *discordgo.Session, channelID string, num int) error {
	messages, _ := s.ChannelMessages(channelID, num, "", "", "")
	var messageIDs []string
	if len(messages) == 0 {
		return nil
	}

	for i := 0; i < len(messages); i++ {
		messageIDs = append(messageIDs, messages[i].ID)
	}

	err := s.ChannelMessagesBulkDelete(channelID, messageIDs)
	if err != nil {
		return err
	}
	return nil
}

func memberHasPermission(s *discordgo.Session, guildID string, userID string, permission int) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

func isMemeFriendly(guildID string, channelID string) bool {
	for _, i := range config.Guilds[guildID].MemeFriendlyChannels {
		if i == channelID {
			return true
		}
	}
	return false
}
