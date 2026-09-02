package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/provider"
	"iot-rd-backend/internal/repository"
)

type ExchangeRateService struct {
	repo         *repository.ExchangeRateRepository
	provider     provider.ExchangeRateProvider
	activityRepo *repository.ActivityRepository
}

func NewExchangeRateService(
	repo *repository.ExchangeRateRepository,
	provider provider.ExchangeRateProvider,
	activityRepo *repository.ActivityRepository,
) *ExchangeRateService {
	return &ExchangeRateService{
		repo:         repo,
		provider:     provider,
		activityRepo: activityRepo,
	}
}

// GetRate returns the cached exchange rate if fresh (< 4 hours) or fetches from provider.
func (s *ExchangeRateService) GetRate(base, quote string, forceRefresh bool) (*model.ExchangeRate, error) {
	base = strings.ToUpper(strings.TrimSpace(base))
	quote = strings.ToUpper(strings.TrimSpace(quote))

	if base == "" {
		base = "USD"
	}
	if quote == "" {
		quote = "IDR"
	}

	if base == quote {
		return &model.ExchangeRate{
			ID:            "EXR-" + base + "-" + quote,
			BaseCurrency:  base,
			QuoteCurrency: quote,
			Rate:          1.0,
			Provider:      "System Parity",
			Source:        "Direct 1:1",
			FetchedAt:     time.Now(),
			EffectiveDate: time.Now().Format("2006-01-02"),
		}, nil
	}

	if !forceRefresh {
		cached, err := s.repo.FindRate(base, quote)
		if err == nil && cached != nil {
			// Cache valid for 4 hours
			if time.Since(cached.FetchedAt) < 4*time.Hour && cached.Rate > 0 {
				return cached, nil
			}
		}
	}

	// Fetch from provider
	rateVal, effDate, err := s.provider.GetLatestRate(base, quote)
	if err != nil {
		// Fallback to cache if available
		if cached, cErr := s.repo.FindRate(base, quote); cErr == nil && cached != nil {
			return cached, nil
		}
		rateVal = 16500.00
		effDate = time.Now().Format("2006-01-02")
	}

	record := &model.ExchangeRate{
		ID:            "EXR-" + base + "-" + quote,
		BaseCurrency:  base,
		QuoteCurrency: quote,
		Rate:          math.Round(rateVal*10000) / 10000,
		Provider:      s.provider.GetName(),
		Source:        fmt.Sprintf("Live FX (%s -> %s)", base, quote),
		FetchedAt:     time.Now(),
		EffectiveDate: effDate,
	}

	_ = s.repo.UpsertRate(record)
	return record, nil
}

// Convert converts an amount between two currencies using either AUTO or MANUAL exchange rate.
func (s *ExchangeRateService) Convert(amount float64, fromCur, toCur string, mode string, manualRate float64) (float64, float64, string, error) {
	from := strings.ToUpper(strings.TrimSpace(fromCur))
	to := strings.ToUpper(strings.TrimSpace(toCur))

	if from == "" {
		from = "USD"
	}
	if to == "" {
		to = "USD"
	}

	if from == to {
		return amount, 1.0, "Direct 1:1", nil
	}

	mode = strings.ToUpper(strings.TrimSpace(mode))
	if mode == "MANUAL" && manualRate > 0 {
		var converted float64
		if from == "USD" && to == "IDR" {
			converted = amount * manualRate
		} else if from == "IDR" && to == "USD" {
			converted = amount / manualRate
		} else {
			converted = amount * manualRate
		}
		return math.Round(converted*100) / 100, manualRate, "Manual Override", nil
	}

	// AUTO mode: retrieve rate
	rateObj, err := s.GetRate(from, to, false)
	if err != nil {
		return amount, 1.0, "Fallback", err
	}

	converted := math.Round((amount*rateObj.Rate)*100) / 100
	return converted, rateObj.Rate, rateObj.Provider, nil
}

