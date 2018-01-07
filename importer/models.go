package main

import "time"

type base struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Event describes a football event
type Event struct {
	base
	Time    time.Time `json:"-"`
	Markets []int     `json:"markets"`
}

// PopulatedEvent are events, which are already populated and ready to be
// sent to the store
type PopulatedEvent struct {
	Event
	TimeRepresentation string   `json:"time"`
	Markets            []Market `json:"markets"`
}

// Market describes a football market
type Market struct {
	ID      string   `json:"ID"`
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

// rawOption describes an option where its odds represented in a string
type rawOption struct {
	base
	Odds string
}

// Option describes a football market's option
type Option struct {
	base
	Numerator   int `json:"num"`
	Denominator int `json:"den"`
}
