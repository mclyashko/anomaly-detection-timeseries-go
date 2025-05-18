package generator

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/mclyashko/anomaly-detection-timeseries-go/internal/evaluator"
)

type TimeSeriesGenerator interface {
	Generate() error
}

type timeSeriesGenerator struct {
	cfg       GeneratorConfig
	evaluator evaluator.ExpressionEvaluator
}

type GeneratorConfig struct {
	StartTime  time.Time
	EndTime    time.Time
	Delta      time.Duration
	Jitter     time.Duration
	MetricExpr string
	RulesExpr  string
	Output     io.Writer
}

func NewGenerator(
	cfg GeneratorConfig,
	evaluator evaluator.ExpressionEvaluator,
) TimeSeriesGenerator {
	return &timeSeriesGenerator{
		cfg:       cfg,
		evaluator: evaluator,
	}
}

func (g *timeSeriesGenerator) Generate() error {
	curr := g.cfg.StartTime
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for !curr.After(g.cfg.EndTime) {
		jitter := time.Duration(rnd.Int63n(int64(g.cfg.Jitter)*2)) - g.cfg.Jitter
		ts := curr.Add(jitter)

		value, err := g.evaluator.EvaluateMetric(ts)
		if err != nil {
			return fmt.Errorf("ошибка вычисления метрики: %w", err)
		}

		isAnomaly, err := g.evaluator.EvaluateRule(value)
		if err != nil {
			return fmt.Errorf("ошибка проверки правила аномалии: %w", err)
		}

		line := fmt.Sprintf("%s,%.4f,%t\n", ts.Format(time.DateTime), value, isAnomaly)

		if _, err := g.cfg.Output.Write([]byte(line)); err != nil {
			return fmt.Errorf("write error: %w", err)
		}

		curr = curr.Add(g.cfg.Delta)
	}

	return nil
}