// CalculatePricing processes an interactive pricing request and computes base markup, selling price, and IDR references.
func (s *ExchangeRateService) CalculatePricing(req model.PricingCalculationRequest) (*model.PricingCalculationResult, error) {
	pCur := strings.ToUpper(strings.TrimSpace(req.PurchaseCurrency))
	if pCur == "" {
		pCur = "USD"
	}
	sCur := strings.ToUpper(strings.TrimSpace(req.SellingCurrency))
	if sCur == "" {
		sCur = pCur
	}

	margin := req.ProfitMarginPercentage
	if margin < 0 {
		margin = 0
	}

	// 1. Base selling price in Purchase Currency = PurchasePrice * (1 + margin / 100)
	baseSellingPrice := math.Round((req.PurchasePrice*(1.0+(margin/100.0)))*10000) / 10000

	// 2. Resolve Effective Exchange Rate (Primary Reference: USD <-> IDR)
	mode := strings.ToUpper(strings.TrimSpace(req.ExchangeRateMode))
	if mode == "" {
		mode = "AUTO"
	}

	liveRateObj, _ := s.GetRate("USD", "IDR", false)
	liveRate := 17749.00
	providerName := "ExchangeRate-API (Live)"
	var updatedAt *time.Time
	if liveRateObj != nil {
		liveRate = liveRateObj.Rate
		providerName = liveRateObj.Provider
		updatedAt = &liveRateObj.FetchedAt
	}

	manualRate := req.ManualExchangeRate
	if manualRate <= 0 {
		manualRate = liveRate
	}

	effectiveRate := liveRate
	if mode == "MANUAL" {
		effectiveRate = manualRate
		providerName = "Manual Override"
	}

	// 3. Calculate Final Selling Price in Selling Currency
	var sellingPrice float64
	if pCur == sCur {
		sellingPrice = math.Round(baseSellingPrice*100) / 100
	} else if pCur == "USD" && sCur == "IDR" {
		sellingPrice = math.Round((baseSellingPrice*effectiveRate)*100) / 100
	} else if pCur == "IDR" && sCur == "USD" {
		sellingPrice = math.Round((baseSellingPrice/effectiveRate)*100) / 100
	} else {
		// Generic cross rate
		conv, _, _, _ := s.Convert(baseSellingPrice, pCur, sCur, mode, effectiveRate)
		sellingPrice = conv
	}

	// 4. Calculate IDR Equivalents
	var purchaseInIDR, sellingInIDR float64
	if pCur == "IDR" {
		purchaseInIDR = req.PurchasePrice
	} else if pCur == "USD" {
		purchaseInIDR = math.Round((req.PurchasePrice * effectiveRate) * 100) / 100
	} else {
		purchaseInIDR = s.GetIDREquivalent(req.PurchasePrice, pCur)
	}

	if sCur == "IDR" {
		sellingInIDR = sellingPrice
	} else if sCur == "USD" {
		sellingInIDR = math.Round((sellingPrice * effectiveRate) * 100) / 100
	} else {
		sellingInIDR = s.GetIDREquivalent(sellingPrice, sCur)
	}

	return &model.PricingCalculationResult{
		PurchasePrice:          req.PurchasePrice,
		PurchaseCurrency:       pCur,
		ProfitMarginPercentage: margin,
		BaseSellingPrice:       baseSellingPrice,
		SellingPrice:           sellingPrice,
		SellingCurrency:        sCur,
		ExchangeRateMode:       mode,
		LiveExchangeRate:       liveRate,
		ManualExchangeRate:     manualRate,
		EffectiveExchangeRate:  effectiveRate,
		ExchangeRateProvider:   providerName,
		ExchangeRateUpdatedAt:  updatedAt,
		PurchasePriceInIDR:     purchaseInIDR,
		SellingPriceInIDR:      sellingInIDR,
	}, nil
}

// CalculateTradingDetailPricing enriches and recalculates a TradingProductDetail model.
func (s *ExchangeRateService) CalculateTradingDetailPricing(detail *model.TradingProductDetail) error {
	if detail == nil {
		return nil
	}

	if detail.PurchaseCurrency == "" {
		if detail.Currency != "" {
			detail.PurchaseCurrency = detail.Currency
		} else {
			detail.PurchaseCurrency = "USD"
		}
	}
	if detail.SellingCurrency == "" {
		detail.SellingCurrency = detail.PurchaseCurrency
	}
	if detail.ExchangeRateMode == "" {
		detail.ExchangeRateMode = "AUTO"
	}
	if detail.ProfitMarginPercentage <= 0 && detail.PurchasePrice > 0 && detail.SellingPrice > 0 {
		// Infer margin if existing record
		detail.ProfitMarginPercentage = math.Round(((detail.SellingPrice-detail.PurchasePrice)/detail.PurchasePrice)*10000) / 100
	} else if detail.ProfitMarginPercentage <= 0 {
		detail.ProfitMarginPercentage = 25.00
	}

	calc, err := s.CalculatePricing(model.PricingCalculationRequest{
		PurchasePrice:          detail.PurchasePrice,
		PurchaseCurrency:       detail.PurchaseCurrency,
		ProfitMarginPercentage: detail.ProfitMarginPercentage,
		SellingCurrency:        detail.SellingCurrency,
		ExchangeRateMode:       detail.ExchangeRateMode,
		ManualExchangeRate:     detail.ManualExchangeRate,
	})
	if err != nil {
		return err
	}

	detail.SellingPrice = calc.SellingPrice
	detail.LiveExchangeRate = calc.LiveExchangeRate
	if detail.ManualExchangeRate <= 0 {
		detail.ManualExchangeRate = calc.LiveExchangeRate
	}
	detail.EffectiveExchangeRate = calc.EffectiveExchangeRate
	detail.ExchangeRateProvider = calc.ExchangeRateProvider
	detail.ExchangeRateUpdatedAt = calc.ExchangeRateUpdatedAt
	detail.PurchasePriceInIDR = calc.PurchasePriceInIDR
	detail.SellingPriceInIDR = calc.SellingPriceInIDR
	detail.Currency = detail.PurchaseCurrency // sync legacy field

	return nil
}

// GetIDREquivalent converts any amount in a given currency to its estimated IDR value.
func (s *ExchangeRateService) GetIDREquivalent(amount float64, currency string) float64 {
	curr := strings.ToUpper(strings.TrimSpace(currency))
	if curr == "IDR" {
		return amount
	}
	rateObj, err := s.GetRate(curr, "IDR", false)
	if err == nil && rateObj != nil && rateObj.Rate > 0 {
		return math.Round((amount*rateObj.Rate)*100) / 100
	}
	// Fallback standard rate
	return math.Round((amount*17749.00)*100) / 100
}
