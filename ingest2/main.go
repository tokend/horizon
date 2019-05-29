// Package ingest2 provides tools which are used to convert core transactions into primitives which are more useful
// for the client side applications. Main entry points are `producer` - fetches data from the core and
// `consumer` - handles the data and stores it into horizon db
package ingest2

// increase it if you want force reingest (after backward not compatible changes)
const CurrentIngestVersion = 3
