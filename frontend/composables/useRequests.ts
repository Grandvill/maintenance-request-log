import type { 
  MaintenanceRequest, 
  RequestDetail, 
  CreateRequestPayload, 
  UpdateRequestPayload, 
  ReviewRequestPayload 
} from '~/types'

export const useRequests = () => {
  const { apiFetch } = useApi()

  const getRequests = async (filters: { status?: string; priority?: string; search?: string } = {}) => {
    const params = new URLSearchParams()
    if (filters.status && filters.status !== 'All') params.append('status', filters.status)
    if (filters.priority && filters.priority !== 'All') params.append('priority', filters.priority)
    if (filters.search) params.append('search', filters.search)

    const queryStr = params.toString() ? `?${params.toString()}` : ''
    return await apiFetch<MaintenanceRequest[]>(`/requests${queryStr}`)
  }

  const getRequest = async (id: string) => {
    return await apiFetch<RequestDetail>(`/requests/${id}`)
  }

  const createRequest = async (payload: CreateRequestPayload) => {
    return await apiFetch<MaintenanceRequest>('/requests', {
      method: 'POST',
      body: payload
    })
  }

  const updateRequest = async (id: string, payload: UpdateRequestPayload) => {
    return await apiFetch<MaintenanceRequest>(`/requests/${id}`, {
      method: 'PUT',
      body: payload
    })
  }

  const reviewRequest = async (id: string, payload: ReviewRequestPayload) => {
    return await apiFetch<MaintenanceRequest>(`/requests/${id}/review`, {
      method: 'PATCH',
      body: payload
    })
  }

  const deleteRequest = async (id: string) => {
    return await apiFetch<{ message: string }>(`/requests/${id}`, {
      method: 'DELETE'
    })
  }

  return {
    getRequests,
    getRequest,
    createRequest,
    updateRequest,
    reviewRequest,
    deleteRequest
  }
}

