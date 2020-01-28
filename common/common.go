package common

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	customEpoch = int64(1300000000000)
	shardId     = int64(100)
)

// 64bits id = 41bits timestamp + 10bits shardId + 13 bits randomId
func GenerateID() string {
	nowMillis := time.Now().UnixNano() / 1e6
	seqId := rand.Int63n(1000)
	result := (nowMillis - customEpoch) << 23
	result |= shardId << 10
	result |= seqId % 1024
	return strconv.FormatInt(result, 10)
}
