package collymongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type request struct {
	RequestID uint64 `bson:"requestID"`
}

func (m *CollyMongo) IsVisited(requestID uint64) (visited bool, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), m.getWait())
	defer cancel()

	r := &request{
		RequestID: requestID,
	}

	if err = m.request.FindOne(ctx, r, m.findRequestOpts...).Decode(r); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (m *CollyMongo) Visited(requestID uint64) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), m.getWait())
	defer cancel()

	r := &request{
		RequestID: requestID,
	}

	if _, err = m.request.InsertOne(ctx, r, m.insertRequestOpts...); err != nil {
		return err
	}

	return
}
