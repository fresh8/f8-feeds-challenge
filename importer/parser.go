package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Parser describes a parser object which holds functions to parse values
type Parser struct {
}

// String parses a given value into string value
func (p *Parser) String(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	switch value.(type) {
	case float64:
		return strconv.Itoa(int(value.(float64)))
	case int:
		return fmt.Sprintf("%d", value.(int))
	}
	return ""
}

// Time parses a value into a given time form
func (p *Parser) Time(value interface{}) (time.Time, error) {
	if v, ok := value.(time.Time); ok {
		return v, nil
	}
	switch value.(type) {
	case float64:
		return time.Unix(int64(value.(float64)), 0).UTC(), nil
	case int:
		return time.Unix(int64(value.(int)), 0).UTC(), nil
	case string:
		input := value.(string)
		t, err := time.Parse(time.RFC3339, input)
		if err == nil {
			return t, nil
		}
		t, err = time.Parse("2006-01-02:15:00:00Z", input)
		if err == nil {
			return t, nil
		}
		t, err = time.Parse("2006-01-02 15:00:00Z", input)
		if err == nil {
			return t, nil
		}
		i, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			return time.Unix(int64(i), 0).UTC(), nil
		}
	}
	return time.Now(), fmt.Errorf("No time parsing worked")
}

// IntArray is trying to parse integer arrays
func (p *Parser) IntArray(value interface{}) ([]int, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var parsed []int
	if err = json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return parsed, err
}

// Events is to parse an array of events
func (p *Parser) Events(body []byte) ([]int, error) {
	var events []int
	err := json.Unmarshal(body, &events)
	return events, err
}

// Event is to parse a single event object
func (p *Parser) Event(body []byte) (*Event, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	t, err := p.Time(raw["time"])
	if err != nil {
		return nil, err
	}
	markets, err := p.IntArray(raw["markets"])
	if err != nil {
		return nil, err
	}
	event := Event{
		base: base{
			ID:   p.String(raw["id"]),
			Name: p.String(raw["name"]),
		},
		Time:    t,
		Markets: markets,
	}
	return &event, nil
}

// Market parses a market from the given body
func (p *Parser) Market(body []byte) (*Market, error) {
	return nil, nil
}

// Options is to parse options field out from a given raw data
func (p *Parser) Options(raw interface{}) ([]Option, error) {
	return nil, nil
}

// Odds parsing out odds from a given options raw object
func (p *Parser) Odds(raw interface{}) (int, int, error) {
	return 0, 0, nil
}
