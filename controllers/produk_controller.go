package controllers

import (
	"assignmentday23/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProdukController struct {
	DB *gorm.DB
}

func NewProdukController(db *gorm.DB) *ProdukController {
	return &ProdukController{
		DB: db,
	}
}

func (pc *ProdukController) CreateProduk(c *gin.Context) {
	var produk models.Produk
	if err := c.ShouldBindJSON(&produk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := pc.DB.Begin()
	defer tx.Commit()
	if err := tx.Create(&produk).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	inventaris := models.Inventaris{
		ProdukId: produk.ID,
	}

	if err := tx.Create(&inventaris).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": produk})
}

func (pc *ProdukController) GetAllProduk(c *gin.Context) {
	var produk []models.Produk
	if err := pc.DB.Preload("Inventaris").Find(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]models.ProdukResponse, len(produk))
	for i, p := range produk {
		response[i] = models.ProdukResponse{
			ID:        p.ID,
			Nama:      p.Nama,
			Harga:     p.Harga,
			Deskripsi: p.Deskripsi,
			Kategori:  p.Kategori,
			Jumlah:    p.Inventaris.Jumlah,
			Lokasi:    p.Inventaris.Lokasi,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (pc *ProdukController) GetProdukById(c *gin.Context) {
	var produk models.Produk
	produkId := c.Param("id")

	if err := pc.DB.Preload("Inventaris").First(&produk, "id = ?", produkId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ProdukResponse{
		ID:        produk.ID,
		Nama:      produk.Nama,
		Harga:     produk.Harga,
		Deskripsi: produk.Deskripsi,
		Kategori:  produk.Kategori,
		Jumlah:    produk.Inventaris.Jumlah,
		Lokasi:    produk.Inventaris.Lokasi,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})

}

func (pc *ProdukController) UpdateProduk(c *gin.Context) {
	var produk models.Produk
	produkId := c.Param("id")

	if err := pc.DB.Preload("Inventaris").First(&produk, "id = ?", produkId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&produk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var produkrequest models.ProdukRequest
	tx := pc.DB.Begin()
	defer tx.Commit()
	if err := tx.Model(&produk).Updates(&produkrequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Save(&produk).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ProdukResponse{
		ID:        produk.ID,
		Nama:      produk.Nama,
		Harga:     produk.Harga,
		Deskripsi: produk.Deskripsi,
		Kategori:  produk.Kategori,
		Jumlah:    produk.Inventaris.Jumlah,
		Lokasi:    produk.Inventaris.Lokasi,
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (pc *ProdukController) DeleteProduk(c *gin.Context) {
	var produk models.Produk
	produkId := c.Param("id")

	if err := pc.DB.Preload("Inventaris").First(&produk, "id = ?", produkId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := pc.DB.Begin()
	defer tx.Commit()
	if err := tx.Delete(&produk).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Delete(&produk.Inventaris).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk deleted successfully"})
}

func (pc *ProdukController) GetProdukByKategori(c *gin.Context) {
	kategori := c.Param("kategori")
	var produk []models.Produk
	if err := pc.DB.Preload("Inventaris").Where("kategori = ?", kategori).Find(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]models.ProdukResponse, len(produk))
	for i, p := range produk {
		response[i] = models.ProdukResponse{
			ID:        p.ID,
			Nama:      p.Nama,
			Harga:     p.Harga,
			Deskripsi: p.Deskripsi,
			Kategori:  p.Kategori,
			Jumlah:    p.Inventaris.Jumlah,
			Lokasi:    p.Inventaris.Lokasi,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
