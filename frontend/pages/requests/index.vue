<script setup lang="ts">
import { 
  Plus, Search, Filter, RotateCcw, Eye, Edit3, Trash2, 
  CheckCircle, XCircle, AlertTriangle, Loader2 
} from 'lucide-vue-next'
import type { MaintenanceRequest } from '~/types'

definePageMeta({
  middleware: 'auth'
})

const { user, isOperator, canReview, isAdmin } = useAuth()
const { getRequests, deleteRequest, reviewRequest } = useRequests()

const requests = ref<MaintenanceRequest[]>([])
const isLoading = ref(true)
const errorMessage = ref('')

// Filters
const filters = reactive({
  status: 'All',
  priority: 'All',
  search: ''
})

// Review Modal State
const reviewModal = reactive({
  isOpen: false,
  requestId: '',
  assetId: '',
  action: 'Approved' as 'Approved' | 'Rejected',
  note: '',
  isSubmitting: false
})

const fetchRequests = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    requests.value = await getRequests({
      status: filters.status,
      priority: filters.priority,
      search: filters.search
    })
  } catch (err: any) {
    errorMessage.value = err?.data?.error || 'Failed to fetch maintenance requests'
  } finally {
    isLoading.value = false
  }
}

// Watch filters with debounce for search
watch([() => filters.status, () => filters.priority], () => {
  fetchRequests()
})

let searchTimeout: any = null
const onSearchInput = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchRequests()
  }, 350)
}

const resetFilters = () => {
  filters.status = 'All'
  filters.priority = 'All'
  filters.search = ''
  fetchRequests()
}

// Check if current user can edit this specific request
const canEdit = (req: MaintenanceRequest) => {
  if (isAdmin.value) return true
  // Operator / Supervisor can only edit their own request while status is 'Submitted'
  return req.created_by === user.value?.id && req.status === 'Submitted'
}

// Delete request
const handleDelete = async (req: MaintenanceRequest) => {
  if (!confirm(`Are you sure you want to delete maintenance request for asset ${req.asset_id}?`)) {
    return
  }
  try {
    await deleteRequest(req.id)
    await fetchRequests()
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to delete request')
  }
}

// Open Review Modal
const openReviewModal = (req: MaintenanceRequest, action: 'Approved' | 'Rejected') => {
  reviewModal.requestId = req.id
  reviewModal.assetId = req.asset_id
  reviewModal.action = action
  reviewModal.note = ''
  reviewModal.isOpen = true
}

