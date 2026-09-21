<script setup lang="ts">
import type { RequestStatus } from '~/types'

const props = defineProps<{
  status: RequestStatus | string
}>()

const badgeConfig = computed(() => {
  switch (props.status) {
    case 'Approved':
      return {
        bg: 'bg-emerald-50 text-emerald-700 border-emerald-200',
        dot: 'bg-emerald-500',
        label: 'Approved'
      }
    case 'Rejected':
      return {
        bg: 'bg-rose-50 text-rose-700 border-rose-200',
        dot: 'bg-rose-500',
        label: 'Rejected'
      }
    case 'Submitted':
    default:
      return {
        bg: 'bg-amber-50 text-amber-700 border-amber-200',
        dot: 'bg-amber-500',
        label: 'Submitted'
      }
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border"
    :class="badgeConfig.bg"
  >
    <span class="w-1.5 h-1.5 rounded-full" :class="badgeConfig.dot"></span>
    {{ badgeConfig.label }}
  </span>
</template>

