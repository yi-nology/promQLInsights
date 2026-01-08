<template>
  <div class="query-results">
    <div v-if="!result" class="empty-state">
      Run a query to see results
    </div>

    <div v-else class="results-container">
      <div class="view-toggle">
        <button
          :class="['toggle-button', { active: viewMode === 'chart' }]"
          @click="viewMode = 'chart'"
        >
          Chart
        </button>
        <button
          :class="['toggle-button', { active: viewMode === 'table' }]"
          @click="viewMode = 'table'"
        >
          Table
        </button>
      </div>

      <div v-if="viewMode === 'chart'" ref="chartContainer" class="chart-container"></div>

      <div v-if="viewMode === 'table'" class="table-container">
        <table class="results-table">
          <thead>
            <tr>
              <th v-for="key in allLabels" :key="key">{{ key }}</th>
              <th>Values</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(series, index) in result.data.result" :key="index">
              <td v-for="key in allLabels" :key="key">{{ series.metric[key] || '-' }}</td>
              <td class="values-cell">
                <div class="values-list">
                  <div v-for="(value, i) in series.values" :key="i" class="value-item">
                    <span class="timestamp">{{ formatTimestamp(value[0]) }}</span>
                    <span class="value">{{ value[1] }}</span>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { QueryRangeResponse } from '@/types/api'

interface Props {
  result: QueryRangeResponse | null
}

const props = defineProps<Props>()
const viewMode = ref<'chart' | 'table'>('chart')
const chartContainer = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const allLabels = computed(() => {
  if (!props.result?.data.result.length) return []
  const labels = new Set<string>()
  props.result.data.result.forEach(series => {
    Object.keys(series.metric).forEach(key => labels.add(key))
  })
  return Array.from(labels)
})

onMounted(() => {
  if (chartContainer.value) {
    chart = echarts.init(chartContainer.value)
  }
})

watch(() => props.result, () => {
  if (viewMode.value === 'chart') {
    nextTick(() => {
      renderChart()
    })
  }
})

watch(viewMode, () => {
  if (viewMode.value === 'chart') {
    nextTick(() => {
      if (chartContainer.value && !chart) {
        chart = echarts.init(chartContainer.value)
      }
      renderChart()
    })
  }
})

function renderChart() {
  if (!chart || !props.result) return

  const series = props.result.data.result.map(s => ({
    name: Object.entries(s.metric).map(([k, v]) => `${k}="${v}"`).join(', '),
    data: s.values.map(v => [v[0] * 1000, parseFloat(v[1])]),
    type: 'line',
    showSymbol: false
  }))

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
    },
    legend: {
      type: 'scroll',
      top: 10,
      textStyle: { color: '#d4d4d4' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'time',
      axisLine: { lineStyle: { color: '#555' } },
      axisLabel: { color: '#d4d4d4' }
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#555' } },
      axisLabel: { color: '#d4d4d4' },
      splitLine: { lineStyle: { color: '#333' } }
    },
    series
  }

  chart.setOption(option)
}

function formatTimestamp(timestamp: number): string {
  const date = new Date(timestamp * 1000)
  return date.toLocaleTimeString()
}
</script>

<style scoped>
.query-results {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #888;
  font-size: 14px;
}

.results-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.view-toggle {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.toggle-button {
  background: #3c3c3c;
  color: #888;
  border: none;
  padding: 6px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.toggle-button:hover {
  background: #4c4c4c;
}

.toggle-button.active {
  background: #0e639c;
  color: #d4d4d4;
}

.chart-container {
  flex: 1;
  min-height: 400px;
}

.table-container {
  flex: 1;
  overflow: auto;
}

.results-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.results-table th,
.results-table td {
  padding: 8px 12px;
  text-align: left;
  border-bottom: 1px solid #3c3c3c;
}

.results-table th {
  background: #252526;
  color: #d4d4d4;
  font-weight: 600;
  position: sticky;
  top: 0;
}

.results-table td {
  color: #d4d4d4;
}

.values-cell {
  max-width: 400px;
}

.values-list {
  max-height: 200px;
  overflow: auto;
}

.value-item {
  display: flex;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px solid #2d2d2d;
}

.timestamp {
  color: #888;
  min-width: 80px;
}

.value {
  color: #9cdcfe;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}
</style>
