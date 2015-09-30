package models

import (
	"encoding/json"
	"log"
	"time"
)

type Container struct {
	Hostname    string    `json:"hostname"`
	Domainname  string    `json:"domainname"`
	Image       string    `json:"image"`
	Url         string    `json:"url"`
	ContainerId string    `json:"container_id"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	User        string    `json:"user"`
}

func (c *Container) IsActive() bool {
	out := false
	if c.Status == "ACTIVE" {
		out = true
	}
	return out
}

// Serializes a container as a string
func (c *Container) Serialize() string {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

func DeserializeContainer(s string) Container {
	var c Container
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		log.Println(err)
	}
	return c
}
