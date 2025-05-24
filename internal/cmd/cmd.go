package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mclyashko/anomaly-detection-timeseries-go/internal/evaluator"
	"github.com/mclyashko/anomaly-detection-timeseries-go/internal/generator"
	"github.com/spf13/cobra"
)

var (
	startTime string
	endTime   string
	metric    string
	rules     string
	output    string
	delta     string
	jitter    string
)

var rootCmd = &cobra.Command{
	Use:   "anomaly-detection-timeseries-go",
	Short: "CLI для генерации временных рядов с метками аномалий",
	Run: func(cmd *cobra.Command, args []string) {
		start, err := time.Parse(time.DateTime, startTime)
		if err != nil {
			fmt.Println("Некорректный формат времени начала:", err)
			os.Exit(1)
		}
		end, err := time.Parse(time.DateTime, endTime)
		if err != nil {
			fmt.Println("Некорректный формат времени окончания:", err)
			os.Exit(1)
		}
		d, err := time.ParseDuration(delta)
		if err != nil {
			fmt.Println("Некорректный формат delta:", err)
			os.Exit(1)
		}
		j, err := time.ParseDuration(jitter)
		if err != nil {
			fmt.Println("Некорректный формат jitter:", err)
			os.Exit(1)
		}
		if j*2 >= d {
			fmt.Println("jitter должен быть меньше половины delta, чтобы гарантировать строго возрастающее время метрики")
			os.Exit(1)
		}

		f, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Ошибка открытия файла:", err)
			os.Exit(1)
		}
		defer f.Close()

		cfg := generator.GeneratorConfig{
			StartTime:  start,
			EndTime:    end,
			Delta:      d,
			Jitter:     j,
			MetricExpr: metric,
			RulesExpr:  rules,
			Output:     f,
		}

		evaluator, err := evaluator.NewExpressionEvaluator(metric, rules)
		if err != nil {
			fmt.Println("Ошибка инициализации вычислителя:", err)
			os.Exit(1)
		}

		gen := generator.NewGenerator(cfg, evaluator)
		if err := gen.Generate(); err != nil {
			fmt.Println("Ошибка генерации:", err)
			os.Exit(1)
		}

		fmt.Println("Генерация завершена")
	},
}

func init() {
	rootCmd.Flags().StringVarP(&startTime, "start", "s", defaultStartTime(), "Время начала (2006-01-02 15:04:05)")
	rootCmd.Flags().StringVarP(&endTime, "end", "e", defaultEndTime(), "Время окончания (2006-01-02 15:04:05)")
	rootCmd.Flags().StringVarP(&metric, "metric", "m", "5 + sin(t)", "Функция генерации метрики")
	rootCmd.Flags().StringVarP(&rules, "rules", "r", "value > 5", "Правила аномалий")
	rootCmd.Flags().StringVarP(&output, "output", "o", "data/output.csv", "Выходной CSV файл")
	rootCmd.Flags().StringVarP(&delta, "delta", "d", "1s", "Дельта между точками времени")
	rootCmd.Flags().StringVarP(&jitter, "jitter", "j", "0s", "Максимальный джиттер (случайное смещение времени)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func defaultStartTime() string {
	now := time.Now()
	offset := (int(now.Weekday()) + 6) % 7
	weekStart := now.AddDate(0, 0, -offset)
	return weekStart.Format(time.DateTime)
}

func defaultEndTime() string {
	now := time.Now()
	offset := 6 - (int(now.Weekday())+6)%7
	weekEnd := now.AddDate(0, 0, offset)
	return weekEnd.Format(time.DateTime)
}
