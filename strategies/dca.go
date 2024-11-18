package strategies

import (
	"fmt"
	"log"
	"math"
	"time"
	"trading-bot/services"
)

func DCAStrategy(binanceClient *services.BinanceClient, telegramClient *services.TelegramClient, baseCurrency, quoteCurrency string, tradeAmount float64) {
	symbol := fmt.Sprintf("%s%s", baseCurrency, quoteCurrency)

	log.Printf("Iniciando Dollar Cost Averaging (DCA) para el par %s", symbol)

	for {
		// Obtener el precio actual
		price, err := binanceClient.GetPrice(symbol)
		if err != nil {
			log.Printf("Error obteniendo precio de %s: %v", symbol, err)
			continue
		}

		log.Printf("Precio actual de %s: %.2f", symbol, price)

		// Verificar el balance disponible
		quoteBalance, err := binanceClient.GetBalance(quoteCurrency)
		if err != nil {
			log.Printf("Error verificando balance de %s: %v", quoteCurrency, err)
			continue
		}

		// Calcular la cantidad a comprar
		amountToBuy := tradeAmount / price
		if tradeAmount > quoteBalance {
			log.Printf("Saldo insuficiente en %s para ejecutar compra de %.2f %s", quoteCurrency, amountToBuy, baseCurrency)
			telegramClient.SendMessage(fmt.Sprintf("Saldo insuficiente en %s. No se puede realizar la compra DCA.", quoteCurrency))
			time.Sleep(10 * time.Second) // Esperar antes de intentar nuevamente
			continue
		}

		// Realizar la compra
		adjustedAmount := math.Min(amountToBuy, tradeAmount/price) // Ajustar cantidad al límite del balance
		err = binanceClient.Buy(symbol, adjustedAmount)
		if err != nil {
			log.Printf("Error realizando compra de %.2f %s: %v", adjustedAmount, baseCurrency, err)
			telegramClient.SendMessage(fmt.Sprintf("Error ejecutando compra DCA en %s: %v", symbol, err))
			continue
		}

		log.Printf("Compra DCA ejecutada: %.2f %s a precio %.2f", adjustedAmount, baseCurrency, price)
		telegramClient.SendMessage(fmt.Sprintf("Compra ejecutada: %.2f %s en %s a precio %.2f", adjustedAmount, baseCurrency, symbol, price))

		// Esperar antes de la próxima compra
		log.Printf("Esperando 1 minuto para la próxima compra DCA...")
		time.Sleep(1 * time.Minute)
	}
}
