package collymongo

import (
	"context"
	"net/url"
)

type cookie struct {
	Host    string `bson:"host"`
	Cookies string `bson:"cookie"`
}

func (m *CollyMongo) Cookies(u *url.URL) (cookies string) {

	ctx, cancel := context.WithTimeout(context.Background(), m.getWait())
	defer cancel()

	c := &cookie{
		Host: u.Host,
	}

	if err := m.cookie.FindOne(ctx, c, m.findCookieOpts...).Decode(c); err != nil {
		m.CookieErr = err
		return ""
	}

	return c.Cookies
}

func (m *CollyMongo) SetCookies(u *url.URL, cookies string) {

	ctx, cancel := context.WithTimeout(context.Background(), m.getWait())
	defer cancel()

	c := &cookie{
		Cookies: cookies,
		Host:    u.Host,
	}

	if _, err := m.cookie.InsertOne(ctx, c, m.insertCookieOpts...); err != nil {
		m.CookieErr = err
	}
}
