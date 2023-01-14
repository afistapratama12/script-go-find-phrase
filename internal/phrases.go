package internal

import (
	"math/rand"
	"strings"
	"time"
)

func GetPhrase(data []string, length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var temp = make(map[string]*struct{})

	for {
		if _, ok := temp[data[rand.Intn(len(data))]]; !ok {
			temp[data[rand.Intn(len(data))]] = &struct{}{}
		}

		if len(temp) == length {
			return keyToString(temp)
		}
	}
}

func keyToString(m map[string]*struct{}) string {
	var result []string

	for k := range m {
		result = append(result, k)
	}

	return strings.Join(result, " ")
}
