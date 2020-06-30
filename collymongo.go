package collymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultWait = time.Second * 10
)

type CollyMongo struct {
	CookieErr         error
	client            *mongo.Client
	cookie            *mongo.Collection
	cookieCol         string
	cookieOpts        []*options.CollectionOptions
	database          *mongo.Database
	databaseName      string
	databaseOpts      []*options.DatabaseOptions
	findCookieOpts    []*options.FindOneOptions
	findCtxTime       time.Duration
	findRequestOpts   []*options.FindOneOptions
	initCtxTime       time.Duration
	insertCookieOpts  []*options.InsertOneOptions
	insertCtxTime     time.Duration
	insertRequestOpts []*options.InsertOneOptions
	options           *options.ClientOptions
	uri               string
	request           *mongo.Collection
	requestCol        string
	requestOpts       []*options.CollectionOptions
}

func (m *CollyMongo) getWait() time.Duration {
	wait := m.findCtxTime
	if wait == 0 {
		wait = defaultWait
	}
	return wait
}
