package models

import (
	"time"

	"gorm.io/gorm"
)

type Pesanan struct {
	gorm.Model
	Tanggal       time.Time       `json:"tanggal" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	DetailPesanan []DetailPesanan `gorm:"foreignKey:PesananId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DetailPesanan struct {
	gorm.Model
	PesananId uint    `json:"pesanan_id" gorm:"not null"`
	ProdukId  uint    `json:"produk_id" gorm:"not null"`
	Jumlah    int     `json:"jumlah" gorm:"not null"`
	Produk    *Produk `gorm:"foreignKey:ProdukId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TambahPesananRequest struct {
	ProdukName string `json:"produk_name"`
	Jumlah     int    `json:"jumlah"`
}

type DetailPesananResponse struct {
	Nama   string `json:"nama"`
	Jumlah int    `json:"jumlah"`
	Harga  int    `json:"harga"`
	Total  int    `json:"total"`
}

type PesananResponse struct {
	PesananId     uint                    `json:"pesanan_id"`
	DetailPesanan []DetailPesananResponse `json:"detail_pesanan"`
	Total         int                     `json:"total"`
}
