'use client'

import { useState } from 'react'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSubmit = async (e) => {
  e.preventDefault()
  setError('')
  setSuccess('')

  try {
    const res = await fetch('http://localhost/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, password })
    })

    const data = await res.json() // 💡 перемещено сюда

    if (!res.ok) {
      throw new Error(data?.message || 'Ошибка входа')
    }

    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)
    localStorage.setItem('role', data.role)

    setSuccess('Успешный вход!')
  } catch (err) {
    setError(err.message)
  }
}


  return (
    <div className="max-w-md mx-auto mt-10 p-4 border rounded">
      <h2 className="text-xl font-bold mb-4">Вход</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          className="w-full border p-2 rounded"
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={e => setPassword(e.target.value)}
          className="w-full border p-2 rounded"
        />
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
        >
          Войти
        </button>
      </form>

      {error && <p className="text-red-600 mt-4">{error}</p>}
      {success && <p className="text-green-600 mt-4">{success}</p>}
    </div>
  )
}
