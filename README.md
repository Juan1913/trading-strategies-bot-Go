# Trading Bot en Go

Este proyecto implementa un bot de trading en Go que interactúa con Binance y envía alertas a Telegram.

## Estructura

- `main.go`: Punto de entrada principal.
- `strategies/`: Implementación de estrategias de trading.
- `services/`: Conexión con Binance y sistema de alertas.
- `models/`: Modelos de datos como órdenes y balances.
- `utils/`: Funciones auxiliares.

## Configuración

1. Crea un archivo `.env` con tus credenciales de Binance y Telegram.
2. Instala las dependencias necesarias:
   ```bash
   go mod tidy
