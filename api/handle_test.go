package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/yukimochi/Activity-Relay/models"
)

const (
	BlockService models.Config = iota
	ManuallyAccept
	CreateAsAnnounce
)

func TestHandleWebfingerGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleWebfinger))
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL, nil)
	q := req.URL.Query()
	q.Add("resource", "acct:relay@"+GlobalConfig.ServerHostname().Host)
	req.URL.RawQuery = q.Encode()
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("fail - Content-Type is not match")
	}
	if r.StatusCode != 200 {
		t.Fatalf("fail - StatusCode is not match")
	}
	defer r.Body.Close()

	data, _ := io.ReadAll(r.Body)
	var webfinger models.WebfingerResource
	err = json.Unmarshal(data, &webfinger)
	if err != nil {
		t.Fatalf("fail - response is not valid")
	}

	domain, _ := url.Parse(webfinger.Links[0].Href)
	if domain.Host != GlobalConfig.ServerHostname().Host {
		t.Fatalf("fail - host is not match")
	}
}

func TestHandleWebfingerGetBadResource(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleWebfinger))
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL, nil)
	q := req.URL.Query()
	q.Add("resource", "acct:yukimochi@"+os.Getenv("RELAY_DOMAIN"))
	req.URL.RawQuery = q.Encode()
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 404 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleNodeinfoLinkGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleNodeinfoLink))
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("fail - Content-Type not match")
	}
	if r.StatusCode != 200 {
		t.Fatalf("fail - StatusCode is not match")
	}
	defer r.Body.Close()

	data, _ := io.ReadAll(r.Body)
	var nodeinfoLinks models.NodeinfoLinks
	err = json.Unmarshal(data, &nodeinfoLinks)
	if err != nil {
		t.Fatalf("fail - response is not valid")
	}
}

func TestHandleNodeinfoLinkInvalidMethod(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleNodeinfoLink))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleNodeinfoGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleNodeinfo))
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("fail - Content-Type is not match")
	}
	if r.StatusCode != 200 {
		t.Fatalf("fail - StatusCode is not match")
	}
	defer r.Body.Close()

	data, _ := io.ReadAll(r.Body)
	var nodeinfo models.Nodeinfo
	err = json.Unmarshal(data, &nodeinfo)
	if err != nil {
		t.Fatalf("fail - response is not valid")
	}
}

func TestHandleNodeinfoInvalidMethod(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleNodeinfo))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleWebfingerInvalidMethod(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleWebfinger))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleActorGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleActor))
	defer s.Close()

	r, err := http.Get(s.URL)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.Header.Get("Content-Type") != "application/activity+json" {
		t.Fatalf("fail - Content-Type not match")
	}
	if r.StatusCode != 200 {
		t.Fatalf("fail - StatusCode is not match")
	}
	defer r.Body.Close()

	data, _ := io.ReadAll(r.Body)
	var actor models.Actor
	err = json.Unmarshal(data, &actor)
	if err != nil {
		t.Fatalf("fail - response is not valid")
	}

	domain, _ := url.Parse(actor.ID)
	if domain.Host != GlobalConfig.ServerHostname().Host {
		t.Fatalf("fail - host is not match")
	}
}

func TestHandleActorInvalidMethod(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleActor))
	defer s.Close()

	r, err := http.Post(s.URL, "text/plain", nil)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestContains(t *testing.T) {
	data := "nil"
	sData := []string{
		"no",
		"nil",
	}
	invalidData := 0
	result := contains(data, "true")
	if result != false {
		t.Fatalf("fail - not contain but return true")
	}
	result = contains(data, "nil")
	if result != true {
		t.Fatalf("fail - contains but return false")
	}
	result = contains(sData, "true")
	if result != false {
		t.Fatalf("fail - not contain but return true (slice)")
	}
	result = contains(sData, "nil")
	if result != true {
		t.Fatalf("fail - contains but return false (slice)")
	}
	result = contains(invalidData, "hoge")
	if result != false {
		t.Fatalf("fail - given invalid data but return true (slice)")
	}
}

func mockActivityDecoderProvider(activity *models.Activity, actor *models.Actor) func(r *http.Request) (*models.Activity, *models.Actor, []byte, error) {
	return func(r *http.Request) (*models.Activity, *models.Actor, []byte, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		return activity, actor, body, nil
	}
}

