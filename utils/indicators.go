package utils

// RSI (Relative Strength Index) calcula el índice de fuerza relativa
func RSI(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0.0
	}

	gains, losses := 0.0, 0.0
	for i := 1; i < period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	averageGain := gains / float64(period)
	averageLoss := losses / float64(period)

	if averageLoss == 0 {
		return 100.0
	}

	rs := averageGain / averageLoss
	return 100.0 - (100.0 / (1.0 + rs))
}

// EMA (Exponential Moving Average) calcula la media móvil exponencial
func EMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	multiplier := 2.0 / float64(period+1)
	ema := make([]float64, len(prices))

	// Primera EMA basada en el promedio simple
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calcular el resto de la EMA
	for i := period; i < len(prices); i++ {
		ema[i] = ((prices[i] - ema[i-1]) * multiplier) + ema[i-1]
	}

	return ema[period-1:]
}

// SMA (Simple Moving Average) calcula la media móvil simple
func SMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	sma := make([]float64, len(prices)-period+1)
	for i := 0; i < len(sma); i++ {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += prices[j]
		}
		sma[i] = sum / float64(period)
	}

	return sma
}
