package collymongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init follows the implementation of the storage.Storage interface. It will create a client and connect to MongoDB. It
// will also create the database, cookie, and request collections if they do not already exist.
func (m *CollyMongo) Init() (err error) {

	// If no client options were given, use the default.
	if m.clientOptions == nil {
		m.clientOptions = options.Client()
	}

	// If a URI was given, apply it to the client options.
	if len(m.uri) != 0 {
		m.clientOptions = m.clientOptions.ApplyURI(m.uri)
	}

	// Create the MongoDB client.
	if m.client, err = mongo.NewClient(m.clientOptions); err != nil {
		return err
	}

	// Initialize some context related variables.
	var ctx context.Context
	var cancel context.CancelFunc

	// If no context timeout was given for initialization, use the default. Create the context.
	if m.initCtxTime == 0 {
		ctx, cancel = context.WithTimeout(context.Background(), defaultWait)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), m.initCtxTime)
	}
	defer cancel()

	// Connect to the MongoDB server.
	if err = m.client.Connect(ctx); err != nil {
		return err
	}

	// If no database name was given, use the default.
	if len(m.databaseName) == 0 {
		m.databaseName = defaultDatabase
	}

	// Create the database struct.
	m.database = m.client.Database(m.databaseName, m.databaseOpts...)

	// If no cookie collection name was given, use the default.
	if len(m.cookieCol) == 0 {
		m.cookieCol = defaultCookie
	}

	// Create the cookie collection struct.
	m.cookie = m.database.Collection(m.cookieCol, m.cookieOpts...)

	// If no request collection name was give, use the default.
	if len(m.requestCol) == 0 {
		m.requestCol = defaultRequest
	}

	// Create the request collection struct.
	m.request = m.database.Collection(m.requestCol, m.requestOpts...)

	return nil
}
