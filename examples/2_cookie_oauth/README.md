# 2_cookie_oauth

Let's say you want to build a web application with the following authorization flow:
- User wants to log in on your site via external OAuth provider
- User calls `auth.oauth()` method to get redirect URL to the OAuth provider
- User is asked to grant access to your application
- User is redirected back to your site with an authorization code, let's say to `/oauth/callback?code=1234`
- Your site sets a cookie with access token and redirects the user to the home page
- Now user is logged in and can access protected resources, such as `auth.status()`

This doesn't feel right for JSON RPC because of the cookies, but it's quite straightforward to implement, and tooling 
should be flexible enough to support it.

In this example we have single service with auth methods:

- `auth.oauth() -> OAuthResponse`
- `auth.status() -> AuthStatus`

Also, there is a separate handler for OAuth callback, `/oauth/callback`, which sets cookie on successful authentication.