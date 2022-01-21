package proto

//gjrpc:service auth
type Auth interface {
	//gjrpc:method oauth
	Oauth() (OAuthResponse, error)

	//gjrpc:method status
	Status() (AuthStatus, error)
}

type OAuthResponse struct {
	RedirectURL string
}

type AuthStatus struct {
	Email string
}
