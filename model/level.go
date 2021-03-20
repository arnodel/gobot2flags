package model

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Level struct {
	Name        string
	Maze        *Maze
	BoardWidth  int
	BoardHeigth int
	ChipCost    int
	MoveCost    int
}

func LevelFromString(defaultName string, s string) (*Level, error) {
	lvl := Level{
		Name:        defaultName,
		BoardWidth:  9,
		BoardHeigth: 9,
		ChipCost:    10,
		MoveCost:    1,
	}
	var parseErrors []string
	for _, kv := range parseString(s) {
		var err error
		switch strings.ToLower(kv.k) {
		case "", "maze":
			m, err := MazeFromString(kv.v)
			if err == nil {
				lvl.Maze = m
			}
			lvl.Maze = m
		case "name":
			lvl.Name = strings.TrimSpace(kv.v)
		case "boardwidth":
			w, err := parseInt(kv.v)
			if err == nil {
				lvl.BoardWidth = w
			}
		case "boardheight":
			h, err := parseInt(kv.v)
			if err == nil {
				lvl.BoardHeigth = h
			}
		case "chipcost":
			c, err := parseInt(kv.v)
			if err == nil {
				lvl.ChipCost = c
			}
		case "movecost":
			c, err := parseInt(kv.v)
			if err == nil {
				lvl.MoveCost = c
			}
		}
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("%s: %s", kv.k, err))
		}
	}
	if lvl.Maze == nil {
		parseErrors = append(parseErrors, "Maze definition is missing")
	}
	if len(parseErrors) != 0 {
		return nil, errors.New(strings.Join(parseErrors, "\n"))
	}
	return &lvl, nil
}

func (l *Level) BoardSize() (int, int) {
	return l.BoardWidth, l.BoardHeigth
}

var ptn = regexp.MustCompile(`(?m:^[a-zA-Z]+:)`)

type keyVal struct{ k, v string }

func parseString(s string) []keyVal {
	if s == "" {
		return nil
	}
	var (
		kvs        []keyVal
		lastKey    string
		start, end int
	)
	for _, keyPos := range ptn.FindAllStringIndex(s, -1) {
		start = keyPos[0]
		if end != start {
			kvs = append(kvs, keyVal{lastKey, s[end:start]})
		}
		end = keyPos[1]
		lastKey = s[start : end-1]
	}
	return append(kvs, keyVal{lastKey, s[end:]})
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}
