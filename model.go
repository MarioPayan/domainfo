package main

import (
	"sort"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func startDB() {
	const addr = "postgresql://domainfouser@localhost:26257/domainfo?sslmode=disable"
	db, _ = gorm.Open("postgres", addr)
	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&Server{})
}

func createDomain(domain *Domain) Domain {
	db.Create(&domain)
	newDomain, _ := getDomain(domain.Url)
	return newDomain
}

func createServer(server *Server) {
	db.FirstOrCreate(&server, server)
}

func getDomain(domainName string) (Domain, bool) {
	var domain Domain
	exist := false
	db.Preload("Servers").Where("url = ?", domainName).First(&domain)
	if domain.Url != "" {
		exist = true
	}
	return domain, exist
}

func getDomains() []Domain {
	var domains []Domain
	db.Preload("Servers").Find(&domains)
	return domains
}

func updateDomain(domain *Domain) Domain {
	domainToUpdate, _ := getDomain(domain.Url)
	db.Model(&domain).Where("url = ?", domain.Url).Update("servers_changed", false)
	servers1 := domain.Servers
	servers2 := domainToUpdate.Servers
	sort.Slice(servers1, func(i, j int) bool { return servers1[i].Address < servers1[j].Address })
	sort.Slice(servers2, func(i, j int) bool { return servers2[i].Address < servers2[j].Address })
	if len(servers1) == len(servers2) {
		for index, _ := range servers1 {
			if servers1[index].Address != servers2[index].Address ||
				servers1[index].SSLGrade != servers2[index].SSLGrade ||
				servers1[index].Country != servers2[index].Country ||
				servers1[index].Owner != servers2[index].Owner ||
				servers1[index].Status != servers2[index].Status {
				domain.ServersChanged = true
			}
		}
	} else {
		domain.ServersChanged = true
	}
	domain.PreviousSSLGrade = domainToUpdate.SSLGrade
	db.Model(&domainToUpdate).Updates(domain)
	for _, server := range domain.Servers {
		createServer(&server)
	}
	newDomain, _ := getDomain(domain.Url)
	return newDomain
}
func deleteDomain(domainName string) {
	var domain Domain
	db.First(&domain, domainName)
	db.Delete(&domain)
}
