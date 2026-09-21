<script setup lang="ts">
import { 
  ArrowLeft, Edit3, Trash2, CheckCircle, XCircle, 
  Calendar, User as UserIcon, AlertCircle, Save, Loader2 
} from 'lucide-vue-next'
import type { RequestDetail, Priority, RequestStatus } from '~/types'

definePageMeta({
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const requestId = route.params.id as string

const { user, isAdmin, canReview } = useAuth()
const { getRequest, updateRequest, reviewRequest, deleteRequest } = useRequests()

const request = ref<RequestDetail | null>(null)
const isLoading = ref(true)
const errorMessage = ref('')
const isEditMode = ref(route.query.edit === 'true')

// Edit form
const editForm = reactive({
  asset_id: '',
  problem_description: '',
  priority: 'Medium' as Priority,
  status: 'Submitted' as RequestStatus
})
const isSaving = ref(false)

// Review Modal
const reviewModal = reactive({
  isOpen: false,
  action: 'Approved' as 'Approved' | 'Rejected',
  note: '',
  isSubmitting: false
})

const fetchDetail = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const data = await getRequest(requestId)
    request.value = data
    // Pre-populate edit form
    editForm.asset_id = data.asset_id
    editForm.problem_description = data.problem_description
    editForm.priority = data.priority
    editForm.status = data.status
  } catch (err: any) {
    errorMessage.value = err?.data?.error || 'Failed to load request details'
  } finally {
    isLoading.value = false
  }
}

// Can user edit this request?
const canUserEdit = computed(() => {
  if (!request.value || !user.value) return false
  if (isAdmin.value) return true
  // Operator & Supervisor can only edit their own request if status is 'Submitted'
  return request.value.created_by === user.value.id && request.value.status === 'Submitted'
})

// Save edits
const handleSaveEdit = async () => {
  if (!editForm.asset_id.trim() || !editForm.problem_description.trim()) {
    alert('Asset ID and Problem Description are required.')
    return
  }

  isSaving.value = true
  try {
    const payload: any = {
      asset_id: editForm.asset_id.trim(),
      problem_description: editForm.problem_description.trim(),
      priority: editForm.priority
    }
    if (isAdmin.value) {
      payload.status = editForm.status
    }
    await updateRequest(requestId, payload)
    isEditMode.value = false
    await fetchDetail()
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to update request')
  } finally {
    isSaving.value = false
  }
}

// Delete request
const handleDelete = async () => {
  if (!confirm('Are you sure you want to delete this maintenance request permanently?')) {
    return
  }
  try {
    await deleteRequest(requestId)
    router.push('/requests')
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to delete request')
  }
}

// Submit Review (Approve/Reject)
const submitReview = async () => {
  reviewModal.isSubmitting = true
  try {
    await reviewRequest(requestId, {
      status: reviewModal.action,
      note: reviewModal.note
    })
    reviewModal.isOpen = false
    await fetchDetail()
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to review request')
  } finally {
    reviewModal.isSubmitting = false
  }
}

