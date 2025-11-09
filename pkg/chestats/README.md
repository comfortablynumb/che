# chestats - Statistical Functions

Common statistical calculations for numeric data with zero dependencies.

## Features

- Mean, Median, Mode
- Variance and Standard Deviation
- Percentiles and Quartiles
- Min, Max, Range, Sum
- Correlation
- Generic support for all numeric types

## Quick Start

```go
data := []float64{1, 2, 3, 4, 5}

mean := chestats.Mean(data)         // 3.0
median := chestats.Median(data)     // 3.0
stdDev := chestats.StdDev(data)     // ~1.41

q1, q2, q3 := chestats.Quartiles(data)
p95 := chestats.Percentile(data, 95)
```

## API Reference

- `Mean[T Number](values []T) float64`
- `Median[T Number](values []T) float64`
- `Mode[T Number](values []T) T`
- `Variance[T Number](values []T) float64`
- `SampleVariance[T Number](values []T) float64`
- `StdDev[T Number](values []T) float64`
- `SampleStdDev[T Number](values []T) float64`
- `Min[T Number](values []T) T`
- `Max[T Number](values []T) T`
- `Sum[T Number](values []T) T`
- `Percentile[T Number](values []T, p float64) float64`
- `Quartiles[T Number](values []T) (q1, q2, q3 float64)`
- `IQR[T Number](values []T) float64`
- `Range[T Number](values []T) T`
- `Correlation[T Number](x, y []T) float64`
