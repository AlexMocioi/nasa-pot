# nasa-pot
Proof of Technology Go with Nasa API

Collect and format data from https://api.nasa.gov/ Asteroids - NeoWs

#config
Change values in config.yaml
    services.api-key to not use the default one anymore
    slackHookURL if you have a slack channel to post messages on it
    
# run
`go run main.go --config config.yaml`

# test
Access `127.0.0.1:8089/grabLatest/2020-05-18/sync` to call live the NASA endpoint and see a summary result.
The handler will also trigger a new entry saving in big cache.

Access `127.0.0.1:8089/grabLatest/2020-05-18/async` to get data from cache if avaialble.

