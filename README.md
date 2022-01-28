[CyrusLyrics API](https://api.cyruskrauss.com)
==========

The API behind the [CyrusLyrics iOS app](https://github.com/alkrauss48/cyrus-lyrics),
written in Go.

The purpose of this API is to connect to the Google Sheets API and allow the iOS
app user to create their own spreadsheet of data for the application.

## To Run
```
cp .env.example .env

# Set the Google OAuth creds in the .env file
# These map to the keys in a Google OAuth credentials.json file

docker-compose up
```
