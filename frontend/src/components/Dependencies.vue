<template>
  <div class="dependencies">
    <div class="section">
      <h3>Metrics</h3>
      <div class="tag-list">
        <span v-for="metric in dependencies.metrics" :key="metric" class="tag metric-tag">
          {{ metric }}
        </span>
      </div>
    </div>

    <div class="section">
      <h3>Labels</h3>
      <div class="tag-list">
        <span v-for="label in dependencies.labels" :key="label" class="tag label-tag">
          {{ label }}
        </span>
      </div>
    </div>

    <div class="section">
      <h3>Functions</h3>
      <div class="tag-list">
        <span v-for="fn in dependencies.functions" :key="fn" class="tag function-tag">
          {{ fn }}()
        </span>
      </div>
    </div>

    <div class="section">
      <h3>Aggregations</h3>
      <div class="tag-list">
        <span v-for="agg in dependencies.aggregations" :key="agg" class="tag aggregation-tag">
          {{ agg }}()
        </span>
      </div>
    </div>

    <div class="section">
      <h3>Matchers</h3>
      <div class="matcher-list">
        <div v-for="(matcher, index) in dependencies.matchers" :key="index" class="matcher">
          <span class="matcher-label">{{ matcher.label }}</span>
          <span class="matcher-op">{{ matcher.op }}</span>
          <span class="matcher-value">"{{ matcher.value }}"</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Dependencies as DependenciesType } from '@/types/api'

interface Props {
  dependencies: DependenciesType
}

defineProps<Props>()
</script>

<style scoped>
.dependencies {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section h3 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #d4d4d4;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.metric-tag {
  background: #264f78;
  color: #9cdcfe;
}

.label-tag {
  background: #094771;
  color: #4ec9b0;
}

.function-tag {
  background: #569cd6;
  color: #d4d4d4;
}

.aggregation-tag {
  background: #6a9955;
  color: #d4d4d4;
}

.matcher-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.matcher {
  display: flex;
  gap: 8px;
  padding: 6px 10px;
  background: #2d2d2d;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
}

.matcher-label {
  color: #9cdcfe;
}

.matcher-op {
  color: #d4d4d4;
}

.matcher-value {
  color: #ce9178;
}
</style>
