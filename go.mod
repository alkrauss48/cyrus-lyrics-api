module github.com/alkrauss48/cyrus-lyrics-api

go 1.16

require (
	github.com/alkrauss48/cyrus-lyrics-api/oauth v0.0.0
	github.com/alkrauss48/cyrus-lyrics-api/public v0.0.0
	github.com/alkrauss48/cyrus-lyrics-api/sheets v0.0.0
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
)

replace github.com/alkrauss48/cyrus-lyrics-api/helpers => ./helpers

replace github.com/alkrauss48/cyrus-lyrics-api/public => ./public

replace github.com/alkrauss48/cyrus-lyrics-api/oauth => ./oauth

replace github.com/alkrauss48/cyrus-lyrics-api/sheets => ./sheets
