package markovchain

import (
	"github.com/bluele/mecab-golang"
)

func ParseToNode(m *mecab.MeCab, input string) ([]string, error) {
	words := []string{}
	tg, err := m.NewTagger()
	if err != nil {
		return nil, err
	}
	defer	tg.Destroy()

	lt, err := m.NewLattice(input)
	if err != nil {
		return nil, err
	}
	defer	lt.Destroy()

	node := tg.ParseToNode(lt)
	for {
		if node.Surface() != "" {
			words = append(words, node.Surface())
		}

		if node.Next() != nil {
			break
		}
	}

	return words, nil
}