// Submit Review
const submitReview = async () => {
  reviewModal.isSubmitting = true
  try {
    await reviewRequest(reviewModal.requestId, {
      status: reviewModal.action,
      note: reviewModal.note
    })
    reviewModal.isOpen = false
    await fetchRequests()
  } catch (err: any) {
    alert(err?.data?.error || 'Failed to submit review')
  } finally {
    reviewModal.isSubmitting = false
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

onMounted(() => {
  fetchRequests()
})
</script>

<template>
  <div>
    <!-- Header Section -->
    <div class="sm:flex sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Maintenance Requests</h1>
        <p class="mt-1 text-sm text-slate-500">
          <span v-if="isOperator" class="font-medium text-amber-700 bg-amber-50 px-2 py-0.5 rounded">
            Showing requests raised by you
          </span>
          <span v-else class="text-slate-600">
            Factory-wide maintenance requests overview
          </span>
        </p>
      </div>

      <div class="mt-4 sm:mt-0">
        <NuxtLink
          to="/requests/create"
          class="inline-flex items-center gap-2 px-4 py-2 bg-factory-600 hover:bg-factory-700 text-white text-sm font-medium rounded-lg shadow-sm transition-colors"
        >
          <Plus class="w-4 h-4" />
          Raise New Request
        </NuxtLink>
      </div>
    </div>

    <!-- Filters Bar -->
    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 mb-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
        <!-- Search -->
        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Search</label>
          <div class="relative">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-3" />
            <input
              v-model="filters.search"
              @input="onSearchInput"
              type="text"
              placeholder="Asset ID or Problem..."
              class="w-full pl-9 pr-3 py-2 text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-factory-500"
            />
          </div>
        </div>

        <!-- Filter Status -->
        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Status</label>
          <select
            v-model="filters.status"
            class="w-full py-2 px-3 text-sm border border-slate-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-factory-500"
          >
            <option value="All">All Statuses</option>
            <option value="Submitted">Submitted</option>
            <option value="Approved">Approved</option>
            <option value="Rejected">Rejected</option>
          </select>
        </div>

        <!-- Filter Priority -->
        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Priority</label>
          <select
            v-model="filters.priority"
            class="w-full py-2 px-3 text-sm border border-slate-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-factory-500"
          >
            <option value="All">All Priorities</option>
            <option value="Urgent">Urgent</option>
            <option value="High">High</option>
            <option value="Medium">Medium</option>
            <option value="Low">Low</option>
          </select>
        </div>

        <!-- Reset Button -->
        <div class="flex items-end">
          <button
            @click="resetFilters"
            class="w-full inline-flex items-center justify-center gap-1.5 py-2 px-3 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            Reset Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Table Content -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <!-- Loading State -->
      <div v-if="isLoading" class="py-16 text-center text-slate-500">
        <Loader2 class="w-8 h-8 animate-spin mx-auto mb-2 text-factory-600" />
        <p class="text-sm">Loading requests...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage" class="p-8 text-center text-rose-600">
        <AlertTriangle class="w-8 h-8 mx-auto mb-2 text-rose-500" />
        <p class="text-sm font-medium">{{ errorMessage }}</p>
        <button @click="fetchRequests" class="mt-3 text-xs text-factory-600 font-semibold underline">Retry</button>
      </div>

      <!-- Empty State -->
      <div v-else-if="requests.length === 0" class="py-16 text-center">
        <div class="w-12 h-12 rounded-full bg-slate-100 text-slate-400 flex items-center justify-center mx-auto mb-3">
          <Filter class="w-6 h-6" />
        </div>
        <h3 class="text-base font-semibold text-slate-800">No maintenance requests found</h3>
        <p class="text-sm text-slate-500 max-w-sm mx-auto mt-1">
          No requests match your current filters. Try changing filter criteria or create a new request.
        </p>
        <div class="mt-4">
          <button @click="resetFilters" class="text-sm font-medium text-factory-600 hover:text-factory-700">
            Reset all filters
          </button>
        </div>
      </div>

      <!-- Requests Table -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-200 text-left text-sm">
          <thead class="bg-slate-50 text-slate-600 font-semibold uppercase text-[11px] tracking-wider">
            <tr>
              <th scope="col" class="py-3.5 pl-4 pr-3 sm:pl-6">Asset ID</th>
              <th scope="col" class="px-3 py-3.5">Problem Description</th>
              <th scope="col" class="px-3 py-3.5">Priority</th>
              <th scope="col" class="px-3 py-3.5">Status</th>
              <th scope="col" class="px-3 py-3.5">Created By</th>
              <th scope="col" class="px-3 py-3.5">Date</th>
              <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr v-for="req in requests" :key="req.id" class="hover:bg-slate-50/80 transition-colors">
              <!-- Asset ID -->
              <td class="whitespace-nowrap py-4 pl-4 pr-3 sm:pl-6 font-bold text-slate-900">
                <NuxtLink :to="`/requests/${req.id}`" class="hover:text-factory-600">
                  {{ req.asset_id }}
                </NuxtLink>
              </td>

              <!-- Problem Description -->
              <td class="px-3 py-4 text-slate-600 max-w-xs truncate">
                {{ req.problem_description }}
              </td>

              <!-- Priority -->
              <td class="whitespace-nowrap px-3 py-4">
                <RequestsPriorityBadge :priority="req.priority" />
              </td>

              <!-- Status -->
              <td class="whitespace-nowrap px-3 py-4">
                <RequestsStatusBadge :status="req.status" />
              </td>

              <!-- Created By -->
              <td class="whitespace-nowrap px-3 py-4 text-slate-600">
                {{ req.creator_name || 'User' }}
              </td>

              <!-- Date -->
              <td class="whitespace-nowrap px-3 py-4 text-slate-500 text-xs">
                {{ formatDate(req.created_at) }}
              </td>

              <!-- Actions -->
              <td class="whitespace-nowrap py-4 pl-3 pr-4 sm:pr-6 text-right space-x-1">
                <!-- View Detail -->
                <NuxtLink
                  :to="`/requests/${req.id}`"
                  class="inline-flex items-center p-1.5 text-slate-500 hover:text-slate-900 rounded hover:bg-slate-100"
                  title="View Detail"
                >
                  <Eye class="w-4 h-4" />
                </NuxtLink>

                <!-- Edit (Subject to RBAC) -->
                <NuxtLink
                  v-if="canEdit(req)"
                  :to="`/requests/${req.id}?edit=true`"
                  class="inline-flex items-center p-1.5 text-blue-600 hover:text-blue-800 rounded hover:bg-blue-50"
                  title="Edit Request"
                >
                  <Edit3 class="w-4 h-4" />
                </NuxtLink>

                <!-- Quick Review Buttons (Supervisor & Admin only for Submitted requests) -->
                <template v-if="canReview && req.status === 'Submitted'">
                  <button
                    @click="openReviewModal(req, 'Approved')"
                    class="inline-flex items-center p-1.5 text-emerald-600 hover:text-emerald-800 rounded hover:bg-emerald-50"
                    title="Approve Request"
                  >
                    <CheckCircle class="w-4 h-4" />
                  </button>
                  <button
                    @click="openReviewModal(req, 'Rejected')"
                    class="inline-flex items-center p-1.5 text-rose-600 hover:text-rose-800 rounded hover:bg-rose-50"
                    title="Reject Request"
                  >
                    <XCircle class="w-4 h-4" />
                  </button>
                </template>

                <!-- Delete (Admin Only) -->
                <button
                  v-if="isAdmin"
                  @click="handleDelete(req)"
                  class="inline-flex items-center p-1.5 text-rose-500 hover:text-rose-700 rounded hover:bg-rose-50"
                  title="Delete Request"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
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
          Machine / Asset: <strong class="text-slate-800">{{ reviewModal.assetId }}</strong>
        </p>

        <div class="mb-4">
          <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
            Reviewer Note / Justification (Optional)
          </label>
          <textarea
            v-model="reviewModal.note"
            rows="3"
            placeholder="e.g. Approved. Technician scheduled for tomorrow morning..."
            class="w-full text-sm border border-slate-300 rounded-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-factory-500"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2">
          <button
            type="button"
            @click="reviewModal.isOpen = false"
            class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
          >
            Cancel
          </button>
          <button
            type="button"
            :disabled="reviewModal.isSubmitting"
            @click="submitReview"
            class="px-4 py-2 text-sm font-medium text-white rounded-lg shadow-sm transition-colors"
            :class="reviewModal.action === 'Approved' ? 'bg-emerald-600 hover:bg-emerald-700' : 'bg-rose-600 hover:bg-rose-700'"
          >
            Confirm {{ reviewModal.action }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

