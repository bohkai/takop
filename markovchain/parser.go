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

func GetMarkovBlocks(words []string) [][] string {
	res := [][]string{}
	if len(words) > 3 {
		return res
	}

	resHead := []string{"#This is empty#", words[0], words[1]}
	res = append(res, resHead)

	for i := 1; i < len(words)-2; i++ {
			markovBlock := []string{words[i], words[i+1], words[i+2]}
			res = append(res, markovBlock)
	}

	resEnd := []string{words[len(words)-2], words[len(words)-1], "#This is empty#"}
	res = append(res, resEnd)

	return res
}