func mockActivity(req string) models.Activity {
	switch req {
	case "Follow":
		file, _ := os.Open("../misc/test/follow.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	case "Invalid-Follow":
		file, _ := os.Open("../misc/test/followAsActor.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	case "Unfollow":
		file, _ := os.Open("../misc/test/unfollow.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	case "Invalid-Unfollow":
		body := "{\"@context\":\"https://www.w3.org/ns/activitystreams\",\"id\":\"https://mastodon.test.yukimochi.io/c125e836-e622-478e-a22d-2d9fbf2f496f\",\"type\":\"Undo\",\"actor\":\"https://mastodon.test.yukimochi.io/users/yukimochi\",\"object\":{\"@context\":\"https://www.w3.org/ns/activitystreams\",\"id\":\"https://hacked.test.yukimochi.io/c125e836-e622-478e-a22d-2d9fbf2f496f\",\"type\":\"Follow\",\"actor\":\"https://hacked.test.yukimochi.io/users/yukimochi\",\"object\":\"https://www.w3.org/ns/activitystreams#Public\"}}"
		var activity models.Activity
		json.Unmarshal([]byte(body), &activity)
		return activity
	case "UnfollowAsActor":
		body := "{\"@context\":\"https://www.w3.org/ns/activitystreams\",\"id\":\"https://mastodon.test.yukimochi.io/c125e836-e622-478e-a22d-2d9fbf2f496f\",\"type\":\"Undo\",\"actor\":\"https://mastodon.test.yukimochi.io/users/yukimochi\",\"object\":{\"@context\":\"https://www.w3.org/ns/activitystreams\",\"id\":\"https://hacked.test.yukimochi.io/c125e836-e622-478e-a22d-2d9fbf2f496f\",\"type\":\"Follow\",\"actor\":\"https://mastodon.test.yukimochi.io/users/yukimochi\",\"object\":\"https://relay.yukimochi.example.org/actor\"}}"
		var activity models.Activity
		json.Unmarshal([]byte(body), &activity)
		return activity
	case "Create":
		file, _ := os.Open("../misc/test/create.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	case "Create-Article":
		body := "{\"@context\":[\"https://www.w3.org/ns/activitystreams\",\"https://w3id.org/security/v1\",{\"manuallyApprovesFollowers\":\"as:manuallyApprovesFollowers\",\"sensitive\":\"as:sensitive\",\"movedTo\":{\"@id\":\"as:movedTo\",\"@type\":\"@id\"},\"Hashtag\":\"as:Hashtag\",\"ostatus\":\"http://ostatus.org#\",\"atomUri\":\"ostatus:atomUri\",\"inReplyToAtomUri\":\"ostatus:inReplyToAtomUri\",\"conversation\":\"ostatus:conversation\",\"toot\":\"http://joinmastodon.org/ns#\",\"Emoji\":\"toot:Emoji\",\"focalPoint\":{\"@container\":\"@list\",\"@id\":\"toot:focalPoint\"},\"featured\":{\"@id\":\"toot:featured\",\"@type\":\"@id\"},\"schema\":\"http://schema.org#\",\"PropertyValue\":\"schema:PropertyValue\",\"value\":\"schema:value\"}],\"id\":\"https://mastodon.test.yukimochi.io/users/yukimochi/statuses/101075045564444857/activity\",\"type\":\"Create\",\"actor\":\"https://mastodon.test.yukimochi.io/users/yukimochi\",\"published\":\"2018-11-15T11:07:26Z\",\"to\":[\"https://www.w3.org/ns/activitystreams#Public\"],\"cc\":[\"https://mastodon.test.yukimochi.io/users/yukimochi/followers\"],\"object\":{\"id\":\"https://mastodon.test.yukimochi.io/users/yukimochi/statuses/101075045564444857\",\"type\":\"Article\",\"summary\":null,\"inReplyTo\":null,\"published\":\"2018-11-15T11:07:26Z\",\"url\":\"https://mastodon.test.yukimochi.io/@yukimochi/101075045564444857\",\"attributedTo\":\"https://mastodon.test.yukimochi.io/users/yukimochi\",\"to\":[\"https://www.w3.org/ns/activitystreams#Public\"],\"cc\":[\"https://mastodon.test.yukimochi.io/users/yukimochi/followers\"],\"sensitive\":false,\"atomUri\":\"https://mastodon.test.yukimochi.io/users/yukimochi/statuses/101075045564444857\",\"inReplyToAtomUri\":null,\"conversation\":\"tag:mastodon.test.yukimochi.io,2018-11-15:objectId=68:objectType=Conversation\",\"content\":\"<p>Actvity-Relay</p>\",\"contentMap\":{\"en\":\"<p>Actvity-Relay</p>\"},\"attachment\":[],\"tag\":[]},\"signature\":{\"type\":\"RsaSignature2017\",\"creator\":\"https://mastodon.test.yukimochi.io/users/yukimochi#main-key\",\"created\":\"2018-11-15T11:07:26Z\",\"signatureValue\":\"mMgl2GgVPgb1Kw6a2iDIZc7r0j3ob+Cl9y+QkCxIe6KmnUzb15e60UuhkE5j3rJnoTwRKqOFy1PMkSxlYW6fPG/5DBxW9I4kX+8sw8iH/zpwKKUOnXUJEqfwRrNH2ix33xcs/GkKPdedY6iAPV9vGZ10MSMOdypfYgU9r+UI0sTaaC2iMXH0WPnHQuYAI+Q1JDHIbDX5FH1WlDL6+8fKAicf3spBMxDwPHGPK8W2jmDLWdN2Vz4ffsCtWs5BCuqOKZrtTW0Rdd4HWzo40MnRXvBjv7yNlnnKzokANBqiOLWT7kNfK0+Vtnt6c/bNX64KBro53KR7wL3ZBvPVuv5rdQ==\"}}"
		var activity models.Activity
		json.Unmarshal([]byte(body), &activity)
		return activity
	case "Announce":
		file, _ := os.Open("../misc/test/announce.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	case "Undo":
		file, _ := os.Open("../misc/test/undo.json")
		body, _ := io.ReadAll(file)
		var activity models.Activity
		json.Unmarshal(body, &activity)
		return activity
	default:
		panic("fatal - request not registered")
	}
}

func mockActor(req string) models.Actor {
	switch req {
	case "Person":
		file, _ := os.Open("../misc/test/person.json")
		body, _ := io.ReadAll(file)
		var actor models.Actor
		json.Unmarshal(body, &actor)
		return actor
	case "Service":
		file, _ := os.Open("../misc/test/service.json")
		body, _ := io.ReadAll(file)
		var actor models.Actor
		json.Unmarshal(body, &actor)
		return actor
	case "Application":
		file, _ := os.Open("../misc/test/application.json")
		body, _ := io.ReadAll(file)
		var actor models.Actor
		json.Unmarshal(body, &actor)
		return actor
	default:
		panic("fatal - request not registered")
	}
}

func TestSuitableRelayNoBlockService(t *testing.T) {
	activity := mockActivity("Create")
	personActor := mockActor("Person")
	serviceActor := mockActor("Service")
	applicationActor := mockActor("Application")

	RelayState.SetConfig(BlockService, false)

	if isRelayRetransmission(&activity, &personActor) != true {
		t.Fatalf("fail - Person activity should relay")
	}
	if isRelayRetransmission(&activity, &serviceActor) != true {
		t.Fatalf("fail - Service activity should relay")
	}
	if isRelayRetransmission(&activity, &applicationActor) != true {
		t.Fatalf("fail - Service activity should relay")
	}
}

func TestSuitableRelayBlockService(t *testing.T) {
	activity := mockActivity("Create")
	personActor := mockActor("Person")
	serviceActor := mockActor("Service")
	applicationActor := mockActor("Application")

	RelayState.SetConfig(BlockService, true)

	if isRelayRetransmission(&activity, &personActor) != true {
		t.Fatalf("fail - Person activity should relay")
	}
	if isRelayRetransmission(&activity, &serviceActor) != false {
		t.Fatalf("fail - Service activity should not relay when blocking mode")
	}
	if isRelayRetransmission(&activity, &applicationActor) != false {
		t.Fatalf("fail - Application activity should not relay when blocking mode")
	}
	RelayState.SetConfig(BlockService, false)
}

func TestHandleInboxNoSignature(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, decodeActivity)
	}))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleInboxInvalidMethod(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, decodeActivity)
	}))
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 404 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleInboxValidFollow(t *testing.T) {
	activity := mockActivity("Follow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 1 {
		t.Fatalf("fail - follow request not work")
	}
	RelayState.DelSubscription(domain.Host)
}

