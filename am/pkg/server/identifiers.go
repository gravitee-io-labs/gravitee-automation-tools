package server

func (d Domain) Identity() string {
	return d.Key
}

func (c Certificate) Identity() string {
	return c.Key
}

func (i IdentityProvider) Identity() string {
	return i.Key
}

func (r Reporter) Identity() string {
	return r.Key
}
