package bot

import "math/rand"

func shuffle(ids []string) {
	n := len(ids)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		ids[i], ids[j] = ids[j], ids[i]
	}
}
