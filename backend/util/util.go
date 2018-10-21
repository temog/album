package util

import(
  "crypto/sha512"
  "encoding/hex"
  "math/rand"
  "time"
  "strconv"
  "github.com/davecgh/go-spew/spew" // var_dump的なやつー
)

func Sha512String(str string) string {
    bytes := sha512.Sum512([]byte(str))
    return hex.EncodeToString(bytes[:])
}

func Random(i int) int {
  rand.Seed(time.Now().UnixNano())
  return rand.Intn(i)
}

func CastInt(s string) int {
  i, _ := strconv.Atoi(s)
  return i
}

func CastString(i int) string {
  return strconv.Itoa(i)
}

func Dump(src interface{}) {
  spew.Dump(src)
}
