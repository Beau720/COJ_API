package database

import "fmt"

// struct for connection string variables
type Config struct {
	Username       string
	Password       string
	ConnectionType string
	Host           string
	Port           int
	Name           string
}

//fucntion that handles the connection string
func (c *Config) ConnectString() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", c.Username, c.Password, c.ConnectionType, c.Host, c.Port, c.Name)
}
