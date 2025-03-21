package libraryNordvpn

import (
	"encoding/json"
	"fmt"
	libraryHttp "github.com/s-r-engineer/library/http"
	"io"
	"net"
)

const DefaultRecommendationsURL = "https://api.nordvpn.com/v1/servers/recommendations?filters[servers_technologies][identifier]=wireguard_udp&limit=1"

// FetchServerData will return hostname, ip, public key, country code and error
func FetchServerData(country int) (string, string, string, string, error) {
	url := DefaultRecommendationsURL
	if country > 0 {
		url += fmt.Sprintf("&filters[country_id]=%d", country)
	}
	resp, err := libraryHttp.GetUrl(url)
	if err != nil {
		return "", "", "", "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", err
	}
	servers := Servers{}
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return "", "", "", "", err
	}
	hostname := servers[0].Hostname
	var publicKey string
	for _, technology := range servers[0].Technologies {
		if technology.Identifier != "wireguard_udp" {
			continue
		}
		publicKey = technology.Metadata[0].Value
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", "", "", "", err
	}
	return hostname, ips[0].String(), publicKey, servers[0].Locations[0].Country.Code, nil
}

type Creds struct {
	ID                 int    `json:"id,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	Username           string `json:"username,omitempty"`
	Password           string `json:"password,omitempty"`
	NordlynxPrivateKey string `json:"nordlynx_private_key,omitempty"`
}

type Servers []Server

type Server struct {
	ID           int            `json:"id,omitempty"`
	CreatedAt    string         `json:"created_at,omitempty"`
	UpdatedAt    string         `json:"updated_at,omitempty"`
	Name         string         `json:"name,omitempty"`
	Station      string         `json:"station,omitempty"`
	Ipv6Station  string         `json:"ipv6_station,omitempty"`
	Hostname     string         `json:"hostname,omitempty"`
	Load         int            `json:"load,omitempty"`
	Status       string         `json:"status,omitempty"`
	Locations    []Locations    `json:"locations,omitempty"`
	Technologies []Technologies `json:"technologies,omitempty"`
}
type City struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	DNSName   string  `json:"dns_name,omitempty"`
	HubScore  int     `json:"hub_score,omitempty"`
}
type Country struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
	City City   `json:"city,omitempty"`
}
type Locations struct {
	ID        int     `json:"id,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Country   Country `json:"country,omitempty"`
}
type Services struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}
type Pivot struct {
	TechnologyID int    `json:"technology_id,omitempty"`
	ServerID     int    `json:"server_id,omitempty"`
	Status       string `json:"status,omitempty"`
}
type Technologies struct {
	ID         int        `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Identifier string     `json:"identifier,omitempty"`
	CreatedAt  string     `json:"created_at,omitempty"`
	UpdatedAt  string     `json:"updated_at,omitempty"`
	Metadata   []Metadata `json:"metadata,omitempty"`
	Pivot      Pivot      `json:"pivot,omitempty"`
}

type Metadata struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type Type struct {
	ID         int    `json:"id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	Title      string `json:"title,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}
type Groups struct {
	ID         int    `json:"id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	Title      string `json:"title,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Type       Type   `json:"type,omitempty"`
}
type Values struct {
	ID    int    `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}
type Specifications struct {
	ID         int      `json:"id,omitempty"`
	Title      string   `json:"title,omitempty"`
	Identifier string   `json:"identifier,omitempty"`
	Values     []Values `json:"values,omitempty"`
}
type IP struct {
	ID      int    `json:"id,omitempty"`
	IP      string `json:"ip,omitempty"`
	Version int    `json:"version,omitempty"`
}
type Ips struct {
	ID        int    `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	ServerID  int    `json:"server_id,omitempty"`
	IPID      int    `json:"ip_id,omitempty"`
	Type      string `json:"type,omitempty"`
	IP        IP     `json:"ip,omitempty"`
}
