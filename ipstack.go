package locip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	IP_STACK_URL = "https://ipstack.com/ipstack_api.php"
)

func IPLocationFromIPStack(ip string) (*IPStackResponse, error) {
	// 1.2.3.4:8080 -> 1.2.3.4
	if ind := strings.Index(ip, ":"); ind > 0 {
		ip = ip[:ind]
	}

	form := url.Values{}
	form.Add("ip", ip)
	sendTo := fmt.Sprintf("%s?%s", IP_STACK_URL, form.Encode())
	resp, err := http.Get(sendTo)
	if err != nil {
		return nil, fmt.Errorf("Failed send to %s: %+v", sendTo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("No 200 (%d) when send to %s: %+v", resp.StatusCode, sendTo, err)
	}

	var isr IPStackResponse
	err = json.NewDecoder(resp.Body).Decode(&isr)
	if err != nil {
		return nil, fmt.Errorf("Failed decode body of %s: %+v", sendTo, err)
	}

	return &isr, isr.Validate()
}

type Location struct {
	GeonameID int    `json:"geoname_id"`
	Capital   string `json:"capital"`
	Languages []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Native string `json:"native"`
	} `json:"languages"`
	CountryFlag             string `json:"country_flag"`
	CountryFlagEmoji        string `json:"country_flag_emoji"`
	CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
	CallingCode             string `json:"calling_code"`
	IsEu                    bool   `json:"is_eu"`
}

type TimeZone struct {
	ID               string `json:"id"`
	CurrentTime      string `json:"current_time"`
	GmtOffset        int    `json:"gmt_offset"`
	Code             string `json:"code"`
	IsDaylightSaving bool   `json:"is_daylight_saving"`
}

type Connection struct {
	Asn int    `json:"asn"`
	Isp string `json:"isp"`
}

type IPStackResponse struct {
	IP            string  `json:"ip"`
	Hostname      string  `json:"hostname"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`

	Location   Location   `json:"location"`
	TimeZone   TimeZone   `json:"time_zone"`
	Connection Connection `json:"connection"`

	/*
		Security   struct {
			IsProxy     bool        `json:"is_proxy"`
			ProxyType   interface{} `json:"proxy_type"`
			IsCrawler   bool        `json:"is_crawler"`
			CrawlerName interface{} `json:"crawler_name"`
			CrawlerType interface{} `json:"crawler_type"`
			IsTor       bool        `json:"is_tor"`
			ThreatLevel string      `json:"threat_level"`
			ThreatTypes interface{} `json:"threat_types"`
		} `json:"security"`

		Currency struct {
			Code         string `json:"code"`
			Name         string `json:"name"`
			Plural       string `json:"plural"`
			Symbol       string `json:"symbol"`
			SymbolNative string `json:"symbol_native"`
		} `json:"currency"`
	*/
}

func (isp *IPStackResponse) Validate() error {
	if isp.IP == "" {
		return fmt.Errorf("IP empty")
	}

	if isp.Hostname == "" {
		return fmt.Errorf("Hostname empty")
	}

	if isp.CountryCode == "" && isp.CountryName == "" {
		return fmt.Errorf("CountryCode/CountryName empty")
	}

	if isp.Latitude == 0 && isp.Longitude == 0 {
		return fmt.Errorf("Latitude/Longitude empty")
	}

	return nil
}
