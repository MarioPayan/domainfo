package main

import (
	"time"
)

type DomainModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Domain struct {
	DomainModel
	Url              string `gorm:"unique;not null"`
	ServersChanged   bool
	SSLGrade         string
	PreviousSSLGrade string
	Logo             string
	Title            string
	IsDown           bool
	Status           string
	Servers          []Server `gorm:"foreignkey:ID"`
}

type Server struct {
	Address  string `gorm:"primary_key"`
	ID       uint
	SSLGrade string
	Country  string
	Owner    string
	Status   string
}

type SSLData struct {
	Host            string     `json:"host,omitempty"`
	Port            int        `json:"port,omitempty"`
	Protocol        string     `json:"protocol,omitempty"`
	IsPublic        bool       `json:"isPublic,omitempty"`
	Status          string     `json:"status,omitempty"`
	StartTime       int64      `json:"startTime,omitempty"`
	TestTime        int64      `json:"testTime,omitempty"`
	EngineVersion   string     `json:"engineVersion,omitempty"`
	CriteriaVersion string     `json:"criteriaVersion,omitempty"`
	Endpoints       []Endpoint `json:"endpoints,omitempty"`
}

type Endpoint struct {
	IPAddress         string `json:"ipAddress,omitempty"`
	ServerName        string `json:"serverName,omitempty"`
	StatusMessage     string `json:"statusMessage,omitempty"`
	Grade             string `json:"grade,omitempty"`
	GradeTrustIgnored string `json:"gradeTrustIgnored,omitempty"`
	HasWarnings       bool   `json:"hasWarnings,omitempty"`
	IsExceptional     bool   `json:"isExceptional,omitempty"`
	Progress          int    `json:"progress,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	Delegation        int    `json:"delegation,omitempty"`
}
