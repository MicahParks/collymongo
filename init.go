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
	if m.ClientOptions == nil {
		m.ClientOptions = options.Client()
	}

	// If a URI was given, apply it to the client options.
	if len(m.Uri) != 0 {
		m.ClientOptions = m.ClientOptions.ApplyURI(m.Uri)
	}

	// Create the MongoDB client.
	if m.client, err = mongo.NewClient(m.ClientOptions); err != nil {
		return err
	}

	// Initialize some context related variables.
	var ctx context.Context
	var cancel context.CancelFunc

	// If no context timeout was given for initialization, use the default. Create the context.
	if m.InitCtxTime == 0 {
		ctx, cancel = context.WithTimeout(context.Background(), DefaultWait)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), m.InitCtxTime)
	}
	defer cancel()

	// Connect to the MongoDB server.
	if err = m.client.Connect(ctx); err != nil {
		return err
	}

	// If no database name was given, use the default.
	if len(m.DatabaseName) == 0 {
		m.DatabaseName = DefaultDatabase
	}

	// Create the database struct.
	m.database = m.client.Database(m.DatabaseName, m.DatabaseOpts...)

	// If no cookie collection name was given, use the default.
	if len(m.CookieCol) == 0 {
		m.CookieCol = DefaultCookie
	}

	// Create the cookie collection struct.
	m.cookie = m.database.Collection(m.CookieCol, m.CookieOpts...)

	// If no request collection name was give, use the default.
	if len(m.RequestCol) == 0 {
		m.RequestCol = DefaultRequest
	}

	// Create the request collection struct.
	m.request = m.database.Collection(m.RequestCol, m.RequestOpts...)

	return nil
}
