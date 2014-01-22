/**
 * Redis To GO Helper
 *
 * This file will parse "REDISTOGO_URL" from String to Host and Port
 */

package helper

import (
  "os"
  "net"
  "net/url"
  "strconv"
)

var redisToGoURL string = os.Getenv("REDISTOGO_URL")

func GetRedisToGoEnv() (host string, port uint, password string) {
  if len(redisToGoURL) <= 0 { // Fallback, for localhsot test
    redisToGoURL = "redis://@localhost:6379/"
  }

  dbURL, _ := url.Parse(redisToGoURL)
  dbHost, dbPort, _ := net.SplitHostPort(dbURL.Host)
  dbPortUint, _ := strconv.Atoi(dbPort)
  if dbURL.User != nil {
    if password, exists := dbURL.User.Password(); exists {
      return dbHost, uint(dbPortUint), password
    }
  }
  return dbHost, uint(dbPortUint), ""
}
