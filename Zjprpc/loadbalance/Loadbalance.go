package loadbalance

import (
	"math/rand"
	"myRpc/Zjprpc/common"
	"time"
)

func Random(urls []common.URL) common.URL {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(urls))
	return urls[i]
}
