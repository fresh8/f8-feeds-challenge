package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {
	parser := Parser{}

	table := []struct {
		Input  []byte
		Output []int
		Fail   bool
	}{
		{
			Input:  []byte(`[1,2,3,4]`),
			Output: []int{1, 2, 3, 4},
			Fail:   false,
		},
		{
			Input:  []byte(`[1,2,3,5]`),
			Output: []int{1, 2, 3, 4},
			Fail:   true,
		},
	}

	for _, data := range table {
		var fail bool
		parsed, err := parser.Events(data.Input)
		if err != nil {
			t.Error(err)
		}
		for i, num := range data.Output {
			if parsed[i] != num {
				fail = true
				break
			}
		}

		// If it failed and it should fail continue
		if fail == data.Fail {
			continue
		}
		fmt.Printf("Failed: || input %v | output %v\n", parsed, data.Output)
		t.Fail()
	}
}

// Checks if an event has the same data as another event
func isEventSame(a Event, b Event) bool {
	if a.base != b.base {
		return false
	}
	if a.Time != b.Time {
		fmt.Println(a.Time, b.Time)
		return false
	}
	for id, market := range a.Markets {
		if b.Markets[id] != market {
			fmt.Println(market, b.Markets[id])
			return false
		}
	}
	return true
}

// Checks if an options has the same data as another options
func isOptionSame(a Option, b Option) bool {
	if a.base != b.base {
		return false
	}
	if a.Numerator != b.Numerator {
		return false
	}
	if a.Denominator != b.Denominator {
		return false
	}
	return true
}

