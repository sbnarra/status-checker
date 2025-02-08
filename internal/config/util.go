package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func withIntDefault(key string, fallback int) (int, error) {
	return withDefault(key, fallback, func(s string) (int, error) {
		if num, err := strconv.Atoi(s); err != nil {
			return 0, err
		} else {
			return num, nil
		}
	})
}

func withBoolDefault(key string, fallback bool) (bool, error) {
	return withDefault(key, fallback, func(s string) (bool, error) {
		return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(s)), nil
	})
}

func withStrDefault(key string, fallback string) (string, error) {
	return withDefault(key, fallback, func(s string) (string, error) { return s, nil })
}

func withStrArrDefault(key string, fallback []string) ([]string, error) {
	return withDefault(key, fallback, func(s string) ([]string, error) { return strings.Split(s, ","), nil })
}

func withDefault[A any](key string, fallback A, conv func(string) (A, error)) (A, error) {
	if val := os.Getenv(key); val == "" {
		return fallback, nil
	} else if ret, err := conv(val); err != nil {
		return fallback, fmt.Errorf("%s error: %w", key, err)
	} else {
		return ret, nil
	}
}

var toBytesRegex = regexp.MustCompile(`(?i)^(\d+)(B|KB|MB|GB)$`)
var toBytesUnits = map[string]uintptr{
	"B":  1,
	"KB": 1024,
	"MB": 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
}

func toBytes(sizeStr string) (uintptr, error) {
	if matches := toBytesRegex.FindStringSubmatch(sizeStr); len(matches) != 3 {
		return 0, errors.New("invalid size format: " + sizeStr + ": 1B/1KB/1MB/1GB")
	} else if value, err := strconv.ParseUint(matches[1], 10, 64); err != nil {
		return 0, err
	} else if multiplier, exists := toBytesUnits[matches[2]]; !exists {
		return 0, errors.New("unknown unit: " + matches[2])
	} else {
		result := uintptr(value) * multiplier
		return result, nil
	}
}
