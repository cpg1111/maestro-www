package main

import (
	"flag"
	"os"
)

var (
	flags = map[string]interface{}{
		"host": flag.String("host", "127.0.0.1", "host for server to bind on"),
		"port": flag.Int("port", "8080", "port for server to listen on"),
		"cert": flag.String("cert-path", "", "path to tls cert"),
		"key":  flag.String("key-path", "", "path to tls key"),
	}
	env = map[string]interface{}{
		"host": os.Getenv("MAESTRO_WWW_HOST"),
		"port": os.Getenv("MAESTRO_WWW_PORT"),
		"cert": os.Getenv("MAESTRO_WWW_CERT_PATH"),
		"key":  os.Getenv("MAESTRO_WWW_KEY_PATH"),
	}
)

func flagEnvMerge() (opts map[string]interface{}) {
	for i := range env {
		if len(env[i]) > 0 {
			opts[i] = env[i]
		} else {
			opts[i] = flags[i]
		}
	}
	return
}

func main() {
	opts := flagEnvMerge()
}
