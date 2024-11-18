package strategies

import (
	"fmt"
	"log"
	"trading-bot/services"
)

func ScalpingStrategy(binanceClient *services.BinanceClient, telegramClient *services.TelegramClient, baseCurrency, quoteCurrency string, tradeAmount float64) {
	symbol := fmt.Sprintf("%s%s", baseCurrency, quoteCurrency)

	for {
		// Obtener el precio actual
		price, err := binanceClient.GetPrice(symbol)
		if err != nil {
			log.Printf("Error obteniendo precio de %s: %v", symbol, err)
			continue
		}

		log.Printf("Precio actual de %s: %.2f", symbol, price)

		// Calcular Stop Loss y Take Profit
		stopLoss := price * 0.99
		takeProfit := price * 1.01

		message := fmt.Sprintf(
			"Scalping Strategy:\n- Mercado: %s\n- Precio actual: %.2f\n- Stop Loss: %.2f\n- Take Profit: %.2f",
			symbol, price, stopLoss, takeProfit,
		)
		err = telegramClient.SendMessage(message)
		if err != nil {
			log.Printf("Error enviando mensaje a Telegram: %v", err)
		}

		// Verificar balance en moneda de cotización (USDT)
		quoteBalance, err := binanceClient.GetBalance(quoteCurrency)
		if err != nil {
			log.Printf("Error verificando balance de %s: %v", quoteCurrency, err)
			continue
		}

		// Si el balance en USDT es insuficiente, intenta vender
		if quoteBalance < tradeAmount {
			log.Printf("Saldo insuficiente en %s. Intentando vender...", quoteCurrency)

			// Obtener balance en la moneda base (ETH)
			baseBalance, err := binanceClient.GetBalance(baseCurrency)
			if err != nil {
				log.Printf("Error verificando balance de %s: %v", baseCurrency, err)
				continue
			}

			if baseBalance > 0 {
				err := binanceClient.Sell(symbol, baseBalance*0.5) // Vende el 50% del balance
				if err != nil {
					log.Printf("Error al vender: %v", err)
				}
			} else {
				log.Printf("Saldo insuficiente en %s para vender. Esperando...", baseCurrency)
			}
		} else {
			// Si hay suficiente USDT, realiza una compra
			err := binanceClient.Buy(symbol, tradeAmount)
			if err != nil {
				log.Printf("Error al comprar: %v", err)
			}
		}

		// Esperar antes de la próxima iteración
		binanceClient.Sleep(5)
	}
}
