package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
)

const baseUrl string = "https://api.twitter.com"


type Configuration struct {
    Key string
    Secret string
}

func getConfig(configFile string) Configuration {

    file, _ := os.Open(configFile)
    contents, _ := ioutil.ReadAll(file)

    var config Configuration
    json.Unmarshal(contents, &config)

    return config
}

func generateBearerCreds(apiKey string, apiSecret string) string {

    /* generate the base64 encoded credentials necessary to request an app
    bearer token */

    var bearerCreds = apiKey + ":" + apiSecret
    // []byte(bearerCreds) casts string to bytes
    var encodedCreds = base64.StdEncoding.EncodeToString([]byte(bearerCreds))
    return encodedCreds
}

func extractToken(response *http.Response) string {

    /* extract a bearer token from a JSON-formatted response from the
    oauth endpoint */

    body, _ := ioutil.ReadAll(response.Body)

    type Response struct {
        Access_token string
        Token_type string
    }

    var responseBody Response
    json.Unmarshal(body, &responseBody)

    var bearerToken = responseBody.Access_token
    return bearerToken
}

func getBearerToken(bearerCreds string) string {

    /* retrieve a bearer token from the oauth endpoint using base64
    encoded bearer credentials*/

    u, _ := url.ParseRequestURI(baseUrl)
    u.Path = "oauth2/token"
    var postUrl = fmt.Sprintf("%v", u)

    var postType = "application/x-www-form-urlencoded"
    var data = url.Values{}
    data.Set("grant_type", "client_credentials")

    var client = &http.Client{}

    request, _ := http.NewRequest("POST", postUrl,
        bytes.NewBufferString(data.Encode()))

    request.Header.Add("Authorization", "Basic " + bearerCreds)
    request.Header.Add("Content-Type", postType)


    response, _ := client.Do(request)

    return extractToken(response)
}

func extractTweet(response *http.Response) string {

    /* extract a tweet's text from a JSON-formatted response from the
    statuses/user_timeline endpoint */

    body, _ := ioutil.ReadAll(response.Body)

    type Tweet struct {
        Text string
    }

    var tweetList = []Tweet{}

    json.Unmarshal(body, &tweetList)

    var firstTweet = tweetList[0]

    return firstTweet.Text
}


func getLastTweet(bearerToken string, username string) string {

    /* retrieve a tweet's text from the statuses/user_timeline endpoint using
    an oauth bearer token*/

    u, _ := url.ParseRequestURI(baseUrl)
    u.Path = "1.1/statuses/user_timeline.json"
    var queryParams = url.Values{}
    queryParams.Set("count", "1")
    queryParams.Set("screen_name", username)
    u.RawQuery = queryParams.Encode()

    var url = fmt.Sprintf("%v", u)

    var client = &http.Client{}

    request, _ := http.NewRequest("GET", url, nil)
    request.Header.Add("Authorization", "Bearer " + bearerToken)
    response, _ := client.Do(request)

    return extractTweet(response)
}

func main() {

    /* given a twitter screen name, return the text of the last tweet
    for said user from the twitter API*/

    // extract the username from command line arguments
    var username = flag.String("u", "Jwpe", "Twitter Screen Name")
    var configFile = flag.String("c", "config.json", "Path to config file")
    flag.Parse()

    var config = getConfig(*configFile)
    var bearerCreds string = generateBearerCreds(config.Key, config.Secret)
    var bearerToken string = getBearerToken(bearerCreds)
    var tweet string = getLastTweet(bearerToken, *username)

    fmt.Println(tweet)
}