<template>
  <div class="promql-workbench">
    <div class="header">
      <h1>PromQL Insights</h1>
      <div class="header-controls">
        <select v-model="selectedDatasource" class="datasource-select">
          <option value="default">Default Prometheus</option>
        </select>
        <div class="time-range">
          <input type="datetime-local" v-model="timeRange.start" />
          <span>to</span>
          <input type="datetime-local" v-model="timeRange.end" />
        </div>
        <input type="number" v-model="step" placeholder="Step (s)" class="step-input" />
        <button @click="runQuery" class="run-button">Run</button>
      </div>
    </div>

    <div class="main-content">
      <div class="editor-panel">
        <div class="editor-toolbar">
          <button @click="formatQuery" class="toolbar-button">Format</button>
          <button @click="parseQuery" class="toolbar-button">Parse</button>
        </div>
        <div ref="editorContainer" class="editor-container"></div>
      </div>

      <div class="analysis-panel">
        <div class="tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            :class="['tab-button', { active: activeTab === tab.id }]"
            @click="activeTab = tab.id"
          >
            {{ tab.label }}
          </button>
        </div>

        <div class="tab-content">
          <div v-if="activeTab === 'ast'" class="ast-view">
            <ASTTree v-if="parseResult?.ast" :ast="parseResult.ast" />
            <div v-else class="empty-state">No AST data available</div>
          </div>

          <div v-if="activeTab === 'dependencies'" class="dependencies-view">
            <Dependencies v-if="parseResult?.dependencies" :dependencies="parseResult.dependencies" />
            <div v-else class="empty-state">No dependencies data available</div>
          </div>

          <div v-if="activeTab === 'issues'" class="issues-view">
            <Issues v-if="parseResult?.issues" :issues="parseResult.issues" />
            <div v-else class="empty-state">No issues found</div>
          </div>

          <div v-if="activeTab === 'results'" class="results-view">
            <QueryResults v-if="queryResult" :result="queryResult" />
            <div v-else class="empty-state">Run a query to see results</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as monaco from 'monaco-editor'
import type { ParseResponse, QueryRangeResponse } from '@/types/api'
import ASTTree from '@/components/ASTTree.vue'
import Dependencies from '@/components/Dependencies.vue'
import Issues from '@/components/Issues.vue'
import QueryResults from '@/components/QueryResults.vue'

const editorContainer = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | null = null

const selectedDatasource = ref('default')
const timeRange = ref({
  start: '',
  end: ''
})
const step = ref(15)

const tabs = [
  { id: 'ast', label: 'AST' },
  { id: 'dependencies', label: 'Dependencies' },
  { id: 'issues', label: 'Issues' },
  { id: 'results', label: 'Results' }
]

const activeTab = ref('ast')
const parseResult = ref<ParseResponse | null>(null)
const queryResult = ref<QueryRangeResponse | null>(null)

onMounted(() => {
  if (editorContainer.value) {
    editor = monaco.editor.create(editorContainer.value, {
      value: 'sum(rate(http_requests_total{job="api"}[5m])) by (status)',
      language: 'promql',
      theme: 'vs-dark',
      minimap: { enabled: false },
      automaticLayout: true,
      fontSize: 14,
      lineNumbers: 'on',
      scrollBeyondLastLine: false
    })

    monaco.languages.register({ id: 'promql' })
    monaco.languages.setMonarchTokensProvider('promql', {
      keywords: ['sum', 'avg', 'min', 'max', 'count', 'rate', 'irate', 'increase'],
      operators: ['by', 'without', 'and', 'or', 'unless', 'on', 'ignoring'],
      tokenizer: {
        root: [
          [/[a-zA-Z_][a-zA-Z0-9_]*/, { cases: { '@keywords': 'keyword', '@default': 'identifier' } }],
          [/[{}()\[\]]/, '@brackets'],
          [/\d+/, 'number'],
          [/"/, 'string', '@string'],
          [/\/\/.*$/, 'comment']
        ],
        string: [
          [/[^"]+/, 'string'],
          [/"/, 'string', '@pop']
        ]
      }
    })
  }

  const now = new Date()
  const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000)
  timeRange.value = {
    start: formatDateTimeLocal(oneHourAgo),
    end: formatDateTimeLocal(now)
  }
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
})

function formatDateTimeLocal(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

async function parseQuery() {
  if (!editor) return

  const query = editor.getValue()
  try {
    const response = await fetch('/api/promql/parse', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        query,
        dialect: 'prometheus',
        options: { format: true, analyze: true }
      })
    })

    const data = await response.json()
    parseResult.value = data

    if (data.formatted) {
      editor.setValue(data.formatted)
    }
  } catch (error) {
    console.error('Parse error:', error)
  }
}

async function formatQuery() {
  await parseQuery()
}

async function runQuery() {
  if (!editor) return

  const query = editor.getValue()
  const start = new Date(timeRange.value.start).getTime() / 1000
  const end = new Date(timeRange.value.end).getTime() / 1000

  try {
    const response = await fetch('/api/query/range', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        query,
        start,
        end,
        step: step.value,
        datasource: selectedDatasource.value
      })
    })

    const data = await response.json()
    queryResult.value = data
    activeTab.value = 'results'
  } catch (error) {
    console.error('Query error:', error)
  }
}
</script>

<style scoped>
.promql-workbench {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1e1e1e;
  color: #d4d4d4;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.header h1 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.datasource-select,
.time-range input,
.step-input {
  background: #3c3c3c;
  border: 1px solid #555;
  color: #d4d4d4;
  padding: 6px 10px;
  border-radius: 4px;
  font-size: 13px;
}

.time-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-range span {
  font-size: 12px;
  color: #888;
}

.run-button {
  background: #0e639c;
  color: white;
  border: none;
  padding: 6px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.run-button:hover {
  background: #1177bb;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.editor-panel {
  width: 50%;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #3c3c3c;
}

.editor-toolbar {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.toolbar-button {
  background: #3c3c3c;
  color: #d4d4d4;
  border: 1px solid #555;
  padding: 4px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.toolbar-button:hover {
  background: #4c4c4c;
}

.editor-container {
  flex: 1;
  overflow: hidden;
}

.analysis-panel {
  width: 50%;
  display: flex;
  flex-direction: column;
}

.tabs {
  display: flex;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.tab-button {
  background: transparent;
  color: #888;
  border: none;
  padding: 10px 16px;
  cursor: pointer;
  font-size: 13px;
  border-bottom: 2px solid transparent;
}

.tab-button:hover {
  color: #d4d4d4;
}

.tab-button.active {
  color: #d4d4d4;
  border-bottom-color: #0e639c;
}

.tab-content {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.empty-state {
  color: #888;
  text-align: center;
  padding: 40px;
  font-size: 14px;
}
</style>
