<script setup lang="ts">
import type { Priority } from '~/types'

const props = defineProps<{
  priority: Priority | string
}>()

const priorityConfig = computed(() => {
  switch (props.priority) {
    case 'Urgent':
      return {
        bg: 'bg-red-100 text-red-800 border-red-300 font-semibold',
        indicator: 'bg-red-600 animate-pulse'
      }
    case 'High':
      return {
        bg: 'bg-orange-100 text-orange-800 border-orange-200',
        indicator: 'bg-orange-500'
      }
    case 'Medium':
      return {
        bg: 'bg-blue-50 text-blue-700 border-blue-200',
        indicator: 'bg-blue-500'
      }
    case 'Low':
    default:
      return {
        bg: 'bg-slate-100 text-slate-700 border-slate-200',
        indicator: 'bg-slate-400'
      }
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-xs border"
    :class="priorityConfig.bg"
  >
    <span class="w-1.5 h-1.5 rounded-full" :class="priorityConfig.indicator"></span>
    {{ priority }}
  </span>
</template>

