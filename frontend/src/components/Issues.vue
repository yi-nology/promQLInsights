<template>
  <div class="issues">
    <div v-if="issues.length === 0" class="no-issues">
      <div class="icon">✓</div>
      <div class="text">No issues found</div>
    </div>

    <div v-else class="issue-list">
      <div
        v-for="(issue, index) in issues"
        :key="index"
        :class="['issue-item', `issue-${issue.level}`]"
        @click="highlightIssue(issue)"
      >
        <div class="issue-header">
          <span :class="['issue-level', `level-${issue.level}`]">
            {{ issue.level.toUpperCase() }}
          </span>
          <span class="issue-code">{{ issue.code }}</span>
        </div>
        <div class="issue-message">{{ issue.message }}</div>
        <div v-if="issue.span" class="issue-span">
          Position: {{ issue.span[0] }} - {{ issue.span[1] }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Issue as IssueType } from '@/types/api'

interface Props {
  issues: IssueType[]
}

const props = defineProps<Props>()

function highlightIssue(issue: IssueType) {
  console.log('Highlight issue:', issue)
}
</script>

<style scoped>
.issues {
  height: 100%;
  overflow: auto;
}

.no-issues {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #6a9955;
}

.no-issues .icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-issues .text {
  font-size: 14px;
}

.issue-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.issue-item {
  padding: 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.issue-item:hover {
  background: #2d2d2d;
}

.issue-error {
  border-left: 3px solid #f14c4c;
}

.issue-warn {
  border-left: 3px solid #cca700;
}

.issue-info {
  border-left: 3px solid #3794ff;
}

.issue-header {
  display: flex;
  gap: 12px;
  margin-bottom: 6px;
}

.issue-level {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 3px;
}

.level-error {
  background: #f14c4c;
  color: #d4d4d4;
}

.level-warn {
  background: #cca700;
  color: #1e1e1e;
}

.level-info {
  background: #3794ff;
  color: #d4d4d4;
}

.issue-code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #9cdcfe;
}

.issue-message {
  font-size: 13px;
  color: #d4d4d4;
  margin-bottom: 4px;
}

.issue-span {
  font-size: 11px;
  color: #888;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}
</style>