func TestHandleInboxValidManuallyFollow(t *testing.T) {
	activity := mockActivity("Follow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	// Switch Manually
	RelayState.SetConfig(ManuallyAccept, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:pending:" + domain.Host).Result()
	if res != 1 {
		t.Fatalf("fail - manually follow not work")
	}
	res, _ = RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 0 {
		t.Fatalf("fail - manually follow not work")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.SetConfig(ManuallyAccept, false)
}

func TestHandleInboxInvalidFollow(t *testing.T) {
	activity := mockActivity("Invalid-Follow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.SetConfig(ManuallyAccept, false)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 0 {
		t.Fatalf("fail - follow request not blocked")
	}
}

func TestHandleInboxValidFollowBlocked(t *testing.T) {
	activity := mockActivity("Follow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.SetBlockedDomain(domain.Host, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 0 {
		t.Fatalf("fail - follow request not blocked")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.SetBlockedDomain(domain.Host, false)
}

func TestHandleInboxValidUnfollow(t *testing.T) {
	activity := mockActivity("Unfollow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 0 {
		t.Fatalf("fail - unfollow request not works")
	}
	RelayState.DelSubscription(domain.Host)
}

func TestHandleInboxValidManuallyUnFollow(t *testing.T) {
	activity := mockActivity("Unfollow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.RedisClient.HMSet("relay:pending:"+domain.Host, map[string]interface{}{
		"inbox_url":   actor.Endpoints.SharedInbox,
		"activity_id": activity.ID,
		"type":        "Follow",
		"actor":       actor.ID,
		"object":      mockActivity("Follow"),
	})

	// Switch Manually
	RelayState.SetConfig(ManuallyAccept, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:pending:" + domain.Host).Result()
	if res != 0 {
		t.Fatalf("fail - pending follow request not deleted")
	}
	RelayState.RedisClient.Del("relay:pending:" + domain.Host)
	RelayState.SetConfig(ManuallyAccept, false)
}

func TestHandleInboxInvalidUnfollow(t *testing.T) {
	activity := mockActivity("Invalid-Unfollow")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 1 {
		t.Fatalf("fail - invalid unfollow request should be blocked")
	}
	RelayState.DelSubscription(domain.Host)
}

func TestHandleInboxUnfollowAsActor(t *testing.T) {
	activity := mockActivity("UnfollowAsActor")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 1 {
		t.Fatalf("fail - invalid unfollow request should be blocked")
	}
	RelayState.DelSubscription(domain.Host)
}

func TestHandleInboxValidCreate(t *testing.T) {
	activity := mockActivity("Create")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})
	RelayState.AddSubscription(models.Subscription{
		Domain:   "example.org",
		InboxURL: "https://example.org/inbox",
	})

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.DelSubscription("example.org")
	RelayState.RedisClient.Del("relay:subscription:" + domain.Host).Result()
	RelayState.RedisClient.Del("relay:subscription:example.org").Result()
}

func TestHandleInboxLimitedCreate(t *testing.T) {
	activity := mockActivity("Create")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})
	RelayState.SetLimitedDomain(domain.Host, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.SetLimitedDomain(domain.Host, false)
}

func TestHandleInboxValidCreateAsAnnounceNote(t *testing.T) {
	activity := mockActivity("Create")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})
	RelayState.AddSubscription(models.Subscription{
		Domain:   "example.org",
		InboxURL: "https://example.org/inbox",
	})
	RelayState.SetConfig(CreateAsAnnounce, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.DelSubscription("example.org")
	RelayState.SetConfig(CreateAsAnnounce, false)
}

