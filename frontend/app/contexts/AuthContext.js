'use client'

import { createContext, useContext, useState, useEffect } from 'react'

const AuthContext = createContext()

export function AuthProvider({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('token')
    setIsAuthenticated(!!token)
    setLoading(false)
  }, [])

  const login = async (email, password) => {
    try {
      const response = await fetch('http://localhost/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      if (!response.ok) {
        throw new Error('Ошибка входа')
      }

      const data = await response.json()
      localStorage.setItem('token', data.token)
      setIsAuthenticated(true)
      return data
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const register = async (formData) => {
    try {
      const response = await fetch('http://localhost/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      const data = await response.json()

      if (!response.ok) {
        // Handle validation errors from the backend
        if (data.errors) {
          const errorMessage = Object.values(data.errors).join(', ')
          throw new Error(errorMessage)
        }
        throw new Error(data.message || 'Ошибка регистрации')
      }

      localStorage.setItem('token', data.token)
      setIsAuthenticated(true)
      return { success: true, data }
    } catch (error) {
      console.error('Registration failed:', error)
      return { success: false, error: error.message }
    }
  }

  const logout = () => {
    localStorage.removeItem('token')
    setIsAuthenticated(false)
  }

  if (loading) {
    return null
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
} 