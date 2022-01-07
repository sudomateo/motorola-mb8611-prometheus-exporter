package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ModemMetrics struct {
	GetMultipleHNAPsResponse GetMultipleHNAPsResponse `json:"GetMultipleHNAPsResponse"`
}
type StartupSequenceResponse struct {
	DownstreamFreq           string `json:"MotoConnDSFreq"`
	DownstreamComment        string `json:"MotoConnDSComment"`
	ConnectivityStatus       string `json:"MotoConnConnectivityStatus"`
	ConnectivityComment      string `json:"MotoConnConnectivityComment"`
	BootStatus               string `json:"MotoConnBootStatus"`
	BootComment              string `json:"MotoConnBootComment"`
	ConfigurationFileStatus  string `json:"MotoConnConfigurationFileStatus"`
	ConfigurationFileComment string `json:"MotoConnConfigurationFileComment"`
	SecurityStatus           string `json:"MotoConnSecurityStatus"`
	SecurityComment          string `json:"MotoConnSecurityComment"`
	StartupSequenceResult    string `json:"GetMotoStatusStartupSequenceResult"`
}
type ConnectionInfoResponse struct {
	SystemUpTime         string `json:"MotoConnSystemUpTime"`
	NetworkAccess        string `json:"MotoConnNetworkAccess"`
	ConnectionInfoResult string `json:"GetMotoStatusConnectionInfoResult"`
}
type DownstreamChannelInfoResponse struct {
	DownstreamChannel           string `json:"MotoConnDownstreamChannel"`
	DownstreamChannelInfoResult string `json:"GetMotoStatusDownstreamChannelInfoResult"`
}
type UpstreamChannelInfoResponse struct {
	UpstreamChannel           string `json:"MotoConnUpstreamChannel"`
	UpstreamChannelInfoResult string `json:"GetMotoStatusUpstreamChannelInfoResult"`
}
type LagStatusResponse struct {
	CurrentStatus   string `json:"MotoLagCurrentStatus"`
	LagStatusResult string `json:"GetMotoLagStatusResult"`
}
type GetMultipleHNAPsResponse struct {
	StartupSequenceResponse       StartupSequenceResponse       `json:"GetMotoStatusStartupSequenceResponse"`
	ConnectionInfoResponse        ConnectionInfoResponse        `json:"GetMotoStatusConnectionInfoResponse"`
	DownstreamChannelInfoResponse DownstreamChannelInfoResponse `json:"GetMotoStatusDownstreamChannelInfoResponse"`
	UpstreamChannelInfoResponse   UpstreamChannelInfoResponse   `json:"GetMotoStatusUpstreamChannelInfoResponse"`
	LagStatusResponse             LagStatusResponse             `json:"GetMotoLagStatusResponse"`
	GetMultipleHNAPsResult        string                        `json:"GetMultipleHNAPsResult"`
}

func main() {
	ctx := context.Background()

	ip := os.Getenv("MODEM_IP")
	if ip == "" {
		log.Fatalln("Missing MODEM_IP environment variable!")
	}

	metrics, err := fetch(ctx, ip)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v\n", metrics.GetMultipleHNAPsResponse.DownstreamChannelInfoResponse.DownstreamChannel)
}

func fetch(ctx context.Context, ip string) (*ModemMetrics, error) {
	url := fmt.Sprintf("https://%s", ip)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	d := `{"GetMultipleHNAPs":{"GetMotoStatusStartupSequence":"","GetMotoStatusConnectionInfo":"","GetMotoStatusDownstreamChannelInfo":"","GetMotoStatusUpstreamChannelInfo":"","GetMotoLagStatus":""}}`

	body := strings.NewReader(d)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SOAPACTION", "http://purenetworks.com/HNAP1/GetMultipleHNAPs")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metrics ModemMetrics

	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}
