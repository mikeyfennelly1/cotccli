package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	WebAppPort               int
	CollectorListenPort      int
	DesktopSysinfoListenPort int
}

// Load reads port configuration from ../.env.local relative to the given base directory.
// Pass the directory of the running binary or script (e.g. os.Executable()).
func Load() (*Config, error) {
	wd, _ := os.Getwd()
	envPath := wd + "/../.env.local"

	vals, err := parseEnvFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read %s: %w", envPath, err)
	}

	required := []string{"WEB_APP_PORT", "COLLECTOR_LISTEN_PORT", "DESKTOP_SYSINFO_LISTEN_PORT"}
	for _, key := range required {
		if _, ok := vals[key]; !ok {
			return nil, fmt.Errorf("config: missing required env var %s in %s", key, envPath)
		}
	}

	webAppPort, err := parseInt(vals, "WEB_APP_PORT")
	if err != nil {
		return nil, err
	}
	collectorPort, err := parseInt(vals, "COLLECTOR_LISTEN_PORT")
	if err != nil {
		return nil, err
	}
	desktopPort, err := parseInt(vals, "DESKTOP_SYSINFO_LISTEN_PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		WebAppPort:               webAppPort,
		CollectorListenPort:      collectorPort,
		DesktopSysinfoListenPort: desktopPort,
	}, nil
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vals := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		vals[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"`)
	}
	return vals, scanner.Err()
}

func parseInt(vals map[string]string, key string) (int, error) {
	n, err := strconv.Atoi(vals[key])
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer, got %q", key, vals[key])
	}
	return n, nil
}
