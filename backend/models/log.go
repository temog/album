package models

import(
  "os"
  log "github.com/sirupsen/logrus"
  "time"
  "github.com/davecgh/go-spew/spew" // var_dump的なやつー
)

func LogError(err error) {
  today := Ymd()
  spew.Dump(today)

  file, err := os.OpenFile("./logs/error/" + today + ".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
      log.Fatal("error opening file :", err.Error())
  }

  log.SetFormatter(&log.JSONFormatter{})
  log.SetOutput(file)
  log.Error(err)
  file.Close()
}

func Ymd() string {
  return time.Now().Format("2006-01-02")
}
