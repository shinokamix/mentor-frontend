'use client'

import { useState } from 'react';
import Link from 'next/link';

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    repeat_password: '',
    role: 'user',
    contact: '',
  });

  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // Валидация формы
  const validateForm = () => {
    const { email, password, repeat_password, role } = formData;

    if (!email.includes('@')) {
      setError('Некорректный email');
      return false;
    }

    if (password.length < 6) {
      setError('Пароль должен содержать минимум 6 символов');
      return false;
    }

    if (password !== repeat_password) {
      setError('Пароли не совпадают');
      return false;
    }

    if (!['user', 'mentor', 'admin'].includes(role)) {
      setError('Выберите корректную роль');
      return false;
    }

    return true;
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
    if (error) setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      if (!validateForm()) return;

      const response = await fetch('http://localhost/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Ошибка регистрации');
      }

      // Можно добавить redirect на /login
    } catch (err) {
      setError(err.message || 'Произошла ошибка при регистрации');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md text-black">
        <h2 className="text-2xl font-bold text-center">Регистрация</h2>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium">Email</label>
              <input
                id="email"
                name="email"
                type="email"
                required
                value={formData.email}
                onChange={handleInputChange}
                disabled={isLoading}
                className="mt-1 block w-full border px-3 py-2 rounded-md"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium">Пароль</label>
              <input
                id="password"
                name="password"
                type="password"
                required
                value={formData.password}
                onChange={handleInputChange}
                disabled={isLoading}
                className="mt-1 block w-full border px-3 py-2 rounded-md"
              />
            </div>

            <div>
              <label htmlFor="repeat_password" className="block text-sm font-medium">Повторите пароль</label>
              <input
                id="repeat_password"
                name="repeat_password"
                type="password"
                required
                value={formData.repeat_password}
                onChange={handleInputChange}
                disabled={isLoading}
                className="mt-1 block w-full border px-3 py-2 rounded-md"
              />
            </div>

            <div>
              <label htmlFor="role" className="block text-sm font-medium">Роль</label>
              <select
                id="role"
                name="role"
                value={formData.role}
                onChange={handleInputChange}
                disabled={isLoading}
                className="mt-1 block w-full border px-3 py-2 rounded-md"
              >
                <option value="user">Пользователь</option>
                <option value="mentor">Наставник</option>
                <option value="admin">Админ</option>
              </select>
            </div>

            <div>
              <label htmlFor="contact" className="block text-sm font-medium">Контакт (опционально)</label>
              <input
                id="contact"
                name="contact"
                type="text"
                value={formData.contact}
                onChange={handleInputChange}
                disabled={isLoading}
                className="mt-1 block w-full border px-3 py-2 rounded-md"
              />
            </div>
          </div>

          {error && <div className="text-red-500 text-sm text-center">{error}</div>}

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>
        </form>

        <div className="text-sm text-center text-gray-600">
          Уже есть аккаунт?{' '}
          <Link href="/login" className="text-blue-600 hover:text-blue-500 font-medium">
            Войти
          </Link>
        </div>
      </div>
    </div>
  );
}
