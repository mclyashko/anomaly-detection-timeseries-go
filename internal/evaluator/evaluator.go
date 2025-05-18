package evaluator

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/Knetic/govaluate"
)

type ExpressionEvaluator interface {
	EvaluateMetric(t time.Time) (float64, error)
	EvaluateRule(value float64) (bool, error)
}

type expressionEvaluator struct {
	metricExpr *govaluate.EvaluableExpression
	rulesExpr  *govaluate.EvaluableExpression
}

var functions = map[string]govaluate.ExpressionFunction{
	"sin":  func(args ...interface{}) (interface{}, error) { return math.Sin(args[0].(float64)), nil },
	"cos":  func(args ...interface{}) (interface{}, error) { return math.Cos(args[0].(float64)), nil },
	"tan":  func(args ...interface{}) (interface{}, error) { return math.Tan(args[0].(float64)), nil },
	"sqrt": func(args ...interface{}) (interface{}, error) { return math.Sqrt(args[0].(float64)), nil },
	"pow": func(args ...interface{}) (interface{}, error) {
		return math.Pow(args[0].(float64), args[1].(float64)), nil
	},
	"abs":   func(args ...interface{}) (interface{}, error) { return math.Abs(args[0].(float64)), nil },
	"exp":   func(args ...interface{}) (interface{}, error) { return math.Exp(args[0].(float64)), nil },
	"log":   func(args ...interface{}) (interface{}, error) { return math.Log(args[0].(float64)), nil },
	"ceil":  func(args ...interface{}) (interface{}, error) { return math.Ceil(args[0].(float64)), nil },
	"floor": func(args ...interface{}) (interface{}, error) { return math.Floor(args[0].(float64)), nil },
	"if": func(args ...any) (any, error) {
		condition, ok := args[0].(bool)
		if !ok {
			return nil, fmt.Errorf("first argument of if must be boolean")
		}
		if condition {
			return args[1], nil
		}
		return args[2], nil
	},
	"hour": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Hour()), nil
	},
	"minute": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Minute()), nil
	},
	"second": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Second()), nil
	},
	"day": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Day()), nil
	},
	"month": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Month()), nil
	},
	"year": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Year()), nil
	},
	"weekday": func(args ...any) (any, error) {
		return float64(time.Unix(int64(args[0].(float64)), 0).Weekday()), nil
	},
	"rand": func(args ...any) (any, error) {
		return rand.Float64(), nil
	},
}

func NewExpressionEvaluator(metric string, rules string) (ExpressionEvaluator, error) {
	metricExpr, err := govaluate.NewEvaluableExpressionWithFunctions(metric, functions)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга метрики: %w", err)
	}

	rulesExpr, err := govaluate.NewEvaluableExpressionWithFunctions(rules, functions)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга правил: %w", err)
	}

	return &expressionEvaluator{
		metricExpr: metricExpr,
		rulesExpr:  rulesExpr,
	}, nil
}

func (e *expressionEvaluator) EvaluateMetric(t time.Time) (float64, error) {
	parameters := map[string]interface{}{
		"t": float64(t.Unix()),
	}
	result, err := e.metricExpr.Evaluate(parameters)
	if err != nil {
		return 0, err
	}

	val, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("метрика не является числом")
	}
	return val, nil
}

func (e *expressionEvaluator) EvaluateRule(value float64) (bool, error) {
	parameters := map[string]interface{}{
		"value": value,
	}
	result, err := e.rulesExpr.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	isAnomaly, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("правило аномалии не возвращает bool")
	}
	return isAnomaly, nil
}
