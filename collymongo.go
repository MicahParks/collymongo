// collymongo is a MongoDB storage backend for the Colly framework.
package collymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (

	// DefaultCookie is the default name of the collection to store the hostname cookie relationship in MongoDB.
	DefaultCookie = "cookie"

	// DefaultDatabase is the default name for the database to store Colly persistent data.
	DefaultDatabase = "colly"

	// DefaultRequest is the default name of the collection to store request IDs in MongoDB.
	DefaultRequest = "request"

	// DefaultWait is the amount of time to wait for a call to MongoDB if none was given by the package user.
	DefaultWait = time.Second * 10
)

// CollyMongo implements the Storage interface from github.com/gocolly/colly/storage. It allows for the use of MongoDB
// as a storage backend for Colly's host to cookie relationships and previous request IDs.
type CollyMongo struct {

	// client is the MongoDB client.
	client *mongo.Client

	// ClientOptions are the options used to create the client for MongoDB.  Can safely be left blank.
	ClientOptions *options.ClientOptions

	// cookie is the collection in MongoDB used to store host to cookie relationships.
	cookie *mongo.Collection

	// CookieCol is the name of the collection in MongoDB used to store host to cookie relationships. It defaults to
	// "cookie".  Can safely be left blank.
	CookieCol string

	// CookieOpts are the collection options to pass to the mongo package when creating a collection struct. Can safely
	// left left blank.
	CookieOpts []*options.CollectionOptions

	// database is the database to store Colly's host to cookie relationships and determining if a request has already
	// been made.
	database *mongo.Database

	// DatabaseName is the name of the database to store Colly's host to cookie relationships and determining if a
	// request has already been made. It defaults to "colly".  Can safely be left blank.
	DatabaseName string

	// DatabaseOpts are the database options to pass to the mongo package when creating a database struct. Can safely
	// left left blank.
	DatabaseOpts []*options.DatabaseOptions

	// ErrCookie allows the package user to find out why the call to MongoDB regarding cookies failed. The interface
	// does not return an error.
	ErrCookie error

	// FindCookieOpts are the options to pass to the mongo package when finding one cookie by hostname. Can safely left
	// left blank.
	FindCookieOpts []*options.FindOneOptions

	// FindCtxTime is the amount of time to put in the context timeout that is passed to the mongo package for finding
	// documents. Defaults to 10 seconds. Can safely be left blank.
	FindCtxTime time.Duration

	// FindRequestOpts are the options to pass to the mongo package when finding request IDs that have already been
	// performed. Can safely be left blank.
	FindRequestOpts []*options.FindOneOptions

	// InitCtxTime is the amount of time to put in the context timeout that is passed to the mongo package when
	// initializing that database connection. Defaults to 10 seconds. Can safely be left blank.
	InitCtxTime time.Duration

	// InsertCookieOpts are the options to pass to the mongo package when inserting a cookie into the cookie collection.
	// Can safely be left blank.
	InsertCookieOpts []*options.InsertOneOptions

	// InsertCtxTime is the amount of time to put in the context timeout that is passed to the mongo package when
	// inserting a document. Defaults to 10 seconds. Can safely be left blank.
	InsertCtxTime time.Duration

	// InsertRequestOpts are the options to pass to the mongo package when inserting a request ID into the request
	// collection. Can safely be left blank.
	InsertRequestOpts []*options.InsertOneOptions

	// Uri is the MongoDB URI string. Mandatory.
	Uri string

	// request is the collection to hold previous request IDs.
	request *mongo.Collection

	// RequestCol is the name of the collection in MongoDB used to store host to request IDs. It defaults to "request".
	// Can safely be left blank.
	RequestCol string

	// RequestOpts are the options to pass the mongo package when creating the request collection.  Can safely be left
	// blank.
	RequestOpts []*options.CollectionOptions
}

// findWait determines the proper amount of time to wait for a request to find a document in MongoDB.
func (m *CollyMongo) findWait() time.Duration {

	// Copy the user given time.
	wait := m.FindCtxTime

	// If that time was left zero, use the default time.
	if wait == 0 {
		wait = DefaultWait
	}

	return wait
}

// insertWait determines the proper amount of time to wait for a request to insert a document into MongoDB.
func (m *CollyMongo) insertWait() time.Duration {

	// Copy the user given time.
	wait := m.InsertCtxTime

	// If that time was left zero, use the default time.
	if wait == 0 {
		wait = DefaultWait
	}

	return wait
}
