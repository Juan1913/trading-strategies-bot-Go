package strategies

import (
	"fmt"
	"log"
	"time"
	"trading-bot/services"
)

func BreakoutStrategy(binanceClient *services.BinanceClient, telegramClient *services.TelegramClient, baseCurrency, quoteCurrency string, tradeAmount float64) {
	symbol := fmt.Sprintf("%s%s", baseCurrency, quoteCurrency)

	// Definir niveles iniciales de soporte y resistencia
	resistanceLevel := 240.0 // Nivel de resistencia
	supportLevel := 230.0    // Nivel de soporte

	log.Printf("Iniciando Breakout Strategy para el par %s", symbol)

	for {
		// Obtener el precio actual
		price, err := binanceClient.GetPrice(symbol)
		if err != nil {
			log.Printf("Error obteniendo precio de %s: %v", symbol, err)
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Precio actual de %s: %.2f", symbol, price)

		// Verificar si el precio rompe la resistencia
		if price > resistanceLevel {
			log.Printf("Breakout detectado por encima de la resistencia %.2f. Intentando comprar...", resistanceLevel)

			// Realizar la compra
			err := binanceClient.Buy(symbol, tradeAmount/price)
			if err != nil {
				log.Printf("Error realizando compra: %v", err)
				telegramClient.SendMessage(fmt.Sprintf("Error comprando en %s: %v", symbol, err))
			} else {
				telegramClient.SendMessage(fmt.Sprintf("Breakout detectado y compra ejecutada en %s:\n- Precio: %.2f\n- Resistencia: %.2f\n- Monto: %.2f %s", symbol, price, resistanceLevel, tradeAmount, quoteCurrency))
			}

			// Ajustar el nivel de resistencia después de la compra
			resistanceLevel = price * 1.02 // Incrementar el nivel de resistencia en un 2%
			log.Printf("Nuevo nivel de resistencia ajustado a %.2f", resistanceLevel)
		}

		// Verificar si el precio rompe el soporte
		if price < supportLevel {
			log.Printf("Breakout detectado por debajo del soporte %.2f. Intentando vender...", supportLevel)

			// Realizar la venta
			baseBalance, _ := binanceClient.GetBalance(baseCurrency)
			if baseBalance > 0 {
				err := binanceClient.Sell(symbol, baseBalance*0.5) // Vender la mitad del balance
				if err != nil {
					log.Printf("Error realizando venta: %v", err)
					telegramClient.SendMessage(fmt.Sprintf("Error vendiendo en %s: %v", symbol, err))
				} else {
					telegramClient.SendMessage(fmt.Sprintf("Breakout detectado y venta ejecutada en %s:\n- Precio: %.2f\n- Soporte: %.2f\n- Monto vendido: %.2f %s", symbol, price, supportLevel, baseBalance*0.5, baseCurrency))
				}
			} else {
				log.Printf("Sin saldo disponible en %s para vender.", baseCurrency)
			}

			// Ajustar el nivel de soporte después de la venta
			supportLevel = price * 0.98 // Reducir el nivel de soporte en un 2%
			log.Printf("Nuevo nivel de soporte ajustado a %.2f", supportLevel)
		}

		// Esperar antes de verificar nuevamente
		time.Sleep(10 * time.Second)
	}
}
