'use client'

import { createContext, useContext } from 'react'

const ApiContext = createContext()

export function ApiProvider({ children }) {
  const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api'

  const api = {
    async get(endpoint) {
      try {
        const token = localStorage.getItem('token')
        const headers = new Headers({
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        })
        
        if (token) {
          headers.append('Authorization', `Bearer ${token}`)
        }

        const response = await fetch(`${baseURL}${endpoint}`, {
          method: 'GET',
          headers,
          credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
            throw new Error('Unauthorized')
          }
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        return data
      } catch (error) {
        console.error('API request failed:', error)
        throw error
      }
    },

    async post(endpoint, body) {
      try {
        const token = localStorage.getItem('token')
        const headers = new Headers({
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        })
        
        if (token) {
          headers.append('Authorization', `Bearer ${token}`)
        }

        const response = await fetch(`${baseURL}${endpoint}`, {
          method: 'POST',
          headers,
          body: JSON.stringify(body),
          credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
            throw new Error('Unauthorized')
          }
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        return data
      } catch (error) {
        console.error('API request failed:', error)
        throw error
      }
    },

    async put(endpoint, body) {
      try {
        const token = localStorage.getItem('token')
        const headers = new Headers({
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        })
        
        if (token) {
          headers.append('Authorization', `Bearer ${token}`)
        }

        const response = await fetch(`${baseURL}${endpoint}`, {
          method: 'PUT',
          headers,
          body: JSON.stringify(body),
          credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
            throw new Error('Unauthorized')
          }
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        return data
      } catch (error) {
        console.error('API request failed:', error)
        throw error
      }
    },

    async delete(endpoint) {
      try {
        const token = localStorage.getItem('token')
        const headers = new Headers({
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        })
        
        if (token) {
          headers.append('Authorization', `Bearer ${token}`)
        }

        const response = await fetch(`${baseURL}${endpoint}`, {
          method: 'DELETE',
          headers,
          credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
            throw new Error('Unauthorized')
          }
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        return data
      } catch (error) {
        console.error('API request failed:', error)
        throw error
      }
    },
  }

  return (
    <ApiContext.Provider value={{ api }}>
      {children}
    </ApiContext.Provider>
  )
}

export function useApi() {
  const context = useContext(ApiContext)
  if (!context) {
    throw new Error('useApi must be used within an ApiProvider')
  }
  return context
} 