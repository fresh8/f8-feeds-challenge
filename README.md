# Feeds Coding Challenge

* Write in Golang
* Use any libraries but please justify their use
* Use any dependancy management tool (if required)
* Things to watch out for: Concurrency, Messaging, Approach

## The challenge

Using the fake feed, import the events and markets and send them to the fake
store. Along the way you will encounter some inconsistent data so have a think
about the best way to handle that.

## Getting Started

* Ensure you have golang installed (use the latest version)
* Run the feed server `make run-feed`
* Query the event feed to get all of the events, along with their corresponding
  markets
* Validate the data
* POST the events to the feed store

## Feed Server

Check out the [documentation](feed/README.md)

## Models

### Event

All fields are required

ID: string
Name: string
Time: timestamp (RFC 3339)
Markets: array of markets (at least 1)

### Market

All fields are required

ID: string
Type: string
Options: array of options (at least 1)

### Options

All fields are required

ID: string
Name: string
Numerator: int
Denomenator: int
