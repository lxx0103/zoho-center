package main

import (
	"os"
	"zoho-center/cmd"
)

// @title zoho-center API
// @version 1.0
// @description API for zoho-center.

// @contact.name Lewis
// @contact.email lxx0103@yahoo.com

// @host 0.0.0.0:8080
// @BasePath /
func main() {
	cmd.Run(os.Args)
}
