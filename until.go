package etcd

import (
	"crypto/sha256"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

func removeDuplicationMap(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

func getRandomString(n int) string {
	s := fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewV4().String()+strconv.FormatInt(time.Now().UnixNano(), 10))))

	randBytes := make([]byte, len(s)/2)
	rand.Read(randBytes)
	s1 := fmt.Sprintf("%x", randBytes)

	return s[:n-3] + s1[15:18]
}
