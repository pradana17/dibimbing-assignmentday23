package controllers

import (
	"assignmentday23/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventarisController struct {
	DB *gorm.DB
}

func NewInventarisController(db *gorm.DB) *InventarisController {
	return &InventarisController{
		DB: db,
	}
}

func (ic *InventarisController) GetInventarisByProdukId(c *gin.Context) {
	var inventaris models.Inventaris
	produkId := c.Param("produk_id")

	if err := ic.DB.Preload("Produk").First(&inventaris, "produk_id = ?", produkId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.InventarisResponse{
		ProdukId: inventaris.ProdukId,
		Nama:     inventaris.Produk.Nama,
		Jumlah:   inventaris.Jumlah,
		Lokasi:   inventaris.Lokasi,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (ic *InventarisController) InventarisUpdateStok(c *gin.Context) {
	var req models.InventarisUpdateStok
	produkId := c.Param("produk_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := ic.DB.Begin()
	defer tx.Commit()
	var inventaris models.Inventaris
	if err := tx.Preload("Produk").First(&inventaris, "produk_id = ?", produkId).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	inventaris.Jumlah = req.Jumlah
	if err := tx.Save(&inventaris).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stok produk berhasil di update"})
}

func (ic *InventarisController) InventarisUpdateLokasi(c *gin.Context) {
	var req models.InventarisUpdateLokasi
	produkId := c.Param("produk_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ic.DB.Begin()
	defer tx.Commit()
	var inventaris models.Inventaris
	if err := tx.Preload("Produk").First(&inventaris, "produk_id = ?", produkId).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	inventaris.Lokasi = req.Lokasi

	if err := tx.Save(&inventaris).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lokasi produk berhasil di update"})
}

func (ic *InventarisController) CekStokByLokasi(c *gin.Context) {
	var inventaris []models.Inventaris
	lokasi := c.Param("lokasi")
	if err := ic.DB.Preload("Produk").Where("lokasi = ?", lokasi).Find(&inventaris).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]models.InventarisByLokasi, len(inventaris))
	for i, invent := range inventaris {
		response[i] = models.InventarisByLokasi{
			ProdukId: invent.ProdukId,
			Nama:     invent.Produk.Nama,
			Jumlah:   invent.Jumlah,
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (ic *InventarisController) GetAllInventaris(c *gin.Context) {
	var inventaris []models.Inventaris
	if err := ic.DB.Preload("Produk").Find(&inventaris).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]models.InventarisResponse, len(inventaris))
	for i, invent := range inventaris {
		response[i] = models.InventarisResponse{
			ProdukId: invent.ProdukId,
			Nama:     invent.Produk.Nama,
			Jumlah:   invent.Jumlah,
			Lokasi:   invent.Lokasi,
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}
