package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

// Token represents token response
type Token struct {
	Value string `json:"access_token"`
}

// Site holds a site source key
type Site struct {
	SourceKey string `json:"sourceKey"`
}

// SiteCollection holds a list of sites
type SiteCollection struct {
	Sites []Site `json:"content"`
}

// Asset holds an asset source key
type Asset struct {
	SourceKey string `json:"sourceKey"`
}

// AssetCollection holds a list of assets
type AssetCollection struct {
	Assets []Asset `json:"content"`
}

var (
	uaaURL              = "https://cc076c41-1318-4aab-a64e-18aa6dd254b7.predix-uaa.run.aws-usw02-pr.ice.predix.io/oauth/token"
	username            = "ems-apm-admin2"
	password            = "se3ret"
	timeout             = 5
	enterpriseSourceKey = "ENTERPRISE_da4ab60d-2f69-4bdb-af18-6cafe981af82"
	sitesAPITmpl        = "http://localhost:8008/v1/enterprises/%s/sites"
	assetsAPITmpl       = "http://localhost:8008/v1/sites/%s/assets"
)

func getToken() string {
	var token Token

	request := gorequest.New()
	request.Timeout(time.Duration(timeout) * time.Second)
	request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request.SetBasicAuth(username, password)
	request.Post(uaaURL)
	request.Send("grant_type=client_credentials")
	_, _, errs := request.EndStruct(&token)
	if errs != nil {
		for _, err := range errs {
			log.Fatal(err)
		}
	}
	return token.Value
}

func getSites(token string) <-chan string {
	out := make(chan string)
	go func(t string) {
		sitesAPIEndpoint := fmt.Sprintf(sitesAPITmpl, enterpriseSourceKey)
		var siteCollection SiteCollection
		authorizationStr := fmt.Sprintf("Bearer %s", t)
		request := gorequest.New()
		request.Timeout(time.Duration(timeout) * time.Second)
		request.Get(sitesAPIEndpoint)
		request.Set("Accept", "application/json")
		request.Set("Authorization", authorizationStr)
		_, _, errs := request.EndStruct(&siteCollection)
		if errs != nil {
			for _, err := range errs {
				log.Fatal(err)
			}
		}
		for _, site := range siteCollection.Sites {
			out <- site.SourceKey
		}
		close(out)
	}(token)
	return out
}

func getAsset(token string, siteSourceKey string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var assetCollection AssetCollection
	authorizationStr := fmt.Sprintf("Bearer %s", token)
	request := gorequest.New()
	request.Timeout(time.Duration(timeout) * time.Second)
	assetsAPIEndpoint := fmt.Sprintf(assetsAPITmpl, siteSourceKey)
	request.Get(assetsAPIEndpoint)
	request.Set("Accept", "application/json")
	request.Set("Authorization", authorizationStr)
	_, _, errs := request.EndStruct(&assetCollection)
	if errs != nil {
		for _, err := range errs {
			log.Fatal(err)
		}
	}
	for _, asset := range assetCollection.Assets {
		out <- asset.SourceKey
	}
}

func getAssets(token string, in <-chan string) <-chan string {
	assetChan := make(chan string)
	wg := &sync.WaitGroup{}

	for siteSourceKey := range in {
		wg.Add(1)
		go getAsset(token, siteSourceKey, assetChan, wg)
	}

	go func() {
		wg.Wait()
		close(assetChan)
	}()

	return assetChan
}

func main() {

	token := getToken()
	c2 := getSites(token)
	c3 := getAssets(token, c2)

	for s := range c3 {
		fmt.Println(s)
	}
}
