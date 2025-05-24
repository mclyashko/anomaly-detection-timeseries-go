#!/bin/bash
set -e

GENERATOR_CMD="go run cmd/tsgen/main.go"
VISUALIZER_CMD="go run cmd/csvviz/main.go"

OUTDIR="data"

# Очистка папки OUTDIR (если нет - создаст)
rm -rf "$OUTDIR"
mkdir -p "$OUTDIR"

# Массив 0: Постоянное значение с редкими пиками (имитация стабильной метрики с аномальными всплесками)
cmd_args_0=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-14 23:59:59"
  "--metric=5 + 10*if(hour(t) == 15, 1, 0)"
  "--rules=value > 6"
  "--output=$OUTDIR/constant_with_spikes.csv"
  "--delta=600s"
  "--jitter=15s"
)

# Массив 1: Синусоида с суточной периодичностью (имитация нагрузки с дневными/ночными циклами)
cmd_args_1=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-14 23:59:59"
  "--metric=5 + 4*sin(2*3.1415926535*hour(t)/24)"
  "--rules=value > 8"
  "--output=$OUTDIR/sine_daynight.csv"
  "--delta=600s"
  "--jitter=30s"
)

# Массив 2: Пилообразная метрика (имитация накопления с сбросом, например, счетчик запросов)
cmd_args_2=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-07 23:59:59"
  "--metric=floor(hour(t)) + if(hour(t) > 12, 24 - hour(t), 0)"
  "--rules=value > 15"
  "--output=$OUTDIR/sawtooth.csv"
  "--delta=3600s"
  "--jitter=1s"
)

# Массив 3: Линейное уменьшение с редкими пиками (имитация уменьшения нагрузки с аномалиями)
cmd_args_3=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-07 23:59:59"
  "--metric=10 - t/14400 + 10*if(and(hour(t) == 3, minute(t) == 0), 1, 0)"
  "--rules=value > -121255"
  "--output=$OUTDIR/linear_decay_with_spikes.csv"
  "--delta=3600s"
  "--jitter=1s"
)

# Массив 4: Периодические ступенчатые всплески (имитация резких изменений нагрузки)
cmd_args_4=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-07 23:59:59"
  "--metric=10 + 5*if(floor(t/7200) - 2*floor(t/14400) == 1, 1, 0)"
  "--rules=value > 12"
  "--output=$OUTDIR/step_spikes.csv"
  "--delta=1800s"
  "--jitter=300s"
)

# Массив 5: Синусоида с медленным линейным ростом (имитация нагрузки с трендом)
cmd_args_5=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-21 23:59:59"
  "--metric=5 + 3*sin(2*3.1415926535*hour(t)/24) + t/86400"
  "--rules=value > 20230"
  "--output=$OUTDIR/sine_plus_linear.csv"
  "--delta=600s"
  "--jitter=15s"
)

# Массив 6: Метрика с дневными и ночными колебаниями, шумом и регулярными аномалиями
cmd_args_6=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-07 23:59:59"
  "--metric=2 + 0.5*rand() 
+ 1.2*sin(2*3.1415926535*t/86400) \
+ 0.8*sin(2*3.1415926535*t/172800 + 3.1415926535/4) \
+ 3*if(and(minute(t) % 37 == 0, hour(t) % 5 == 0), 1, 0) \
+ 9*if(and(hour(t) == 10, minute(t) == 0), 1, 0) \
+ 1.5*if(and(minute(t) % 17 == 0, hour(t) % 4 == 0), 1, 0) \
+ 3*if(and(hour(t) >= 0, hour(t) <= 4), 1, 0) \
+ 5*if(and(minute(t) >= 30, minute(t) <= 45, hour(t) == 14), 1, 0)"
  "--rules=or( \
  and(hour(t) == 10, minute(t) == 0, value > 8), \
  and(hour(t) == 14, minute(t) >= 30, minute(t) <= 45, value > 4), \
  and(minute(t) % 37 == 0, hour(t) % 5 == 0, value > 2), \
  and(minute(t) % 17 == 0, hour(t) % 4 == 0, value > 1.5), \
  and(hour(t) >= 0, hour(t) <= 4, value > 5), \
  and(hour(t) >= 10, hour(t) <= 12, value > 5) \
  )"
  "--output=$OUTDIR/daily_cycles_with_noise_and_anomalies.csv"
  "--delta=600s"
  "--jitter=100s"
)

# Массив 7: Метрика с трендом и случайными пиками (имитация использования памяти)
cmd_args_7=(
  "--start=2025-05-01 00:00:00"
  "--end=2025-05-07 23:59:59"
  "--metric=10 + t/86400 + 15*if(rand() < 0.01, 1, 0)"
  "--rules=value > 20221"
  "--output=$OUTDIR/memory_with_spikes.csv"
  "--delta=4800s"
  "--jitter=1s"
)

# Массив, содержащий все массивы аргументов
CMD_ARGS_LIST=(
  cmd_args_0[@]
  cmd_args_1[@]
  cmd_args_2[@]
  cmd_args_3[@]
  cmd_args_4[@]
  cmd_args_5[@]
  cmd_args_6[@]
  cmd_args_7[@]
)

for ((i=0; i<${#CMD_ARGS_LIST[@]}; i++)); do
  echo "Генерация CSV файла #$((i+1))..."
  # Получаем массив аргументов по имени
  args=("${!CMD_ARGS_LIST[$i]}")
  # Выводим аргументы для отладки
  echo "Аргументы: ${args[@]}"
  # Выполняем команду с аргументами
  $GENERATOR_CMD "${args[@]}" || exit 1

  # Извлекаем путь до файла из массива аргументов
  output_file=""
  for arg in "${args[@]}"; do
    if [[ $arg == --output=* ]]; then
      output_file="${arg#--output=}"
      break
    fi
  done

  if [[ -z "$output_file" ]]; then
    echo "Ошибка: путь до файла (--output=) не найден в команде #$((i+1))."
    exit 1
  fi

  if [[ ! -f "$output_file" ]]; then
    echo "Ошибка: файл $output_file не был создан."
    exit 1
  fi

  echo "Визуализация файла #$((i+1))..."
  $VISUALIZER_CMD --file="$output_file" || exit 1
  echo "---------------------------"
done
