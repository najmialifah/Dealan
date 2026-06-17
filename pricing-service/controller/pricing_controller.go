package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/pricing-service/models"
	"github.com/najmialifah/Dealan/pricing-service/repository"
	"github.com/najmialifah/Dealan/pricing-service/service"
)

// PricingController menangani request HTTP untuk kalkulasi harga dan negosiasi.
type PricingController struct {
	pricingService service.PricingService
	repo           repository.PricingRepository
}

// NewPricingController membuat instance baru untuk HTTP Handler Pricing.
func NewPricingController(pricingService service.PricingService, repo repository.PricingRepository) *PricingController {
	return &PricingController{
		pricingService: pricingService,
		repo:           repo,
	}
}

// CalculateEstimate menangani POST /pricing/estimate untuk mengkalkulasi harga dinamis.
func (ctrl *PricingController) CalculateEstimate(c *gin.Context) {
	var req models.PricingEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	res, err := ctrl.pricingService.CalculateEstimate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// NegotiatePrice menangani POST /pricing/negotiate untuk memvalidasi dan menyimpan negosiasi harga.
func (ctrl *PricingController) NegotiatePrice(c *gin.Context) {
	var req models.NegotiationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	res, err := ctrl.pricingService.NegotiatePrice(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetRules menangani GET /pricing/rules/:service_type untuk melihat aturan tarif aktif.
func (ctrl *PricingController) GetRules(c *gin.Context) {
	serviceType := c.Param("service_type")
	if serviceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe layanan wajib diisi"})
		return
	}

	rule, err := ctrl.repo.GetRuleByServiceType(c.Request.Context(), serviceType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aturan tarif tidak ditemukan untuk tipe layanan ini"})
		return
	}

	c.JSON(http.StatusOK, rule)
}
