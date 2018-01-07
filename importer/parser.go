package main

import "time"

// Parser describes a parser object which holds functions to parse values
type Parser struct {
}

// String parses a given value into string value
func (p *Parser) String(value interface{}) string {
	return ""
}

// Time parses a value into a given time form
func (p *Parser) Time(value interface{}) (time.Time, error) {
	return time.Now(), nil
}

// IntArray is trying to parse integer arrays
func (p *Parser) IntArray(value interface{}) ([]int, error) {
	return nil, nil
}

// Events is to parse an array of events
func (p *Parser) Events(body []byte) ([]int, error) {
	return nil, nil
}

// Event is to parse a single event object
func (p *Parser) Event(body []byte) (*Event, error) {
	return nil, nil
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
