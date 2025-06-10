'use client'

import { useAuth } from '../contexts/AuthContext'

export function useApi() {
  const { user } = useAuth()

  const fetchWithAuth = async (url, options = {}) => {
    const token = localStorage.getItem('access_token')
    
    const headers = {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (response.status === 401) {
      // Токен истек или недействителен
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('role')
      window.location.href = '/login'
      throw new Error('Unauthorized')
    }

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || 'Something went wrong')
    }

    return data
  }

  return {
    get: (url, options = {}) => fetchWithAuth(url, { ...options, method: 'GET' }),
    post: (url, data, options = {}) => 
      fetchWithAuth(url, { 
        ...options, 
        method: 'POST', 
        body: JSON.stringify(data) 
      }),
    put: (url, data, options = {}) => 
      fetchWithAuth(url, { 
        ...options, 
        method: 'PUT', 
        body: JSON.stringify(data) 
      }),
    delete: (url, options = {}) => 
      fetchWithAuth(url, { ...options, method: 'DELETE' }),
  }
} 