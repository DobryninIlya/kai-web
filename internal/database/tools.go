package database

import (
	"fmt"
	"strings"
)

func GetShortenName(name string) string {
	parts := strings.Split(name, " ")
	fmt.Println(parts)
	return parts[0]
}
