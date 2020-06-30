package collymongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *CollyMongo) Init() (err error) {

	if m.clientOptions == nil {
		m.clientOptions = options.Client()
	}

	if len(m.uri) != 0 {
		m.clientOptions = m.clientOptions.ApplyURI(m.uri)
	}

	if m.client, err = mongo.NewClient(m.clientOptions); err != nil {
		return err
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if m.initCtxTime == 0 {
		ctx, cancel = context.WithTimeout(context.Background(), defaultWait)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), m.initCtxTime)
	}
	defer cancel()

	if err = m.client.Connect(ctx); err != nil {
		return err
	}

	if len(m.databaseName) == 0 {
		m.databaseName = "colly"
	}

	m.database = m.client.Database(m.databaseName, m.databaseOpts...)

	if len(m.cookieCol) == 0 {
		m.cookieCol = "cookie"
	}

	m.cookie = m.database.Collection(m.cookieCol, m.cookieOpts...)

	if len(m.requestCol) == 0 {
		m.requestCol = "request"
	}

	m.request = m.database.Collection(m.requestCol, m.requestOpts...)

	return nil
}
