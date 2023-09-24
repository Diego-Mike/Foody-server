package helpers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnoprstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate random number
func RandomInt(min, max int64) int {
	return int(min + rand.Int63n(max-min+1))
}

// generates random string based on alphabet
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}

func RandomUserName() string {
	return RandomString(7)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@hotmail.com", RandomUserName())
}

func RandomPicture() string {
	return fmt.Sprintf("https://platform-lookaside.fbsbx.com/platform/profilepic/?asid=2536413786521623&height=50&width=50&ext=1696650888&hash=AeTzVMPPwrnQML1Lh1w%s", RandomUserName())
}

func RandomProvider() string {
	providers := []string{"google", "facebook", "tik tok", "instagram"}
	k := len(providers)

	return providers[rand.Intn(k)]
}
