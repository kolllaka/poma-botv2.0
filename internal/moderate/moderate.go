package moderate

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type ApiClient struct {
	baseURL *url.URL

	broadcasterId string
	moderateId    string
	clientId      string
	token         string
}

func NewApiClient(broadcasterId string, moderateId string, clientId string, token string) *ApiClient {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   "helix",
	}

	return &ApiClient{
		baseURL: baseURL,

		broadcasterId: broadcasterId,
		moderateId:    moderateId,
		clientId:      clientId,
		token:         token,
	}
}

func (c *ApiClient) Ban(userId string, duration int, reason string) error {
	resource := "moderation/bans"
	params := url.Values{}
	params.Set("moderator_id", c.moderateId)
	params.Set("broadcaster_id", c.broadcasterId)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	bs, err := json.Marshal(map[string]map[string]any{
		"data": {
			"user_id":  userId,
			"duration": duration,
			"reason":   reason,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%+v\n", string(bs))

	r, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(bs))
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	r.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var post interface{}
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		return err
	}

	log.Println(res.StatusCode)
	log.Printf("%+v\n", post)
	log.Printf("%+v\n", post.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["user_id"])

	return nil
}

func (c *ApiClient) UnBan(userId string) error {
	resource := "moderation/bans"
	params := url.Values{}
	params.Set("moderator_id", c.moderateId)
	params.Set("broadcaster_id", c.broadcasterId)
	params.Set("user_id", userId)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	r, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var post interface{}
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		return err
	}

	log.Println(res.StatusCode)
	log.Printf("%+v\n", post)

	return nil
}

func (c *ApiClient) GetUserId(userId string) (string, error) {
	resource := "users"
	params := url.Values{}
	params.Set("login", userId)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	type Resp struct {
		Data []struct {
			UserId string `json:"id"`
		} `json:"data"`
	}

	r, err := http.NewRequest(http.MethodGet, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resp := Resp{}
	json.NewDecoder(res.Body).Decode(&resp)

	log.Println(res.StatusCode)
	log.Printf("%+v\n", resp)

	return resp.Data[0].UserId, nil
}

func (c *ApiClient) GetVips(broadcaster_id string, usersId []string) (map[string]string, error) {
	resource := "channels/vips"
	params := url.Values{}
	params.Set("broadcaster_id", broadcaster_id)
	for _, userId := range usersId {
		params.Add("user_id", userId)
	}

	vips := make(map[string]string)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	type Resp struct {
		Data []struct {
			UserId    string `json:"user_id"`
			UserLogin string `json:"user_login"`
		} `json:"data"`
	}
	log.Printf("[vips] %s\n", url.String())

	r, err := http.NewRequest(http.MethodGet, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return vips, err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return vips, err
	}
	defer res.Body.Close()

	resp := Resp{}
	json.NewDecoder(res.Body).Decode(&resp)

	log.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return vips, errors.New("no vips found")
	}
	log.Printf("[vips] %+v\n", resp)
	for _, vip := range resp.Data {
		vips[vip.UserLogin] = vip.UserId
	}

	return vips, nil
}

func (c *ApiClient) GetModerators(broadcaster_id string, usersId []string) (map[string]string, error) {
	resource := "moderation/moderators"
	params := url.Values{}
	params.Set("broadcaster_id", broadcaster_id)
	for _, userId := range usersId {
		params.Add("user_id", userId)
	}

	moderators := make(map[string]string)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	type Resp struct {
		Data []struct {
			UserId    string `json:"user_id"`
			UserLogin string `json:"user_login"`
		} `json:"data"`
	}
	log.Printf("[moderators] %s\n", url.String())

	r, err := http.NewRequest(http.MethodGet, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return moderators, err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return moderators, err
	}
	defer res.Body.Close()

	resp := Resp{}
	json.NewDecoder(res.Body).Decode(&resp)

	log.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return moderators, errors.New("no moderators found")
	}
	log.Printf("[moderators] %+v\n", resp)
	for _, vip := range resp.Data {
		moderators[vip.UserLogin] = vip.UserId
	}

	return moderators, nil
}

func (c *ApiClient) GetSubscriptions(broadcaster_id string, usersId []string) (map[string]string, error) {
	resource := "subscriptions"
	params := url.Values{}
	params.Set("broadcaster_id", broadcaster_id)
	for _, userId := range usersId {
		params.Add("user_id", userId)
	}

	subscriptions := make(map[string]string)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	type Resp struct {
		Data []struct {
			UserId    string `json:"user_id"`
			UserLogin string `json:"user_login"`
		} `json:"data"`
	}
	log.Printf("[subscriptions] %s\n", url.String())

	r, err := http.NewRequest(http.MethodGet, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return subscriptions, err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return subscriptions, err
	}
	defer res.Body.Close()

	resp := Resp{}
	json.NewDecoder(res.Body).Decode(&resp)

	log.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return subscriptions, errors.New("no moderators found")
	}
	log.Printf("[subscriptions] %+v\n", resp)
	for _, vip := range resp.Data {
		subscriptions[vip.UserLogin] = vip.UserId
	}

	return subscriptions, nil
}

func (c *ApiClient) AddModerator(broadcaster_id string, userId string) error {
	resource := "moderation/moderators"
	params := url.Values{}
	params.Set("broadcaster_id", broadcaster_id)
	params.Set("user_id", userId)

	url := c.baseURL.JoinPath(resource)
	url.RawQuery = params.Encode()

	r, err := http.NewRequest(http.MethodPost, url.String(), nil)
	r.Header.Add("Authorization", c.token)
	r.Header.Add("Client-Id", c.clientId)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var resp interface{}
	json.NewDecoder(res.Body).Decode(&resp)

	log.Println(res.StatusCode)
	log.Printf("[moderators] %+v\n", resp)

	return nil
}
