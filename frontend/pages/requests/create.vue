<script setup lang="ts">
import { ArrowLeft, Send, AlertCircle, Loader2 } from 'lucide-vue-next'
import type { Priority } from '~/types'

definePageMeta({
  middleware: 'auth'
})

const { createRequest } = useRequests()
const router = useRouter()

const form = reactive({
  asset_id: '',
  priority: 'Medium' as Priority,
  problem_description: ''
})

const priorities: Priority[] = ['Low', 'Medium', 'High', 'Urgent']

const isSubmitting = ref(false)
const errorMessage = ref('')

const handleSubmit = async () => {
  if (!form.asset_id.trim()) {
    errorMessage.value = 'Please provide an Asset / Machine ID.'
    return
  }
  if (!form.problem_description.trim()) {
    errorMessage.value = 'Please describe the problem in detail.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const created = await createRequest({
      asset_id: form.asset_id.trim(),
      priority: form.priority,
      problem_description: form.problem_description.trim()
    })
    router.push(`/requests/${created.id}`)
  } catch (err: any) {
    errorMessage.value = err?.data?.error || 'Failed to submit maintenance request'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <!-- Breadcrumb / Back button -->
    <div class="mb-6">
      <NuxtLink
        to="/requests"
        class="inline-flex items-center gap-1.5 text-sm font-medium text-slate-600 hover:text-slate-900 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        Back to Requests List
      </NuxtLink>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 sm:p-8">
      <div class="border-b border-slate-200 pb-5 mb-6">
        <h1 class="text-xl font-bold text-slate-900">Raise Maintenance Request</h1>
        <p class="mt-1 text-sm text-slate-500">
          Submit an incident or machinery issue. Supervisors will be notified for review.
        </p>
      </div>

      <!-- Error Alert -->
      <div
        v-if="errorMessage"
        class="mb-6 p-4 rounded-lg bg-rose-50 border border-rose-200 flex items-start gap-3"
      >
        <AlertCircle class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" />
        <div class="text-sm text-rose-700 font-medium leading-tight">
          {{ errorMessage }}
        </div>
      </div>

      <form class="space-y-6" @submit.prevent="handleSubmit">
        <!-- Asset ID -->
        <div>
          <label for="asset_id" class="block text-sm font-semibold text-slate-700">
            Machine or Asset ID <span class="text-rose-500">*</span>
          </label>
          <div class="mt-1">
            <input
              id="asset_id"
              v-model="form.asset_id"
              type="text"
              required
              placeholder="e.g. CNC-003, PRESS-LINE-A, PUMP-007"
              class="w-full px-3.5 py-2.5 text-sm border border-slate-300 rounded-lg shadow-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-factory-500 focus:border-factory-500 font-mono uppercase"
            />
          </div>
          <p class="mt-1 text-xs text-slate-500">Unique tag or machine identifier found on the factory asset plate.</p>
        </div>

        <!-- Priority Selection -->
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-2">
            Priority Level <span class="text-rose-500">*</span>
          </label>
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <label
              v-for="p in priorities"
              :key="p"
              class="relative flex cursor-pointer rounded-lg border p-3 shadow-sm focus:outline-none transition-all text-center"
              :class="form.priority === p 
                ? 'border-factory-600 ring-2 ring-factory-600 bg-factory-50/50' 
                : 'border-slate-200 bg-white hover:border-slate-300'"
            >
              <input
                type="radio"
                name="priority"
                :value="p"
                v-model="form.priority"
                class="sr-only"
              />
              <span class="w-full">
                <span class="block text-sm font-medium" :class="form.priority === p ? 'text-factory-900 font-bold' : 'text-slate-700'">
                  {{ p }}
                </span>
              </span>
            </label>
          </div>
        </div>

        <!-- Problem Description -->
        <div>
          <label for="description" class="block text-sm font-semibold text-slate-700">
            Problem Description <span class="text-rose-500">*</span>
          </label>
          <div class="mt-1">
            <textarea
              id="description"
              v-model="form.problem_description"
              rows="4"
              required
              placeholder="Describe symptoms, error codes, noise level, or conditions leading to the issue..."
              class="w-full px-3.5 py-2.5 text-sm border border-slate-300 rounded-lg shadow-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-factory-500 focus:border-factory-500"
            ></textarea>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="pt-4 border-t border-slate-200 flex items-center justify-end gap-3">
          <NuxtLink
            to="/requests"
            class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
          >
            Cancel
          </NuxtLink>

          <button
            type="submit"
            :disabled="isSubmitting"
            class="inline-flex items-center gap-2 px-5 py-2 text-sm font-medium text-white bg-factory-600 hover:bg-factory-700 rounded-lg shadow-sm transition-colors disabled:opacity-50"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <Send v-else class="w-4 h-4" />
            <span>{{ isSubmitting ? 'Submitting...' : 'Submit Request' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

