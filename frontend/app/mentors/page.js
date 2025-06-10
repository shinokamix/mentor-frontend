'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useAuth } from '../contexts/AuthContext'

export default function MentorsList() {
  const { isAuthenticated } = useAuth()
  const [mentors, setMentors] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [showCopyNotification, setShowCopyNotification] = useState(false)
  const [copiedContact, setCopiedContact] = useState('')

  useEffect(() => {
    const fetchMentors = async () => {
      try {
        setLoading(true)
        const token = localStorage.getItem('token')
        const headers = {
          'Content-Type': 'application/json',
        }
        if (token) {
          headers['Authorization'] = `Bearer ${token}`
        }

        const response = await fetch('http://localhost/api/mentors/get', {
          method: 'GET',
          headers,
        })

        if (!response.ok) {
          const errorData = await response.json()
          throw new Error(errorData.message || 'Не удалось загрузить список менторов')
        }

        const data = await response.json()
        setMentors(data.mentors || [])
        setError(null)
      } catch (err) {
        console.error('Failed to fetch mentors:', err)
        setError(err.message || 'Не удалось загрузить список менторов')
      } finally {
        setLoading(false)
      }
    }

    fetchMentors()
  }, [])

  const handleCopyContact = async (contact) => {
    try {
      await navigator.clipboard.writeText(contact)
      setCopiedContact(contact)
      setShowCopyNotification(true)
      setTimeout(() => {
        setShowCopyNotification(false)
        setCopiedContact('')
      }, 3000)
    } catch (err) {
      console.error('Failed to copy contact:', err)
      alert('Не удалось скопировать контактные данные')
    }
  }

  const filteredMentors = mentors.filter((mentor) =>
    mentor.mentor_email.toLowerCase().includes(searchQuery.toLowerCase()) ||
    mentor.contact?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Загрузка списка менторов...</p>
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

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 text-white py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-4xl font-bold text-center mb-4">
            Наши менторы
          </h1>
          <p className="text-xl text-center text-blue-100 max-w-3xl mx-auto">
            Найдите подходящего ментора для вашего развития
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Search and Actions */}
        <div className="mb-8 flex flex-col sm:flex-row gap-4 items-center justify-between">
          <div className="w-full sm:w-96">
            <input
              type="text"
              placeholder="Поиск по email или контакту..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          {isAuthenticated && (
            <Link
              href="/become-mentor"
              className="w-full sm:w-auto px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200 text-center"
            >
              Стать ментором
            </Link>
          )}
        </div>

        {/* Mentors Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredMentors.map((mentor) => (
            <div
              key={mentor.mentor_email}
              className="bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow duration-200"
            >
              <div className="p-6">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  {mentor.mentor_email}
                </h3>
                <div className="flex items-center mb-4">
                  <div className="flex items-center">
                    <svg className="h-5 w-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                    <span className="ml-1 text-gray-900">{mentor.average_rating?.toFixed(1) || '0.0'}</span>
                  </div>
                  <span className="mx-2 text-gray-300">•</span>
                  <span className="text-gray-600">{mentor.reviews_count || 0} отзывов</span>
                </div>
                <p className="text-gray-900 mb-4 line-clamp-2">{mentor.contact}</p>
                <div className="flex justify-between items-center">
                  <span className="text-lg font-semibold text-gray-900">
                    {mentor.price ? `${mentor.price} ₽/час` : 'Цена договорная'}
                  </span>
                  <button
                    onClick={() => handleCopyContact(mentor.contact)}
                    className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200 flex items-center"
                  >
                    <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                    Связаться
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        {filteredMentors.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-600">Менторы не найдены</p>
          </div>
        )}
      </div>

      {/* Copy Notification */}
      {showCopyNotification && (
        <div className="fixed bottom-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg flex items-center">
          <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
          </svg>
          Контактные данные скопированы в буфер обмена
        </div>
      )}
    </div>
  )
} 