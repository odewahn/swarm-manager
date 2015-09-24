package models

import (
  "encoding/json"
  "log"
)

type Container struct {
  Hostname string `json:"hostname"`
  Domainname string `json:"domainname"`
  Image string `json:"image"`
  ContainerId string
	Status string
}

// Serializes a container as a string
func (c *Container) Serialize() (string) {
	out, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

func DeserializeContainer(s string)(Container) {
  var c Container
  err := json.Unmarshal([]byte(s),&c)
  if err !=nil {
    log.Println(err)
  }
  return c
}
