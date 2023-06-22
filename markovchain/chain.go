package markovchain

import(
	"math/rand"
	"time"
)

func FindBlocks(array [][]string, target string) [][]string {
	blocks := [][]string{}
	for _, s := range array {
			if s[0] == target {
					blocks = append(blocks, s)
			}
	}

	return blocks
}

func ConnectBlocks(array [][]string, dist []string) []string {
	rand.Seed((time.Now().Unix()))
	i := 0

	for _, word := range array[rand.Intn(len(array))] {
			if i != 0 {
					dist = append(dist, word)
			}
			i += 1
	}

	return dist
}

func MarkovChainExec(array [][]string) []string {
	ret := []string{}
	block := [][]string{}
	count := 0

	block = FindBlocks(array, "#This is empty#")
	ret = ConnectBlocks(block, ret)

	for ret[len(ret)-1] != "#This is empty#" {
			block = FindBlocks(array, ret[len(ret)-1])
			if len(block) == 0 {
					break
			}
			ret = ConnectBlocks(block, ret)

			count++
			if count == 150 {
					break
			}
	}

	return ret
}