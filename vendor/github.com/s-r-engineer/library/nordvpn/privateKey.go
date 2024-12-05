package libraryNordvpn

import (
	"encoding/json"
	"io"
	"net/http"
)

const DefaultUserCredentialsURL = "https://api.nordvpn.com/v1/users/services/credentials"

func FetchOwnPrivateKey(token string) (string, error) {
	url := DefaultUserCredentialsURL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	servers := Creds{}
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return "", err
	}
	return servers.NordlynxPrivateKey, nil
}
