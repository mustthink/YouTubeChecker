# YouTubeChecker

Checking new videos from tracked channels by query.

## Run
1) Set your configuration 
2) `go mod tidy`
3) `go run cmd/main.go`
>If this is your first run:
> 4) Get link from logs application
> 5) Follow the link 
> 6) Choose your account
> 7) Give the necessary permissions
> 8) Copy code and put it to application console
> 
> After first run you don't need to do this again.

## Requirements
1) Go version 1.21 or above

## Flags
1) `config` - path to config
2) `d` - activate debug mode

## Configuration 

### App config struct:
1) "interval_in_seconds" - interval between requests to YouTube
2) "caller" - caller configuration
   1) "count_result" - count of YouTube results in request
   2) "query" - query of content
   3) "api_key" - YouTube API key
3) "sheet" - sheet configuration
   1) "credentials_path" - path to credentials 
   2) "name" - sheet name,
   3) "id" - sheed ID

### Credentials config struct:
just default credentials file from Google
1) "installed"
   1) "client_id" - ID of your client
   2) "project_id" - ID of your project 
   3) "auth_uri" - Google Auth URI
   4) "token_uri" - Google token URI
   5) "auth_provider_x509_cert_url" - Google Auth Provider URL
   6) "client_secret" - client secret\
   7) "redirect_uris": - redirect uris

### Tracked channels
just string array
      