package database

import (
	"time"
	"gorm.io/gorm"
)

type Shell struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	URL          string    `gorm:"not null" json:"url"`
	Password     string    `gorm:"not null" json:"-"`
	EncodeType   string    `gorm:"default:''" json:"encode_type"`
	Protocol     string    `gorm:"default:multipart" json:"protocol"`
	Method       string    `gorm:"default:POST" json:"method"`
	Group        string    `json:"group"`
	Name         string    `json:"name"`
	Status       string    `gorm:"default:offline" json:"status"`
	LastActive   time.Time `json:"last_active"`
	Latency      int       `json:"latency"`
	SystemInfo   string    `gorm:"type:text" json:"system_info"`
	SystemInfoJSON string  `gorm:"type:text" json:"system_info_json"`
	CustomHeaders string   `gorm:"type:text" json:"custom_headers"`
	ProxyID      uint      `json:"proxy_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShellID   uint      `gorm:"not null" json:"shell_id"`
	Type      string    `gorm:"not null" json:"type"`
	Command   string    `gorm:"type:text" json:"command"`
	Result    string    `gorm:"type:text" json:"result"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Proxy struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	Host      string    `gorm:"not null" json:"host"`
	Port      int       `gorm:"not null" json:"port"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

