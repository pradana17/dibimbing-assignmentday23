package models

import "gorm.io/gorm"

type Produk struct {
	gorm.Model
	Nama       string     `json:"nama" gorm:"unique;not null;type:varchar(50)"`
	Harga      int        `json:"harga" gorm:"not null"`
	Deskripsi  string     `json:"deskripsi" gorm:"type:text"`
	Kategori   string     `json:"kategori" gorm:"type:varchar(25)"`
	Inventaris Inventaris `gorm:"foreignKey:ProdukId"`
}

type ProdukResponse struct {
	ID        uint   `json:"id"`
	Nama      string `json:"nama"`
	Harga     int    `json:"harga"`
	Deskripsi string `json:"deskripsi"`
	Kategori  string `json:"kategori"`
	Jumlah    int    `json:"jumlah"`
	Lokasi    string `json:"lokasi"`
}

type ProdukRequest struct {
	Nama      string `json:"nama"`
	Harga     int    `json:"harga"`
	Deskripsi string `json:"deskripsi"`
	Kategori  string `json:"kategori"`
}
