'use client'

import { useAuth } from './contexts/AuthContext'
import Link from 'next/link'
import Image from 'next/image'

export default function Home() {
  const { user } = useAuth()

  const features = [
    {
      title: 'Найдите идеального ментора',
      description: 'Выбирайте из проверенных экспертов в различных областях',
      icon: '🎯'
    },
    {
      title: 'Персональный подход',
      description: 'Индивидуальные программы обучения под ваши цели',
      icon: '👤'
    },
    {
      title: 'Отслеживайте прогресс',
      description: 'Система рейтингов и отзывов для оценки качества обучения',
      icon: '📈'
    }
  ]

  const benefits = [
    {
      title: 'Для учеников',
      items: [
        'Доступ к опытным наставникам',
        'Гибкий график обучения',
        'Разнообразие направлений',
        'Гарантия качества обучения'
      ]
    },
    {
      title: 'Для менторов',
      items: [
        'Возможность монетизировать знания',
        'Гибкий график работы',
        'Поддержка платформы',
        'Построение личного бренда'
      ]
    }
  ]

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-blue-600 to-blue-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-6">
              Найди своего идеального ментора
            </h1>
            <p className="text-xl md:text-2xl mb-8 text-blue-50">
              Развивайся с опытными наставниками в удобном формате
            </p>
            {!user && (
              <div className="space-x-4">
                <Link
                  href="/register"
                  className="inline-block bg-white text-blue-600 px-8 py-3 rounded-lg font-medium hover:bg-blue-50 transition-colors"
                >
                  Начать обучение
                </Link>
                <Link
                  href="/mentors"
                  className="inline-block border-2 border-white text-white px-8 py-3 rounded-lg font-medium hover:bg-white hover:text-blue-600 transition-colors"
                >
                  Найти ментора
                </Link>
              </div>
            )}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Почему выбирают нас
            </h2>
            <p className="text-xl text-gray-800">
              Мы создаем эффективную среду для обучения и развития
            </p>
          </div>
          <div className="grid md:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="text-4xl mb-4">{feature.icon}</div>
                <h3 className="text-xl font-semibold mb-2 text-gray-900">{feature.title}</h3>
                <p className="text-gray-800">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Преимущества платформы
            </h2>
            <p className="text-xl text-gray-800">
              Мы создаем возможности для всех участников
            </p>
          </div>
          <div className="grid md:grid-cols-2 gap-12">
            {benefits.map((benefit, index) => (
              <div key={index} className="bg-gray-50 p-8 rounded-xl shadow-sm">
                <h3 className="text-2xl font-semibold mb-6 text-blue-700">
                  {benefit.title}
                </h3>
                <ul className="space-y-4">
                  {benefit.items.map((item, itemIndex) => (
                    <li key={itemIndex} className="flex items-center text-gray-800">
                      <svg
                        className="h-5 w-5 text-green-600 mr-3"
                        fill="none"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path d="M5 13l4 4L19 7"></path>
                      </svg>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-blue-600 text-white py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold mb-4">
            Готовы начать свой путь к успеху?
          </h2>
          <p className="text-xl mb-8 text-blue-50">
            Присоединяйтесь к нашей платформе уже сегодня
          </p>
          {!user && (
            <div className="space-x-4">
              <Link
                href="/register"
                className="inline-block bg-white text-blue-600 px-8 py-3 rounded-lg font-medium hover:bg-blue-50 transition-colors"
              >
                Зарегистрироваться
              </Link>
              <Link
                href="/mentors"
                className="inline-block border-2 border-white text-white px-8 py-3 rounded-lg font-medium hover:bg-white hover:text-blue-600 transition-colors"
              >
                Найти ментора
              </Link>
            </div>
          )}
        </div>
      </section>
    </div>
  )
}
