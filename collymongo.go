// collymongo is a MongoDB storage backend for the Colly framework.
package collymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (

	// defaultCookie is the default name of the collection to store the hostname cookie relationship in MongoDB.
	defaultCookie = "cookie"

	// defaultDatabase is the default name for the database to store Colly persistent data.
	defaultDatabase = "colly"

	// defaultRequest is the default name of the collection to store request IDs in MongoDB.
	defaultRequest = "request"

	// defaultWait is the amount of time to wait for a call to MongoDB if none was given by the package user.
	defaultWait = time.Second * 10
)

// CollyMongo implements the Storage interface from github.com/gocolly/colly/storage. It allows for the use of MongoDB
// as a storage backend for Colly's host to cookie relationships and previous request IDs.
type CollyMongo struct {

	// ErrCookie allows the package user to find out why the call to MongoDB regarding cookies failed. The interface
	// does not return an error.
	ErrCookie error

	// client is the MongoDB client.
	client *mongo.Client

	// client options are the options used to create the client for MongoDB.
	clientOptions *options.ClientOptions

	// cookie is the collection in MongoDB used to store host to cookie relationships.
	cookie *mongo.Collection

	// cookieCol is the name of the collection in MongoDB used to store host to cookie relationships. It defaults to
	// "cookie".
	cookieCol string

	// cookieOpts are the collection options to pass to the mongo package when creating a collection struct.
	cookieOpts []*options.CollectionOptions

	// database is the database to store Colly's host to cookie relationships and determining if a request has already
	// been made.
	database *mongo.Database

	// databaseName is the name of the database to store Colly's host to cookie relationships and determining if a
	// request has already been made. It defaults to "colly".
	databaseName string

	// databaseOpts are the database options to pass to the mongo package when creating a database struct.
	databaseOpts []*options.DatabaseOptions

	// findCookieOpts are the options to pass to the mongo package when finding one cookie by hostname.
	findCookieOpts []*options.FindOneOptions

	// findCtxTime is the amount of time to put in the context timeout that is passed to the mongo package for finding
	// documents.
	findCtxTime time.Duration

	// findRequestOpts are the options to pass to the mongo package when finding request IDs that have already been
	// performed.
	findRequestOpts []*options.FindOneOptions

	// initCtxTime is the amount of time to put in the context timeout that is passed to the mongo package when
	// initializing that database connection.
	initCtxTime time.Duration

	// insertCookieOpts are the options to pass to the mongo package when inserting a cookie into the cookie collection.
	insertCookieOpts []*options.InsertOneOptions

	// insertCtxTime is the amount of time to put in the context timeout that is passed to the mongo package when
	// inserting a document.
	insertCtxTime time.Duration

	// insertRequestOpts are the options to pass to the mongo package when inserting a request ID into the request
	// collection.
	insertRequestOpts []*options.InsertOneOptions

	// uri is the MongoDB URI string.
	uri string

	// request is the collection to hold previous request IDs.
	request *mongo.Collection

	// requestCol is the name of the collection in MongoDB used to store host to request IDs. It defaults to "request".
	requestCol string

	// requestOpts are the options to pass the mongo package when creating the request collection.
	requestOpts []*options.CollectionOptions
}

// findWait determines the proper amount of time to wait for a request to find a document in MongoDB.
func (m *CollyMongo) findWait() time.Duration {

	// Copy the user given time.
	wait := m.findCtxTime

	// If that time was left zero, use the default time.
	if wait == 0 {
		wait = defaultWait
	}

	return wait
}

// insertWait determines the proper amount of time to wait for a request to insert a document into MongoDB.
func (m *CollyMongo) insertWait() time.Duration {

	// Copy the user given time.
	wait := m.insertCtxTime

	// If that time was left zero, use the default time.
	if wait == 0 {
		wait = defaultWait
	}

	return wait
}
