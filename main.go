package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Host struct {
	HostId             string `xml:"HostId,attr"`
	Name               string `xml:"Name,attr"`
	Type               string `xml:"Type,attr"`
	Address            string `xml:"Address,attr"`
	MXPref             string `xml:"MXPref,attr"`
	TTL                string `xml:"TTL,attr"`
	AssociatedAppTitle string `xml:"AssociatedAppTitle,attr"`
	FriendlyName       string `xml:"FriendlyName,attr"`
	IsActive           string `xml:"IsActive,attr"`
	IsDDNSEnabled      string `xml:"IsDDNSEnabled,attr"`
}

type DomainDNSGetHostsResult struct {
	Domain string `xml:"Domain,attr"`
	Hosts  []Host `xml:"host"`
}

type CommandResponse struct {
	Type                    string                  `xml:"Type,attr"`
	DomainDNSGetHostsResult DomainDNSGetHostsResult `xml:"DomainDNSGetHostsResult"`
}

type ApiResponse struct {
	XMLName         xml.Name        `xml:"ApiResponse"`
	Status          string          `xml:"Status,attr"`
	CommandResponse CommandResponse `xml:"CommandResponse"`
}

var SetHost string
var GetHost string
var AuthURL string

func GetAuthURL(apiUser, apiKey, clientIp, sld, tld string) string {
	const baseURL = "https://api.namecheap.com/xml.response?apiUser=%s&apiKey=%s&UserName=%s&ClientIp=%s&SLD=%s&TLD=%s&Command=namecheap.domains.dns."
	return fmt.Sprintf(baseURL, apiUser, apiKey, apiUser, clientIp, sld, tld)
}

func main() {
	apiUser := flag.String("u", "", "API User")
	apiKey := flag.String("k", "", "API Key")
	clientIp := flag.String("c", "", "Client IP")
	sld := flag.String("s", "", "SLD")
	tld := flag.String("t", "", "TLD")
	host := flag.String("h", "", "Host")
	value := flag.String("v", "", "Value")

	// Parse flags
	flag.Parse()
	AuthURL = GetAuthURL(*apiUser, *apiKey, *clientIp, *sld, *tld)
	SetHost = AuthURL + "setHosts"
	GetHost = AuthURL + "getHosts"
	if *apiUser == "" || *apiKey == "" || *clientIp == "" || *sld == "" || *tld == "" || *value == "" {
		fmt.Println("Please provide all necessary flags.")
		flag.Usage()
		return
	}

	// 读取环境
	dd, err := HttpGet(GetHost)
	var response ApiResponse
	SetName := "_acme-challenge"
	if host != nil && *host != "" {
		SetName += "." + *host
	}

	err = xml.Unmarshal(dd, &response)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	for i, h := range response.CommandResponse.DomainDNSGetHostsResult.Hosts {
		if h.Name == SetName {
			fmt.Println("find: " + SetName + " Host")
			fmt.Println("Set: " + SetName + " Value:" + *value)
			response.CommandResponse.DomainDNSGetHostsResult.Hosts[i].Address = *value
		}
	}
	time.Sleep(5 * time.Second)
	_, err = HttpGet(generateURL(response.CommandResponse.DomainDNSGetHostsResult.Hosts))
	if err != nil {
		fmt.Println("Send SetHosts URL Error:", err)
		return
	}
	//fmt.Println(string(res))
	time.Sleep(5 * time.Second)
	ddc, err := HttpGet(GetHost)
	var response2 ApiResponse
	err = xml.Unmarshal(ddc, &response2)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	success := false
	for _, host := range response.CommandResponse.DomainDNSGetHostsResult.Hosts {
		if host.Name == SetName && host.Address == *value {
			fmt.Println("Set Success!!!!")
			success = true
		}
	}
	if !success {
		fmt.Println("Set Fail!!!!")
	}

}
func generateURL(hosts []Host) string {

	params := ""
	for i, host := range hosts {
		params += fmt.Sprintf("&HostName%d=%s&RecordType%d=%s&Address%d=%s&TTL%d=%s", i+1, host.Name, i+1, host.Type, i+1, host.Address, i+1, host.TTL)
	}
	return SetHost + params
}

func HttpGet(domain string) ([]byte, error) {
	// 使用http.Get函数发起GET请求
	resp, err := http.Get(domain)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return body, nil
}
