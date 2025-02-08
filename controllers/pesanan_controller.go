package controllers

import (
	"assignmentday23/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PesananController struct {
	DB *gorm.DB
}

type DetailPesananController struct {
	DB *gorm.DB
}

func NewPesananController(db *gorm.DB) *PesananController {
	return &PesananController{
		DB: db,
	}
}

func NewDetailPesananController(db *gorm.DB) *DetailPesananController {
	return &DetailPesananController{
		DB: db,
	}
}

func (pc *PesananController) CreatePesanan(c *gin.Context) {

	var pesanan models.Pesanan
	if err := c.ShouldBindJSON(&pesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := pc.DB.Begin()
	defer tx.Commit()
	if err := tx.Create(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": pesanan})

}

func (pc *DetailPesananController) CreateDetailPesananById(c *gin.Context) {

	var req models.TambahPesananRequest
	pesananId := c.Param("pesanan_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := pc.DB.Begin()
	defer tx.Commit()
	var pesanan models.Pesanan
	if err := tx.Preload("DetailPesanan").First(&pesanan, "id = ?", pesananId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if add := tx.Create(&pesanan).Error; add != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": add.Error()})
				return
			}
			var produk models.Produk

			if err := tx.Preload("Inventaris").First(&produk, "nama = ?", req.ProdukName).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if produk.Inventaris.Jumlah < req.Jumlah {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Stok produk tidak mencukupi"})
				return
			}

			produk.Inventaris.Jumlah -= req.Jumlah
			if err := tx.Save(&produk.Inventaris).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			detailPesanan := models.DetailPesanan{
				PesananId: pesanan.ID,
				ProdukId:  produk.ID,
				Jumlah:    req.Jumlah,
			}

			if tx.Find(&detailPesanan, "pesanan_id = ? AND produk_id = ?", pesanan.ID, produk.ID).RowsAffected > 0 {
				detailPesanan.Jumlah += req.Jumlah
				if err := tx.Save(&detailPesanan).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				if err := tx.Create(&detailPesanan).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

		} else {

			var produk models.Produk
			if err := tx.Preload("Inventaris").First(&produk, "nama = ?", req.ProdukName).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if produk.Inventaris.Jumlah < req.Jumlah {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Stok produk tidak mencukupi"})
				return
			}

			produk.Inventaris.Jumlah -= req.Jumlah
			if err := tx.Save(&produk.Inventaris).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			detailPesanan := models.DetailPesanan{
				PesananId: pesanan.ID,
				ProdukId:  produk.ID,
				Jumlah:    req.Jumlah,
			}

			if tx.Find(&detailPesanan, "pesanan_id = ? AND produk_id = ?", pesanan.ID, produk.ID).RowsAffected > 0 {
				detailPesanan.Jumlah += req.Jumlah
				if err := tx.Save(&detailPesanan).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				if err := tx.Create(&detailPesanan).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Berhasil menambahkan pesanan"})
}

func (pc *PesananController) GetDetailPesananById(c *gin.Context) {

	var pesanan models.Pesanan
	pesananId := c.Param("pesanan_id")

	if err := pc.DB.Preload("DetailPesanan.Produk").First(&pesanan, "id = ?", pesananId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detailpesanan := make([]models.DetailPesananResponse, len(pesanan.DetailPesanan))
	for i, res := range pesanan.DetailPesanan {
		if res.Produk == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Produk not found"})
			return
		}
		detailpesanan[i] = models.DetailPesananResponse{
			Nama:   res.Produk.Nama,
			Jumlah: res.Jumlah,
			Harga:  res.Produk.Harga,
			Total:  res.Jumlah * res.Produk.Harga,
		}
	}
	response := models.PesananResponse{
		PesananId:     pesanan.ID,
		DetailPesanan: detailpesanan,
	}

	response.Total = 0
	for _, detail := range response.DetailPesanan {
		response.Total += detail.Total
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
