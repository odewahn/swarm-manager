package db

import (
  "github.com/odewahn/swarm-manager/models"
  "gopkg.in/redis.v3"
  "net/url"
  "log"
  "os"
)

var (
  redisClient *redis.Client
  initialized bool
)

// Helper function to parse a Redis string and return the host and pw
// needed for the redis options.
func ParseRedisUrl(redisUrl string) (string, string, error) {
	uri, uriErr := url.Parse(redisUrl)
	if uriErr != nil {
		return "", "", uriErr
	}
	var pwString string
	if uri.User != nil {
		pwString, _ = uri.User.Password()
	}
	return uri.Host, pwString, nil
}

// Initializes the redis client
func Init() {
  host, password, err := ParseRedisUrl(os.Getenv("REDIS_URL"))
  if err != nil {
    log.Fatal("Cannot parse redis URL: ", err)
  }
  // Initialize the redis client
  redisClient = redis.NewClient(&redis.Options{
    Addr:     host,
    Password: password,
    DB:       0,
  })
  initialized = true
}

func SaveContainer( c *models.Container) {
  if !initialized {
    log.Fatal("Data model is not initialized!  Call model.Init() before using this function.")
  }
  err := redisClient.HSet("containers",c.Hostname, c.Serialize()).Err()
  if err != nil {
    log.Println(err)
  }
}

func GetContainer (k string) (models.Container) {
  if !initialized {
    log.Fatal("Data model is not initialized!  Call model.Init() before using this function.")
  }
  s, err := redisClient.HGet("containers",k).Result()
  if err != nil {
    log.Println(err)
  }
  return models.DeserializeContainer(s)
}


func GetContainers() []models.Container {
  if !initialized {
    log.Fatal("Data model is not initialized!  Call model.Init() before using this function.")
  }
  var out []models.Container
  hostnames, _ := redisClient.HKeys("containers").Result()
  for _, h := range hostnames {
     c := GetContainer(h)
     out = append(out,c)
  }
  return out
}
