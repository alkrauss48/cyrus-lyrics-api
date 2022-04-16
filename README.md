[CyrusLyrics API](https://api.cyruskrauss.com)
==========

The API behind the [CyrusLyrics iOS app](https://github.com/alkrauss48/cyrus-lyrics),
written in Go.

The purpose of this API is to connect to the Google Sheets API and allow the iOS
app user to create their own spreadsheet of data for the application.

This project uses the [Gin Web Framework](https://github.com/gin-gonic/gin).

## Before You Begin

This project requires the Google Sheets API. To use it, you will need to create
a set of OAuth credentials, as well as enable the Google Sheets API for those
credentials.

[More info on creating Google OAuth credentials here](
https://developers.google.com/workspace/guides/create-credentials#oauth-client-id)

## To Run
```
cp .env.example .env

# Next, set the PROJECT_ID; this corresponds to the Google Cloud project
# under which your OAuth credentials are created.
#
# Next, set the following Google OAuth creds in the .env file:
#
# CLIENT_ID
# CLIENT_SECRET
# REDIRECT_URI
#
# You will receive all of these from Google when creating OAuth credentials.

docker-compose up
```

## Available Routes

General Routes
```
GET /                   # Root route
GET /sheets/default     # List the publicly available default sheets
```

OAuth Routes
```
GET /oauth/google             # Initiate the Google OAuth login
GET /oauth/google/callback    # Complete the Google OAuth login
```

Authenticated Sheets Routes

Note: This app uses the **drive.file** Google Drive scope,
which allows access only to the files created under with this Google app.
```
GET /sheets/          # Get all sheet IDs and names
GET /sheets/:id       # Get a single sheet's data by ID
POST /sheets/         # Create a sheet, with a name
DELETE /sheets/:id    # Delete a sheet by ID
```