func TestHandleInboxValidCreateAsAnnounceNoNote(t *testing.T) {
	activity := mockActivity("Create-Article")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})
	RelayState.AddSubscription(models.Subscription{
		Domain:   "example.org",
		InboxURL: "https://example.org/inbox",
	})
	RelayState.SetConfig(CreateAsAnnounce, true)

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	RelayState.DelSubscription(domain.Host)
	RelayState.DelSubscription("example.org")
	RelayState.SetConfig(CreateAsAnnounce, false)
}

func TestHandleInboxUnsubscriptionCreate(t *testing.T) {
	activity := mockActivity("Create")
	actor := mockActor("Person")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 400 {
		t.Fatalf("fail - StatusCode is not match")
	}
}

func TestHandleInboxUndo(t *testing.T) {
	activity := mockActivity("Undo")
	actor := mockActor("Person")
	domain, _ := url.Parse(activity.Actor)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInbox(w, r, mockActivityDecoderProvider(&activity, &actor))
	}))
	defer s.Close()

	RelayState.AddSubscription(models.Subscription{
		Domain:   domain.Host,
		InboxURL: "https://mastodon.test.yukimochi.io/inbox",
	})

	req, _ := http.NewRequest("POST", s.URL, nil)
	client := new(http.Client)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("fail - " + err.Error())
	}
	if r.StatusCode != 202 {
		t.Fatalf("fail - StatusCode is not match")
	}
	res, _ := RelayState.RedisClient.Exists("relay:subscription:" + domain.Host).Result()
	if res != 1 {
		t.Fatalf("fail - undo activity not works")
	}
	RelayState.DelSubscription(domain.Host)
}
