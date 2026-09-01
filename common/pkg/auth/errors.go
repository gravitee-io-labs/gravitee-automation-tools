package auth

const NoAuthProvided = ClientAuthError("no auth configured: provide basic auth or bearer auth credentials")

type ClientAuthError string

func (e ClientAuthError) Error() string {
	return string(e)
}
