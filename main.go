package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"trading-bot/services"
	"trading-bot/strategies"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	// Leer variables de entorno
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	tradeAmountStr := os.Getenv("TRADE_AMOUNT")
	tradeAmount, err := strconv.ParseFloat(tradeAmountStr, 64)
	if err != nil {
		log.Fatalf("TRADE_AMOUNT debe ser un número: %v", err)
	}

	// Mostrar banner
	fmt.Println(`
===========================================================
      ██████╗ ███████╗████████╗
      ██╔══██╗██╔════╝╚══██╔══╝
      ██████╔╝█████╗     ██║   
      ██╔═══╝ ██╔══╝     ██║   
      ██║     ███████╗   ██║   
      ╚═╝     ╚══════╝   ╚═╝   
   Trading Bot - Juan Torres - GitHub: Juan1913
===========================================================
`)

	// Solicitar par de monedas al usuario
	var baseCurrency, quoteCurrency string
	fmt.Println("Ingrese la moneda base (ejemplo: BTC):")
	_, err = fmt.Scan(&baseCurrency)
	if err != nil || baseCurrency == "" {
		log.Fatalf("Error al leer la moneda base: %v", err)
	}

	fmt.Println("Ingrese la moneda de cotización (ejemplo: USDT):")
	_, err = fmt.Scan(&quoteCurrency)
	if err != nil || quoteCurrency == "" {
		log.Fatalf("Error al leer la moneda de cotización: %v", err)
	}

	symbol := baseCurrency + quoteCurrency
	fmt.Printf("Operando con el par: %s\n", symbol)

	// Inicializar servicios
	binanceClient := services.NewBinanceClient(apiKey, apiSecret)
	telegramClient := services.NewTelegramClient(os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"))

	// Mostrar menú de estrategias
	fmt.Println("Seleccione una estrategia:")
	fmt.Println("1. Grid Trading")
	fmt.Println("2. Breakout")
	fmt.Println("3. DCA")
	fmt.Println("4. Scalping")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		log.Fatalf("Error al leer la selección: %v", err)
	}

	// Ejecutar la estrategia seleccionada
	switch choice {
	case 1:
		strategies.GridTrading(binanceClient, telegramClient, baseCurrency, quoteCurrency, tradeAmount)
	case 2:
		strategies.BreakoutStrategy(binanceClient, telegramClient, baseCurrency, quoteCurrency, tradeAmount)
	case 3:
		strategies.DCAStrategy(binanceClient, telegramClient, baseCurrency, quoteCurrency, tradeAmount)
	case 4:
		strategies.ScalpingStrategy(binanceClient, telegramClient, baseCurrency, quoteCurrency, tradeAmount)
	default:
		fmt.Println("Selección no válida. Por favor, elija entre 1 y 4.")
	}
}
