package cites_storage

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

//go:embed cites.txt
var citesData string

type Storage struct {
	data []string
}

func Load() *Storage {
	cites, err := readCites(citesData)
	if err != nil {
		log.Panicf("Failed to read cites: %s\n", err)
	}

	return &Storage{
		data: cites,
	}
}

func (s *Storage) RandomCite() string {
	index := rand.Intn(len(s.data))

	return s.data[index]
}

func readCites(content string) ([]string, error) {
	ret := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(content))

	buffer := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if buffer != "" {
				ret = append(ret, buffer)
			}
			buffer = ""
			continue
		}

		buffer += fmt.Sprintf("%s\n", line)
	}
	if buffer != "" {
		ret = append(ret, buffer)
	}

	return ret, scanner.Err()
}
