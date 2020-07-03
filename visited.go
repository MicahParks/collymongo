package collymongo

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// request simply holds a request ID.
type request struct {

	// RequestID is the ID of a request that has already been preformed.
	RequestID int64 `bson:"_id"`
}

// IsVisited follows the implementation of the storage.Storage interface. It checks if a request has been preformed yet.
func (m *CollyMongo) IsVisited(requestID uint64) (visited bool, err error) {

	// Create a context with the appropriate timeout.
	ctx, cancel := context.WithTimeout(context.Background(), m.findWait())
	defer cancel()

	// Filter for request ID matching the given.
	r := &request{

		// bson doesn't support uint64 and this will misrepresent the number that's in the database, but happens on both
		// ends of the code so it effectively works.
		RequestID: int64(requestID),
	}

	// Check to see if MongoDB has this request ID. Decode the response to a struct.
	if err = m.request.FindOne(ctx, r, m.FindRequestOpts...).Decode(r); err != nil {

		// If the document wasn't found, the given ID's request has not been preformed.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		// An error with MongoDB occurred.
		return false, err
	}

	// No errors means the ID was in MongoDB.
	return true, nil
}

// Visited follows the implementation of the storage.Storage interface. It tells MongoDB the ID of a request that has
// already been preformed.
func (m *CollyMongo) Visited(requestID uint64) (err error) {

	// Create a context with the appropriate timeout.
	ctx, cancel := context.WithTimeout(context.Background(), m.insertWait())
	defer cancel()

	// The request ID as a struct.
	r := &request{

		// bson doesn't support uint64 and this will misrepresent the number that's in the database, but happens on both
		// ends of the code so it effectively works.
		RequestID: int64(requestID),
	}

	// Insert the request ID into MongoDB.
	if _, err = m.request.InsertOne(ctx, r, m.InsertRequestOpts...); err != nil {

		// If it's a duplicate key, the it's already been inserted.
		if strings.Contains(err.Error(), "duplicate key") {
			return nil
		}

		return err
	}

	return nil
}
