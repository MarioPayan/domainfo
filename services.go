package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func getData(domainName string) *Domain {
	domain := new(Domain)
	sslData := getSSLlabsData(domainName)
	domain.Url = domainName
	domain.Status = sslData.Status
	if domain.Status != READY && domain.Status != IN_PROGRESS {
		return domain
	}
	html := getHtml("https://www." + domainName + "/")
	domain.Title = getTitle(html)
	domain.Logo = getLogo(html)
	if sslData.Status == READY || sslData.Status == IN_PROGRESS {
		domain.IsDown = false
	} else {
		domain.IsDown = true
	}
	var SSLGrades []string
	for _, endpoint := range sslData.Endpoints {
		server := new(Server)
		whois := string(getWhoIsString(endpoint.IPAddress))
		server.Address = endpoint.IPAddress
		server.SSLGrade = endpoint.Grade
		server.Country = getCountry(whois)
		server.Owner = getOwner(whois)
		server.Status = endpoint.StatusMessage
		domain.Servers = append(domain.Servers, *server)
		SSLGrades = append(SSLGrades, server.SSLGrade)
	}
	domain.SSLGrade = getMaxSSLGrade(SSLGrades)
	return domain
}

func getHtml(domain string) string {
	resp, err1 := http.Get(domain)
	manageError(err1, "Can't get HTML from "+domain)
	defer resp.Body.Close()
	htmlDataBytes, err2 := ioutil.ReadAll(resp.Body)
	manageError(err2, "Can't get HTML from "+domain)
	htmlData := strings.ReplaceAll(string(htmlDataBytes), "\n", "")
	htmlData = strings.ReplaceAll(htmlData, " ", "")
	return string(htmlData)
}

func getSSLlabsData(domain string) *SSLData {
	sslData := new(SSLData)
	client := &http.Client{Timeout: 15 * time.Second}
	baseUrl := "https://api.ssllabs.com/api/v3/analyze"
	url := baseUrl + "?host=" + domain
	response, err := client.Get(url)
	manageError(err, "Can't get data from ssllabs")
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(sslData)
	return sslData
}

func getWhoIsString(arg string) []byte {
	cmd := exec.Command("whois", arg)
	out, err := cmd.CombinedOutput()
	manageError(err, "cmd.Run() failed")
	return out
}

func getCountry(whoisString string) string {
	return getRegexMatch([]string{"Country:\\s+(\\w+)"}, string(whoisString), "country", "Whois")
}

func getOwner(whoisString string) string {
	return getRegexMatch([]string{"OrgName:\\s+(\\w+.*)", "organisation:\\s+(\\w+.*)"}, string(whoisString), "owner", "Whois")
}

func getTitle(htmlData string) string {
	return getRegexMatch([]string{"<title>(\\w*)</title>"}, htmlData, "title", "HTML")
}

func getLogo(htmlData string) string {
	return getRegexMatch([]string{"<head>.*\"(http.*\\.png)\".*</head>"}, htmlData, "logo", "HTML")
}

func getMaxSSLGrade(SSLGrades []string) string {
	sort.Strings(SSLGrades)
	SSLFirstGrade := SSLGrades[0]
	for _, SSLGrade := range SSLGrades {
		if SSLGrade == SSLFirstGrade+"+" {
			return SSLGrade
		} else if SSLGrade != SSLFirstGrade {
			return SSLFirstGrade
		}
	}
	return SSLGrades[0]
}
