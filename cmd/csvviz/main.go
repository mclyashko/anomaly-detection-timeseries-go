package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/guptarohit/asciigraph"
)

type DataPoint struct {
	Timestamp time.Time
	Value     float64
	Anomaly   bool
}

func main() {
	filePath := flag.String("file", "", "Путь до CSV файла с данными")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Пожалуйста, укажите путь до CSV файла через --file")
		os.Exit(1)
	}

	data, err := readCSV(*filePath)
	if err != nil {
		fmt.Println("Ошибка чтения CSV:", err)
		os.Exit(1)
	}

	plotAll(data)
	plotNormal(data)
	plotAnomalies(data)
}

func readCSV(path string) ([]DataPoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataPoint
	for i, row := range rows {
		if len(row) < 3 {
			continue
		}

		ts, err := time.Parse(time.DateTime, row[0])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга времени в строке %d: %w", i+1, err)
		}

		val, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга значения в строке %d: %w", i+1, err)
		}

		anomaly, err := strconv.ParseBool(row[2])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга аномалии в строке %d: %w", i+1, err)
		}

		data = append(data, DataPoint{Timestamp: ts, Value: val, Anomaly: anomaly})
	}

	return data, nil
}

func plotAll(data []DataPoint) {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.Value
	}

	fmt.Println("\n=== График: Все данные ===")
	graph := asciigraph.Plot(values, asciigraph.Width(160), asciigraph.Height(30), asciigraph.Caption("Метрика (все точки)"))
	fmt.Println(graph)
}

func plotNormal(data []DataPoint) {
	values := make([]float64, len(data))
	for i, d := range data {
		if d.Anomaly {
			values[i] = math.NaN() // пропускаем аномалии
		} else {
			values[i] = d.Value
		}
	}

	fmt.Println("\n=== График: Только нормальные точки ===")
	graph := asciigraph.Plot(values, asciigraph.Width(160), asciigraph.Height(30), asciigraph.Caption("Метрика (без аномалий)"))
	fmt.Println(graph)
}

func plotAnomalies(data []DataPoint) {
	values := make([]float64, len(data))
	for i, d := range data {
		if d.Anomaly {
			values[i] = d.Value
		} else {
			values[i] = math.NaN() // пропускаем нормальные точки
		}
	}

	fmt.Println("\n=== График: Только аномалии ===")
	graph := asciigraph.Plot(values, asciigraph.Width(160), asciigraph.Height(30), asciigraph.Caption("Метрика (только аномалии)"))
	fmt.Println(graph)
}