const formatDate = (dateStr?: string | null) => {
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

onMounted(() => {
  fetchDetail()
})
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <!-- Back Button -->
    <div class="mb-6 flex items-center justify-between">
      <NuxtLink
        to="/requests"
        class="inline-flex items-center gap-1.5 text-sm font-medium text-slate-600 hover:text-slate-900 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        Back to Requests List
      </NuxtLink>

      <div class="flex items-center gap-2" v-if="request && !isLoading">
        <!-- Edit Button Toggle -->
        <button
          v-if="canUserEdit && !isEditMode"
          @click="isEditMode = true"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-blue-700 bg-blue-50 hover:bg-blue-100 rounded-lg border border-blue-200 transition-colors"
        >
          <Edit3 class="w-4 h-4" />
          Edit Request
        </button>

        <!-- Delete Button (Admin Only) -->
        <button
          v-if="isAdmin"
          @click="handleDelete"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-rose-700 bg-rose-50 hover:bg-rose-100 rounded-lg border border-rose-200 transition-colors"
        >
          <Trash2 class="w-4 h-4" />
          Delete
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="bg-white rounded-2xl shadow-sm border border-slate-200 p-12 text-center">
      <Loader2 class="w-8 h-8 animate-spin mx-auto mb-2 text-factory-600" />
      <p class="text-sm text-slate-500">Loading request details...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="errorMessage" class="bg-white rounded-2xl shadow-sm border border-slate-200 p-8 text-center text-rose-600">
      <AlertCircle class="w-8 h-8 mx-auto mb-2 text-rose-500" />
      <p class="text-sm font-semibold">{{ errorMessage }}</p>
      <NuxtLink to="/requests" class="mt-4 inline-block text-xs font-semibold text-factory-600 underline">
        Return to requests list
      </NuxtLink>
    </div>

    <!-- Request Content -->
    <div v-else-if="request" class="space-y-6">
      <!-- Main Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
        <!-- Card Header -->
        <div class="p-6 sm:p-8 border-b border-slate-200 bg-slate-50/50 flex flex-wrap items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-3">
              <h1 class="text-2xl font-extrabold text-slate-900 font-mono tracking-tight">
                {{ request.asset_id }}
              </h1>
              <RequestsStatusBadge :status="request.status" />
              <RequestsPriorityBadge :priority="request.priority" />
            </div>
            <p class="mt-1 text-xs text-slate-500">
              Request ID: <span class="font-mono">{{ request.id }}</span>
            </p>
          </div>

          <!-- Supervisor/Admin Review Buttons (if Submitted) -->
          <div v-if="canReview && request.status === 'Submitted'" class="flex items-center gap-2">
            <button
              @click="reviewModal.action = 'Approved'; reviewModal.note = ''; reviewModal.isOpen = true"
              class="inline-flex items-center gap-1.5 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-sm font-medium rounded-lg shadow-sm transition-colors"
            >
              <CheckCircle class="w-4 h-4" />
              Approve
            </button>
            <button
              @click="reviewModal.action = 'Rejected'; reviewModal.note = ''; reviewModal.isOpen = true"
              class="inline-flex items-center gap-1.5 px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white text-sm font-medium rounded-lg shadow-sm transition-colors"
            >
              <XCircle class="w-4 h-4" />
              Reject
            </button>
          </div>
        </div>

        <!-- Normal View vs Edit Mode -->
        <div class="p-6 sm:p-8">
          <!-- Edit Form Mode -->
          <div v-if="isEditMode" class="space-y-5">
            <h2 class="text-lg font-bold text-slate-900">Editing Maintenance Request</h2>
            <div>
              <label class="block text-sm font-semibold text-slate-700 mb-1">Asset ID</label>
              <input
                v-model="editForm.asset_id"
                type="text"
                class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-factory-500 font-mono uppercase"
              />
            </div>

            <div>
              <label class="block text-sm font-semibold text-slate-700 mb-1">Priority</label>
              <select
                v-model="editForm.priority"
                class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-factory-500"
              >
                <option value="Low">Low</option>
                <option value="Medium">Medium</option>
                <option value="High">High</option>
                <option value="Urgent">Urgent</option>
              </select>
            </div>

            <div v-if="isAdmin">
              <label class="block text-sm font-semibold text-slate-700 mb-1">Status (Admin Override)</label>
              <select
                v-model="editForm.status"
                class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-factory-500"
              >
                <option value="Submitted">Submitted</option>
                <option value="Approved">Approved</option>
                <option value="Rejected">Rejected</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-semibold text-slate-700 mb-1">Problem Description</label>
              <textarea
                v-model="editForm.problem_description"
                rows="5"
                class="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-factory-500"
              ></textarea>
            </div>

            <div class="flex items-center gap-3 pt-3">
              <button
                @click="handleSaveEdit"
                :disabled="isSaving"
                class="inline-flex items-center gap-2 px-4 py-2 bg-factory-600 hover:bg-factory-700 text-white text-sm font-medium rounded-lg shadow-sm disabled:opacity-50"
              >
                <Loader2 v-if="isSaving" class="w-4 h-4 animate-spin" />
                <Save v-else class="w-4 h-4" />
                Save Changes
              </button>
              <button
                @click="isEditMode = false"
                class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg"
              >
                Cancel
              </button>
            </div>
          </div>

          <!-- View Mode -->
          <div v-else class="space-y-6">
            <!-- Metadata Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 rounded-xl bg-slate-50 border border-slate-200 text-sm">
              <div class="flex items-start gap-3">
                <UserIcon class="w-4 h-4 text-slate-400 mt-0.5" />
                <div>
                  <span class="text-xs text-slate-500 block uppercase font-medium">Created By</span>
                  <span class="font-semibold text-slate-800">{{ request.creator_name || 'Operator' }}</span>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <Calendar class="w-4 h-4 text-slate-400 mt-0.5" />
                <div>
                  <span class="text-xs text-slate-500 block uppercase font-medium">Created At</span>
                  <span class="font-semibold text-slate-800">{{ formatDate(request.created_at) }}</span>
                </div>
              </div>

              <div v-if="request.reviewed_by" class="flex items-start gap-3">
                <UserIcon class="w-4 h-4 text-slate-400 mt-0.5" />
                <div>
                  <span class="text-xs text-slate-500 block uppercase font-medium">Last Reviewed By</span>
                  <span class="font-semibold text-slate-800">{{ request.reviewer_name || 'Supervisor' }}</span>
                </div>
              </div>

              <div v-if="request.reviewed_at" class="flex items-start gap-3">
                <Calendar class="w-4 h-4 text-slate-400 mt-0.5" />
                <div>
                  <span class="text-xs text-slate-500 block uppercase font-medium">Reviewed At</span>
                  <span class="font-semibold text-slate-800">{{ formatDate(request.reviewed_at) }}</span>
                </div>
              </div>
            </div>

            <!-- Problem Description -->
            <div>
              <h3 class="text-sm font-bold text-slate-900 uppercase tracking-wider mb-2">
                Problem Description
              </h3>
              <div class="p-4 rounded-xl bg-white border border-slate-200 text-slate-800 text-sm leading-relaxed whitespace-pre-wrap">
                {{ request.problem_description }}
              </div>
            </div>

            <!-- Reviewer Note (if present) -->
            <div v-if="request.review_note" class="p-4 rounded-xl bg-blue-50/50 border border-blue-200 text-sm">
              <h3 class="text-xs font-bold text-blue-900 uppercase tracking-wider mb-1">
                Reviewer Note
              </h3>
              <p class="text-blue-800 italic">
                "{{ request.review_note }}"
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Audit Trail Timeline Section (Bonus Task!) -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 sm:p-8">
        <h2 class="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
          Audit Trail & Status History
        </h2>
        <RequestsAuditTimeline :logs="request.audit_logs" />
      </div>
    </div>

    <!-- Review Confirmation Modal -->
    <div
      v-if="reviewModal.isOpen"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-xl border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900 mb-2">
          {{ reviewModal.action === 'Approved' ? 'Approve' : 'Reject' }} Maintenance Request
        </h3>
        <p class="text-sm text-slate-600 mb-4">
          Machine / Asset: <strong class="text-slate-800">{{ request?.asset_id }}</strong>
        </p>

        <div class="mb-4">
          <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
            Reviewer Note / Justification (Optional)
          </label>
          <textarea
            v-model="reviewModal.note"
            rows="3"
            placeholder="e.g. Disetujui. Teknisi dijadwalkan besok..."
            class="w-full text-sm border border-slate-300 rounded-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-factory-500"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2">
          <button
            type="button"
            @click="reviewModal.isOpen = false"
            class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg"
          >
            Cancel
          </button>
          <button
            type="button"
            :disabled="reviewModal.isSubmitting"
            @click="submitReview"
            class="px-4 py-2 text-sm font-medium text-white rounded-lg shadow-sm"
            :class="reviewModal.action === 'Approved' ? 'bg-emerald-600 hover:bg-emerald-700' : 'bg-rose-600 hover:bg-rose-700'"
          >
            Confirm {{ reviewModal.action }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

