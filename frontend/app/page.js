'use client'

import Image from "next/image";
import { useAuth } from "@/hooks/auth";

export default function Home() {

  const { isAuthenticated, loading } = useAuth();




  return (
    <div className="flex h-screen">
      <div className="my-auto ml-10">
        

        <h1 className="text-9xl">MentorLink</h1>
        {loading ? (
          <p className="text-gray-500 text-xl">Проверка авторизации...</p>
        ) : isAuthenticated ? (
          <p className="text-green-600 text-xl">Вы авторизованы</p>
        ) : (
          <p className="text-red-600 text-xl">Вы не авторизованы</p>
        )}
      </div>
    </div>
  );
}
