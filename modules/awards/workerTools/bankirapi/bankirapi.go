package bankirapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// API is a magic structure
type API struct {
	token  string
	client *http.Client
}

// RateLimit is a structure for 429 error
type RateLimit struct {
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

// JSONBalanse is a structure for changing user balance
type JSONBalanse struct {
	Cash   int    `json:"cash"`
	Bank   int    `json:"bank"`
	Reason string `json:"reason"`
}

// Balance balance of user
type Balance struct {
	Rank   string `json:"rank"`
	UserID string `json:"user_id"`
	Cash   int    `json:"cash"`
	Bank   int    `json:"bank"`
	Total  int    `json:"total"`
}

// New creates new API object
func New(token string) *API {
	return &API{
		token:  token,
		client: &http.Client{},
	}
}

func (api *API) request(protocol, guildID, userID string, reqBodyBytes io.Reader) (*Balance, error) {
	var (
		err   error
		b     Balance
		limit RateLimit
	)

	req, err := http.NewRequest(protocol, "https://unbelievaboat.com/api/v1/guilds/"+guildID+"/users/"+userID, reqBodyBytes)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", api.token)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		resBodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resBodyBytes, &b)
		if err != nil {
			return nil, err
		}
	}

	if res.StatusCode == 429 {
		resBodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resBodyBytes, &limit)
		if err != nil {
			return nil, err
		}

		time.Sleep(time.Duration(limit.RetryAfter) * time.Millisecond)
		return api.request(protocol, guildID, userID, reqBodyBytes)
	}

	if res.StatusCode != http.StatusOK {
		return &b, errors.New("Strange status code: " + strconv.Itoa(res.StatusCode))
	}

	return &b, nil
}

// GetBalance return balance of user
func (api *API) GetBalance(guildID, userID string) (*Balance, error) {
	return api.request("GET", guildID, userID, nil)
}

// SetBalance sets balance of user
func (api *API) SetBalance(guildID, userID string, cash, bank int, reason string) (*Balance, error) {
	jsonBalanse := JSONBalanse{
		Cash:   cash,
		Bank:   bank,
		Reason: reason,
	}

	reqBodyBytes, err := json.Marshal(jsonBalanse)
	if err != nil {
		return nil, err
	}

	return api.request("PUT", guildID, userID, bytes.NewBuffer(reqBodyBytes))
}

// AddToBalance adds money to users balance
func (api *API) AddToBalance(guildID, userID string, cash, bank int, reason string) (*Balance, error) {
	jsonBal := JSONBalanse{
		Cash:   cash,
		Bank:   bank,
		Reason: reason,
	}

	reqBodyBytes, err := json.Marshal(jsonBal)
	if err != nil {
		return nil, err
	}

	return api.request("PATCH", guildID, userID, bytes.NewBuffer(reqBodyBytes))
}
