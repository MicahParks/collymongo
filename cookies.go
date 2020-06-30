package collymongo

import (
	"context"
	"errors"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
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

	// Insert the cookie.
	if _, err := m.cookie.InsertOne(ctx, c, m.InsertCookieOpts...); err != nil {

		// Log the error, if present.
		m.log(err)
	}
}
