import { useState, useEffect } from 'react'

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const checkAuth = async () => {
      const accessToken = localStorage.getItem('access_token')

      // 1. Если access_token есть — пробуем запрос
      if (accessToken) {
        const success = await verifyAccessToken(accessToken)

        if (success) {
          setIsAuthenticated(true)
          setLoading(false)
          return
        }
      }

      // 2. Пробуем обновить токен через refresh_token
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        const newAccessToken = await refreshAccessToken(refreshToken)

        if (newAccessToken) {
          localStorage.setItem('access_token', newAccessToken)
          setIsAuthenticated(true)
        } else {
          // refresh не сработал — удаляем токены
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
        }
      }

      setLoading(false)
    }

    checkAuth()
  }, [])

  return { isAuthenticated, loading }
}

// Проверка access_token на работоспособность
async function verifyAccessToken(token) {
  try {
    const res = await fetch('http://localhost/api/protected', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    return res.ok
  } catch (err) {
    return false
  }
}

// Обновление access_token
async function refreshAccessToken(refreshToken) {
  try {
    const res = await fetch('http://localhost/api/auth/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ refresh_token: refreshToken })
    })

    if (!res.ok) return null

    const data = await res.json()
    return data.access_token
  } catch (err) {
    return null
  }
}
