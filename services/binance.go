package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	client *binance.Client
}

func NewBinanceClient(apiKey, apiSecret string) *BinanceClient {
	client := binance.NewClient(apiKey, apiSecret)
	return &BinanceClient{client: client}
}

func (b *BinanceClient) GetPrice(symbol string) (float64, error) {
	price, err := b.client.NewAveragePriceService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return 0, err
	}
	log.Printf("Precio actual de %s: %s", symbol, price.Price)
	return strconv.ParseFloat(price.Price, 64)
}

func (b *BinanceClient) GetMinQuantity(symbol string) (minQty float64, stepSize float64, err error) {
	exchangeInfo, err := b.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, 0, err
	}

	for _, s := range exchangeInfo.Symbols {
		if s.Symbol == symbol {
			for _, filter := range s.Filters {
				if filter["filterType"] == "LOT_SIZE" {
					minQty, err = strconv.ParseFloat(filter["minQty"].(string), 64)
					if err != nil {
						return 0, 0, fmt.Errorf("error parseando minQty: %v", err)
					}
					stepSize, err = strconv.ParseFloat(filter["stepSize"].(string), 64)
					if err != nil {
						return 0, 0, fmt.Errorf("error parseando stepSize: %v", err)
					}
					return minQty, stepSize, nil
				}
			}
		}
	}
	return 0, 0, fmt.Errorf("símbolo no encontrado o sin filtro LOT_SIZE: %s", symbol)
}

func (b *BinanceClient) AdjustQuantityForLotSize(qty, stepSize float64) float64 {
	return math.Floor(qty/stepSize) * stepSize
}

func (b *BinanceClient) Buy(symbol string, amount float64) error {
	minQty, stepSize, err := b.GetMinQuantity(symbol)
	if err != nil {
		return fmt.Errorf("error obteniendo restricciones de LOT_SIZE: %v", err)
	}

	adjustedQty := b.AdjustQuantityForLotSize(amount, stepSize)
	if adjustedQty < minQty {
		return fmt.Errorf("cantidad ajustada %.6f menor al mínimo permitido %.6f", adjustedQty, minQty)
	}

	_, err = b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.6f", adjustedQty)).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("error ejecutando orden de compra: %v", err)
	}

	log.Printf("Orden de compra ejecutada: %s, cantidad: %.6f", symbol, adjustedQty)
	return nil
}

func (b *BinanceClient) Sell(symbol string, amount float64) error {
	minQty, stepSize, err := b.GetMinQuantity(symbol)
	if err != nil {
		return fmt.Errorf("error obteniendo restricciones de LOT_SIZE: %v", err)
	}

	adjustedQty := b.AdjustQuantityForLotSize(amount, stepSize)
	if adjustedQty < minQty {
		return fmt.Errorf("cantidad ajustada %.6f menor al mínimo permitido %.6f", adjustedQty, minQty)
	}

	_, err = b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.6f", adjustedQty)).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("error ejecutando orden de venta: %v", err)
	}

	log.Printf("Orden de venta ejecutada: %s, cantidad: %.6f", symbol, adjustedQty)
	return nil
}

func (b *BinanceClient) GetBalance(asset string) (float64, error) {
	account, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("error obteniendo balance: %v", err)
	}

	for _, balance := range account.Balances {
		if balance.Asset == asset {
			available, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, fmt.Errorf("error parsing balance: %v", err)
			}
			return available, nil
		}
	}
	return 0, fmt.Errorf("activo no encontrado: %s", asset)
}

func (b *BinanceClient) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
