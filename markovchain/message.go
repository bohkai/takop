package markovchain

import (
	"github.com/bluele/mecab-golang"
	"github.com/bwmarrin/discordgo"
	"fmt"
)

func Chain(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, v := range m.Mentions {
		if v.ID != s.State.User.ID {
			return
		}

		mecab, err := mecab.New("-Owakati")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "mecabなんで死んだ……")
			return
		}
		defer mecab.Destroy()

		messages, err := s.ChannelMessages(m.ChannelID, 100, m.Message.ID, "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Discord元気出すッピ！")
			return
		}

		markovBlocks := [][]string{}
		for _, v := range messages {
			if len(v.Content) <= 0 {
				continue
			}

			firstStr := v.Content[0]
			if v.MentionEveryone ||
				len(v.Mentions) > 0 ||
				firstStr == '@' ||
				firstStr == '#' ||
				firstStr == '?' ||
				firstStr == ','{
				continue
			}

			data, err := ParseToNode(mecab, v.Content)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "元気出すッピ！")
				return
			}

			if len(data) <= 1 {
				continue
			}

			elems := GetMarkovBlocks(data)
			markovBlocks = append(markovBlocks, elems...)
		}

		elm := MarkovChainExec(markovBlocks)
		text := TextGenerate(elm)

		fmt.Println(text)
		s.ChannelMessageSend(m.ChannelID, text)
	}
}

func TextGenerate(array []string) string {
	ret := ""
	for _, s := range array {
		if s == "#This is empty#" {
			continue
		}

		if len([]rune(ret)) >= 90 {
			break
		}

		ret += s
	}

	return ret
}
