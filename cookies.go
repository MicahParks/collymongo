package collymongo

import (
	"context"
	"errors"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// cookie contains a hostname to cookie relationship.
type cookie struct {

	// Host is the hostname from the URL given by Colly.
	Host string `bson:"_id"`

	// Cookies is the string of cookies given by Colly.
	Cookies string `bson:"cookie"`
}

// Cookies follows the implementation of the storage.Storage interface. It accepts a URL and returns the cookies
// associated with the host, if any.
func (m *CollyMongo) Cookies(u *url.URL) (cookies string) {

	// Create a context with the appropriate timeout.
	ctx, cancel := context.WithTimeout(context.Background(), m.findWait())
	defer cancel()

	// Filter for cookies by the given URL's host.
	c := &cookie{
		Host: u.Host,
	}

	// Ask MongoDB for said host's cookies. Put them into the cookie struct.
	if err := m.cookie.FindOne(ctx, c, m.FindCookieOpts...).Decode(c); err != nil {

		// Should an error occur that wasn't the fact that the hostname hasn't been seen, log it.
		if !errors.Is(err, mongo.ErrNoDocuments) {
			m.log(err)
		}

		// A document with this hostname was not found.
		return ""
	}

	return c.Cookies
}

// SetCookies follows the implementation of the storage.Storage interface. It takes in a cookies and their URL then
// stores the URL's hostname and cookies in a document in MongoDB.
func (m *CollyMongo) SetCookies(u *url.URL, cookies string) {

	// Create a context with the appropriate timeout.
	ctx, cancel := context.WithTimeout(context.Background(), m.insertWait())
	defer cancel()

	// Create the cookie.
	c := &cookie{
		Cookies: cookies,
		Host:    u.Host,
	}

	// Create the filter.
	filter := &cookie{
		Host: c.Host,
	}

	// Upsert the document should it not already be there.
	opts := m.ReplaceCookieOpts
	if opts == nil {
		opts = make([]*options.FindOneAndReplaceOptions, 0)
	}
	opts = append(opts, options.FindOneAndReplace().SetUpsert(true))

	// The old document before it was replaced.
	var replacement *cookie

	// Insert the cookie.
	if err := m.cookie.FindOneAndReplace(ctx, c, filter, opts...).Decode(replacement); err != nil {

		// If the document was upserted, ignore the error.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return
		}

		// Log the error, if present.
		m.log(err)
	}
}
