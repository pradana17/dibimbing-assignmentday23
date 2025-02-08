package main

import (
	"assignmentday23/config"
	"assignmentday23/controllers"
	"assignmentday23/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	r := gin.Default()
	db := config.ConnectDB()
	db.AutoMigrate(&models.Produk{}, &models.Inventaris{}, &models.Pesanan{}, &models.DetailPesanan{})
	produkController := controllers.NewProdukController(db)
	inventarisController := controllers.NewInventarisController(db)
	pesananController := controllers.NewPesananController(db)
	detailPesananController := controllers.NewDetailPesananController(db)
	r.POST("/produk", produkController.CreateProduk)
	r.GET("/produk", produkController.GetAllProduk)
	r.GET("/produk/id/:id", produkController.GetProdukById)
	r.GET("/produk/kategori/:kategori", produkController.GetProdukByKategori)
	r.DELETE("/produk/:id", produkController.DeleteProduk)
	r.PUT("/produk/:id", produkController.UpdateProduk)
	r.GET("/inventaris", inventarisController.GetAllInventaris)
	r.GET("/inventaris/id/:produk_id", inventarisController.GetInventarisByProdukId)
	r.PUT("/inventaris/stok/:produk_id", inventarisController.InventarisUpdateStok)
	r.PUT("/inventaris/lokasi/:produk_id", inventarisController.InventarisUpdateLokasi)
	r.GET("/inventaris/lokasi/:lokasi", inventarisController.CekStokByLokasi)
	r.POST("/pesanan", pesananController.CreatePesanan)
	r.GET("/pesanan/:pesanan_id", pesananController.GetDetailPesananById)
	r.POST("/pesanan/:pesanan_id", detailPesananController.CreateDetailPesananById)
	r.Run(":8080")
}
