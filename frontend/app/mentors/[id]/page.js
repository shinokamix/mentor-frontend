'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import { useAuth } from '../../contexts/AuthContext'

export default function MentorProfile() {
  const { id } = useParams()
  const { isAuthenticated } = useAuth()
  const [mentor, setMentor] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [rating, setRating] = useState(0)
  const [review, setReview] = useState('')

  useEffect(() => {
    const fetchMentor = async () => {
      try {
        setLoading(true)
        const token = localStorage.getItem('token')
        const headers = {
          'Content-Type': 'application/json',
        }
        if (token) {
          headers['Authorization'] = `Bearer ${token}`
        }

        const response = await fetch(`http://localhost/api/mentors/get/${id}`, {
          method: 'GET',
          headers,
        })

        if (!response.ok) {
          throw new Error('Не удалось загрузить профиль ментора')
        }

        const data = await response.json()
        setMentor(data)
        setError(null)
      } catch (err) {
        console.error('Failed to fetch mentor:', err)
        setError('Не удалось загрузить профиль ментора')
      } finally {
        setLoading(false)
      }
    }

    fetchMentor()
  }, [id])

  const handleSubmitReview = async (e) => {
    e.preventDefault()
    if (!isAuthenticated) {
      alert('Пожалуйста, войдите в систему, чтобы оставить отзыв')
      return
    }

    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost/api/reviews/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          mentor_email: id,
          rating,
          review,
        }),
      })

      if (!response.ok) {
        throw new Error('Не удалось отправить отзыв')
      }

      // Обновляем данные ментора после успешного отзыва
      const updatedMentor = await fetch(`http://localhost/api/mentors/get/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }).then(res => res.json())

      setMentor(updatedMentor)
      setRating(0)
      setReview('')
      alert('Отзыв успешно отправлен')
    } catch (err) {
      console.error('Failed to submit review:', err)
      alert('Не удалось отправить отзыв')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Загрузка профиля...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600">{error}</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!mentor) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <p className="text-yellow-800">Ментор не найден</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 text-white py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-4xl font-bold text-center mb-4">
            Профиль ментора
          </h1>
          <p className="text-xl text-center text-blue-100 max-w-3xl mx-auto">
            Узнайте больше о менторе и оставьте свой отзыв
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <div className="p-8">
            <div className="space-y-8">
              {/* Mentor Info */}
              <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Информация о менторе</h2>
                <div className="bg-gray-50 rounded-lg p-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <p className="text-sm font-medium text-gray-500">Email</p>
                      <p className="mt-1 text-lg text-gray-900">{mentor.mentor_email}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-500">Контактная информация</p>
                      <p className="mt-1 text-lg text-gray-900">{mentor.contact}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-500">Средний рейтинг</p>
                      <p className="mt-1 text-lg text-gray-900">{mentor.average_rating?.toFixed(1) || '0.0'}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-500">Количество отзывов</p>
                      <p className="mt-1 text-lg text-gray-900">{mentor.reviews_count || 0}</p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-500">Цена</p>
                      <p className="mt-1 text-lg text-gray-900">
                        {mentor.price ? `${mentor.price} ₽/час` : 'Цена договорная'}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Review Form */}
              {isAuthenticated && (
                <div>
                  <h2 className="text-2xl font-bold text-gray-900 mb-6">Оставить отзыв</h2>
                  <form onSubmit={handleSubmitReview} className="space-y-6">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Оценка
                      </label>
                      <div className="flex space-x-2">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <button
                            key={star}
                            type="button"
                            onClick={() => setRating(star)}
                            className={`p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${
                              rating >= star ? 'text-yellow-400' : 'text-gray-300'
                            }`}
                          >
                            <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                            </svg>
                          </button>
                        ))}
                      </div>
                    </div>
                    <div>
                      <label htmlFor="review" className="block text-sm font-medium text-gray-700 mb-2">
                        Отзыв
                      </label>
                      <textarea
                        id="review"
                        rows={4}
                        value={review}
                        onChange={(e) => setReview(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        placeholder="Напишите ваш отзыв..."
                      />
                    </div>
                    <button
                      type="submit"
                      className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200"
                    >
                      Отправить отзыв
                    </button>
                  </form>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 