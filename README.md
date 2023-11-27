# Geocoder: City of Toronto

Geocoding is the act of taking an address and returning its coordinates, its
longitude and latitude.

    >> 299 Queen St W
    Location: [-79.39038, 43.64955]

This geocoder supports both street addresses and intersections in the city of
Toronto.

    >> College St / Spadina Ave
    Location: [-79.40005, 43.65795]

## Why build this?

Most geocoding services available on the web are cost prohibitive for offline
analysis of large datasets. For instance, Google's Geocoding API has a [USD$4
cost](https://developers.google.com/maps/documentation/geocoding/usage-and-billing)
for requests even at high volumes.

For instance, a municipal open data set like the City of Toronto's [record of 
parking tickets](https://open.toronto.ca/dataset/parking-tickets/) would be 
very costly to geocode as it contains many records located at thousands of
addresses across the city from many prior years.

With this open source software, that work is free.

## Building the index

The index is built live from real municipal open datasets each time the
indexer is run.

The resulting index is saved on disk in recordIO format using this protocol
buffer for serialization.

    %> go run indexer/main.go
    Working.....
    Wrote ./toronto_geocode.dat

## Usage

First, start the service:

    %> go run indexer/main.go
    Loading index...Ready!
    >> 

The geocoding service is available both on the command line and as a GRPC
service. On the command line, simply enter your query and get a result:

    >> 299 Queen St W
    Location: [-79.39038, 43.64955]

On the RPC interface, send a query and get a response at an extremely high
speed.
