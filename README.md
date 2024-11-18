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


## Características Principales

Interacción con Binance API: Compra, venta y monitoreo de precios en tiempo real.

## Estrategias de Trading:

## Grid Trading: 
Aprovecha las fluctuaciones del mercado colocando órdenes de compra y venta en un rango definido.

## Breakout: 
Detecta rompimientos de niveles de resistencia o soporte y actúa en consecuencia.

## Dollar Cost Averaging (DCA): 
Compra periódicamente una cantidad fija independientemente del precio.

## Scalping: 
Realiza operaciones rápidas para obtener ganancias pequeñas en un corto período de tiempo.

## Notificaciones en Telegram: 
Envía alertas personalizadas sobre las operaciones ejecutadas y el estado de las estrategias.

## Estructura Modular: 
Diseñado para facilitar la adición de nuevas estrategias o servicios.