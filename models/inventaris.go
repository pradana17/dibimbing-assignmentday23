package models

import "gorm.io/gorm"

type Inventaris struct {
	gorm.Model
	ProdukId uint    `json:"produk_id" gorm:"not null;unique"`
	Jumlah   int     `json:"jumlah" gorm:"not null;default:0"`
	Lokasi   string  `json:"lokasi" gorm:"type:varchar(50);default:NULL"`
	Produk   *Produk `gorm:"foreignKey:ProdukId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type InventarisResponse struct {
	ProdukId uint   `json:"produk_id"`
	Nama     string `json:"nama"`
	Jumlah   int    `json:"jumlah"`
	Lokasi   string `json:"lokasi"`
}

type InventarisUpdateStok struct {
	ProdukId uint `json:"produk_id"`
	Jumlah   int  `json:"jumlah"`
}

type InventarisUpdateLokasi struct {
	ProdukId uint   `json:"produk_id"`
	Lokasi   string `json:"lokasi"`
}

type InventarisByLokasi struct {
	ProdukId uint   `json:"produk_id"`
	Nama     string `json:"nama"`
	Jumlah   int    `json:"jumlah"`
}
