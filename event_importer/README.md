# Event Importer

Imports all of the valid events from the feed API into a store.

Make sure that the STORE_ADDR environment variable is correctly set before running the program.

For example you might start the program with:

STORE_ADDR=localhost:8001 ./bin/event_importer

This would POST the events to http:// + STORE_ADDR + /event URL.

## To do

There are a few problems that still need to be resolved
- automatically formatting incorrect date formats
- the id fields are ints but the store specification expects id values of the string type
