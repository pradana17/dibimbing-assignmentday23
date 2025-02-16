package main

import (
	"assignmentday23/config"
	"assignmentday23/controllers"
	"assignmentday23/middleware"
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
	db.AutoMigrate(&models.Produk{}, &models.Inventaris{}, &models.Pesanan{}, &models.DetailPesanan{}, &models.User{})
	authController := controllers.NewAuthController(db)
	produkController := controllers.NewProdukController(db)
	inventarisController := controllers.NewInventarisController(db)
	pesananController := controllers.NewPesananController(db)
	detailPesananController := controllers.NewDetailPesananController(db)
	userController := controllers.NewUserController(db)
	fileController := controllers.NewFileController(db)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			r.POST("/user", userController.CreateUser)
			r.GET("/user", userController.GetUser)
			r.POST("/file", fileController.CreateDirectory)
			r.POST("/file/:directory", fileController.CreateFile)
			r.GET("/file/:directory/:file", fileController.ReadFile)
			r.PUT("/file/:directory/:file", fileController.RenameFile)
			r.POST("/file/:directory/:file", fileController.UploadFile)
			r.GET("/file/:directory/:file", fileController.DownloadFile)
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
		}
	}

	r.Run(":8080")
}
