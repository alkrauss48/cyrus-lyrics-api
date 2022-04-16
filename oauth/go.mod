module github.com/alkrauss48/cyrus-lyrics-api/oauth

go 1.16

require (
	github.com/alkrauss48/cyrus-lyrics-api/helpers v0.0.0
	github.com/gin-gonic/gin v1.7.7
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
)

replace github.com/alkrauss48/cyrus-lyrics-api/helpers => ../helpers
