<script setup lang="ts">
import type { StatusLog } from '~/types'

defineProps<{
  logs: StatusLog[]
}>()

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="flow-root">
    <div v-if="!logs || logs.length === 0" class="text-sm text-slate-500 italic py-2">
      No status change history recorded yet.
    </div>

    <ul v-else role="list" class="-mb-8">
      <li v-for="(log, logIdx) in logs" :key="log.id">
        <div class="relative pb-8">
          <span
            v-if="logIdx !== logs.length - 1"
            class="absolute left-4 top-4 -ml-px h-full w-0.5 bg-slate-200"
            aria-hidden="true"
          />
          <div class="relative flex space-x-3">
            <div>
              <span
                class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white"
                :class="{
                  'bg-emerald-500 text-white': log.to_status === 'Approved',
                  'bg-rose-500 text-white': log.to_status === 'Rejected',
                  'bg-amber-500 text-white': log.to_status === 'Submitted'
                }"
              >
                <span class="text-xs font-bold">{{ log.to_status ? log.to_status[0] : 'S' }}</span>
              </span>
            </div>
            <div class="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
              <div>
                <p class="text-sm text-slate-800">
                  <span class="font-semibold text-slate-900">{{ log.changer_name || 'System' }}</span>
                  <span v-if="!log.from_status"> created the request as </span>
                  <span v-else> changed status from <strong>{{ log.from_status }}</strong> to </span>
                  <span
                    class="font-semibold inline-block px-1.5 py-0.5 rounded text-xs"
                    :class="{
                      'bg-emerald-100 text-emerald-800': log.to_status === 'Approved',
                      'bg-rose-100 text-rose-800': log.to_status === 'Rejected',
                      'bg-amber-100 text-amber-800': log.to_status === 'Submitted'
                    }"
                  >
                    {{ log.to_status }}
                  </span>
                </p>
                <p v-if="log.note" class="mt-1 text-xs text-slate-600 bg-slate-50 p-2 rounded border border-slate-100 italic">
                  "{{ log.note }}"
                </p>
              </div>
              <div class="whitespace-nowrap text-right text-xs text-slate-400">
                {{ formatDate(log.created_at) }}
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

