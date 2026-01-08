<template>
  <div class="ast-tree">
    <div v-if="ast" class="tree-container">
      <ASTNode :node="ast" :level="0" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ASTNode as ASTNodeType } from '@/types/api'

interface Props {
  ast: ASTNodeType
}

defineProps<Props>()
</script>

<style scoped>
.ast-tree {
  height: 100%;
  overflow: auto;
}

.tree-container {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
}
</style>