func isMarketSame(a Market, b Market) bool {
	if a.ID != b.ID {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	for _, option := range a.Options {
		var include bool
		for _, other := range b.Options {
			if isOptionSame(option, other) {
				include = true
			}
		}
		if include {
			return true
		}
	}
	return false
}

func TestEvent(t *testing.T) {
	parser := Parser{}
	cTime, _ := time.Parse("2006-01-02:15:00:00Z", "2017-08-20:15:00:00Z")
	table := []struct {
		Description string
		Input       []byte
		Output      Event
		Fail        bool
	}{
		{
			Input: []byte(`{ "id": 1, "name": "Southampton v Bournemouth", "time": "2017-08-20:15:00:00Z", "markets": [ 101, 102 ] }`),
			Output: Event{
				base: base{
					ID:   "1",
					Name: "Southampton v Bournemouth",
				},
				Time:    cTime,
				Markets: []int{101, 102},
			},
			Fail:        false,
			Description: "Should be parsed with integer as ID",
		},
		{
			Input: []byte(`{ "id": "1", "name": "Southampton v Bournemouth", "time": "2017-08-20:15:00:00Z", "markets": [ 101, 102 ] }`),
			Output: Event{
				base: base{
					ID:   "1",
					Name: "Southampton v Bournemouth",
				},
				Time:    cTime,
				Markets: []int{101, 102},
			},
			Fail:        false,
			Description: "Should be parsed with string as ID",
		},
		{
			Input: []byte(`{ "id": "1", "name": "Southampton v Bournemouth", "time": "2017-08-20 15:00:00Z", "markets": [ 101, 102 ] }`),
			Output: Event{
				base: base{
					ID:   "1",
					Name: "Southampton v Bournemouth",
				},
				Time:    cTime,
				Markets: []int{101, 102},
			},
			Fail:        false,
			Description: "Should be parsed with different time format",
		},
		{
			Input: []byte(`{ "id": "1", "name": "Southampton v Bournemouth", "time": 1503241200, "markets": [ 101, 102 ] }`),
			Output: Event{
				base: base{
					ID:   "1",
					Name: "Southampton v Bournemouth",
				},
				Time:    cTime,
				Markets: []int{101, 102},
			},
			Fail:        false,
			Description: "Should be parsed with unix time as time format",
		},
		{
			Input: []byte(`{ "id": "1", "name": "Southampton v Bournemouth", "time": "1503241200", "markets": [ 101, 102 ] }`),
			Output: Event{
				base: base{
					ID:   "1",
					Name: "Southampton v Bournemouth",
				},
				Time:    cTime,
				Markets: []int{101, 102},
			},
			Fail:        false,
			Description: "Should be parsed with unix time as time format",
		},
	}
	for _, data := range table {
		event, err := parser.Event(data.Input)
		if err != nil {
			t.Error(err)
		}
		if isEventSame(*event, data.Output) && !data.Fail {
			continue
		}
		t.Fail()
	}
}

func TestOption(t *testing.T) {
	parser := Parser{}
	table := []struct {
		Description   string
		Input         []byte
		Output        []Option
		Fail          bool
		ExpectedError error
	}{
		{
			Description: "Should not fail, input correctly formatted",
			Input:       []byte(`{"options": [ { "id": "10101", "name": "Southampton", "odds": "3/5" }, { "id": "10102", "name": "Draw", "odds": "4/5" }, { "id": "10103", "name": "Bournemouth", "odds": "5/1" } ]}`),
			Output: []Option{
				{
					base: base{
						ID:   "10101",
						Name: "Southampton",
					},
					Numerator:   3,
					Denominator: 5,
				},
				{
					base: base{
						ID:   "10102",
						Name: "Draw",
					},
					Numerator:   4,
					Denominator: 5,
				},
				{
					base: base{
						ID:   "10103",
						Name: "Bournemouth",
					},
					Numerator:   5,
					Denominator: 1,
				},
			},
			Fail: false,
		},
		{
			Description:   "Should fail, input odds is not correctly formatted",
			Input:         []byte(`{"options": [ { "id": "10101", "name": "Southampton", "odds": 6 }, { "id": "10102", "name": "Draw", "odds": "4/5" }, { "id": "10103", "name": "Bournemouth", "odds": "5/1" } ]}`),
			Output:        []Option{},
			Fail:          true,
			ExpectedError: ErrOdds,
		},
	}

	var options []Option
	var err error
	for _, data := range table {
		var raw map[string]interface{}
		if err = json.Unmarshal(data.Input, &raw); err != nil {
			t.Error(err)
		}
		options, err = parser.Options(raw)
		if err != nil && data.ExpectedError == nil {
			t.Error(err)
		}
		if data.ExpectedError != nil && data.ExpectedError != err {
			t.Error(data.ExpectedError, err)
		}
		for i, option := range options {
			if isOptionSame(option, data.Output[i]) && !data.Fail {
				continue
			}
			t.Fail()
		}
	}
}

func TestOdds(t *testing.T) {
	type output struct {
		Numerator   int
		Denominator int
		Err         error
	}

	parser := Parser{}
	table := []struct {
		Input  []byte
		Output output
	}{
		{
			Input: []byte(`{ "odds": "2/5" }`),
			Output: output{
				Numerator:   2,
				Denominator: 5,
			},
		},
		{
			Input: []byte(`{ "odds": "13/4" }`),
			Output: output{
				Numerator:   13,
				Denominator: 4,
			},
		},
		{
			Input: []byte(`{ "odds": 6 }`),
			Output: output{
				Numerator:   2,
				Denominator: 5,
				Err:         ErrOdds,
			},
		},
	}

	for _, data := range table {
		var raw map[string]interface{}
		if err := json.Unmarshal(data.Input, &raw); err != nil {
			t.Error(err)
		}
		num, den, err := parser.Odds(raw)
		if err != nil {
			if err == data.Output.Err {
				continue
			}
			t.Error(err)
		}

		if num != data.Output.Numerator || den != data.Output.Denominator {
			t.Fail()
		}
	}
}
func TestMarket(t *testing.T) {
	type output struct {
		Numerator   int
		Denominator int
		Err         error
	}

	parser := Parser{}
	table := []struct {
		Input  []byte
		Output Market
	}{
		{
			Input: []byte(`{ "id": "101", "type": "win-draw-win", "options": [ { "id": "10101", "name": "Southampton", "odds": "3/5" }, { "id": "10102", "name": "Draw", "odds": "4/5" }, { "id": "10103", "name": "Bournemouth", "odds": "5/1" } ] }`),
			Output: Market{
				ID:   "101",
				Type: "win-draw-win",
				Options: []Option{
					{
						base: base{
							ID:   "10101",
							Name: "Southampton",
						},
						Numerator:   3,
						Denominator: 5,
					},
					{
						base: base{
							ID:   "10102",
							Name: "Draw",
						},
						Numerator:   4,
						Denominator: 5,
					},
					{
						base: base{
							ID:   "10103",
							Name: "Bournemouth",
						},
						Numerator:   5,
						Denominator: 1,
					},
				},
			},
		},
	}

	for _, data := range table {
		market, err := parser.Market(data.Input)
		if err != nil {
			t.Error(err)
		}
		if market == nil {
			t.Fail()
		}
		if isMarketSame(*market, data.Output) {
			continue
		}
		t.Fail()
	}
}
