package strategies

import (
	"fmt"
	"log"
	"math"
	"trading-bot/services"
)

func GridTrading(binanceClient *services.BinanceClient, telegramClient *services.TelegramClient, baseCurrency, quoteCurrency string, tradeAmount float64) {
	symbol := fmt.Sprintf("%s%s", baseCurrency, quoteCurrency)

	gridLevels := 10
	gridRange := 0.05

	price, err := binanceClient.GetPrice(symbol)
	if err != nil {
		log.Fatalf("Error obteniendo precio para %s: %v", symbol, err)
	}

	log.Printf("Estrategia Grid Trading iniciada para %s a precio %.2f", symbol, price)

	gridStep := price * gridRange * (1.0 / float64(gridLevels)) // Conversión explícita de gridLevels
	buyLevels := make([]float64, gridLevels)
	sellLevels := make([]float64, gridLevels)
	for i := 0; i < gridLevels; i++ {
		buyLevels[i] = price - gridStep*float64(i+1)
		sellLevels[i] = price + gridStep*float64(i+1)
	}

	log.Printf("Niveles de compra: %v", buyLevels)
	log.Printf("Niveles de venta: %v", sellLevels)

	for {
		currentPrice, err := binanceClient.GetPrice(symbol)
		if err != nil {
			log.Printf("Error obteniendo precio para %s: %v", symbol, err)
			continue
		}

		log.Printf("Precio actual de %s: %.2f", symbol, currentPrice)

		// Verificar niveles de compra
		for _, buyPrice := range buyLevels {
			if currentPrice <= buyPrice {
				quoteBalance, err := binanceClient.GetBalance(quoteCurrency)
				if err != nil {
					log.Printf("Error verificando balance de %s: %v", quoteCurrency, err)
					continue
				}

				if quoteBalance >= tradeAmount {
					err := binanceClient.Buy(symbol, tradeAmount/currentPrice)
					if err != nil {
						log.Printf("Error realizando compra: %v", err)
					} else {
						log.Printf("Orden de compra ejecutada a %.2f para %s", currentPrice, symbol)
						telegramClient.SendMessage(fmt.Sprintf("Compra ejecutada en %.2f para %s", currentPrice, symbol))
					}
				}
			}
		}

		// Verificar niveles de venta
		for _, sellPrice := range sellLevels {
			if currentPrice >= sellPrice {
				baseBalance, err := binanceClient.GetBalance(baseCurrency)
				if err != nil {
					log.Printf("Error verificando balance de %s: %v", baseCurrency, err)
					continue
				}

				if baseBalance > 0 {
					amountToSell := math.Min(baseBalance, tradeAmount/currentPrice)
					err := binanceClient.Sell(symbol, amountToSell)
					if err != nil {
						log.Printf("Error realizando venta: %v", err)
					} else {
						log.Printf("Orden de venta ejecutada a %.2f para %s", currentPrice, symbol)
						telegramClient.SendMessage(fmt.Sprintf("Venta ejecutada en %.2f para %s", currentPrice, symbol))
					}
				}
			}
		}

		// Esperar 5 segundos antes de la siguiente iteración
		binanceClient.Sleep(5)
	}
}
