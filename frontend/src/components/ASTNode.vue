<template>
  <div class="ast-node" :class="{ collapsed }" :style="{ paddingLeft: `${level * 20}px` }">
    <div class="node-header" @click="toggle">
      <span class="toggle-icon">{{ collapsed ? '▶' : '▼' }}</span>
      <span class="node-type">{{ node.type }}</span>
      <span v-if="node.text" class="node-text">: {{ node.text }}</span>
    </div>
    <div v-if="!collapsed" class="node-children">
      <div v-for="(value, key) in nodeChildren" :key="key" class="node-child">
        <span class="child-key">{{ key }}:</span>
        <ASTNode
          v-if="isASTNode(value)"
          :node="value"
          :level="0"
        />
        <span v-else class="child-value">{{ formatValue(value) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ASTNode as ASTNodeType } from '@/types/api'

interface Props {
  node: ASTNodeType
  level: number
}

const props = defineProps<Props>()
const collapsed = ref(false)

const nodeChildren = computed(() => {
  const { type, span, text, meta, ...rest } = props.node
  return rest
})

function toggle() {
  collapsed.value = !collapsed.value
}

function isASTNode(value: any): boolean {
  return value && typeof value === 'object' && value.type
}

function formatValue(value: any): string {
  if (Array.isArray(value)) {
    return `[${value.map(formatValue).join(', ')}]`
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}
</script>

<style scoped>
.ast-node {
  margin: 2px 0;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
}

.node-header:hover {
  background: #2d2d2d;
}

.toggle-icon {
  color: #888;
  font-size: 10px;
  width: 12px;
}

.node-type {
  color: #569cd6;
  font-weight: 600;
}

.node-text {
  color: #ce9178;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.node-children {
  margin-top: 4px;
}

.node-child {
  display: flex;
  gap: 8px;
  padding: 2px 0;
}

.child-key {
  color: #9cdcfe;
  min-width: 80px;
}

.child-value {
  color: #d4d4d4;
  flex: 1;
}
</style>
