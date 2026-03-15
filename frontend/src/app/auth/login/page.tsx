"use client"

import type React from "react"

import { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card"
import { Input } from "~/components/ui/input"
import { Label } from "~/components/ui/label"

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    try {
      // TODO: API呼び出しでログイン処理を実装
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      })

      if (!response.ok) {
        const data = await response.json()
        setError(data.message || "ログインに失敗しました")
        return
      }

      // ログイン成功時
      router.push("/profile")
    } catch (err) {
      setError("エラーが発生しました。もう一度お試しください。")
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50 flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <Card className="shadow-lg border-pink-100">
          <CardHeader className="text-center bg-pink-50 rounded-t-lg">
            <CardTitle className="text-2xl text-pink-700">ログイン</CardTitle>
            <CardDescription>パパクイズにログインする</CardDescription>
          </CardHeader>
          <CardContent className="pt-6">
            <form onSubmit={handleLogin} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email" className="text-gray-700">
                  メールアドレス
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="example@email.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="border-gray-300"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-gray-700">
                  パスワード
                </Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="パスワード"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="border-gray-300"
                />
              </div>

              {error && (
                <div className="p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">{error}</div>
              )}

              <Button type="submit" disabled={loading} className="w-full bg-pink-600 hover:bg-pink-700">
                {loading ? "ログイン中..." : "ログイン"}
              </Button>
            </form>

            <div className="mt-6 pt-6 border-t border-gray-200">
              <p className="text-center text-gray-600 text-sm">
                アカウントをお持ちでないですか？{" "}
                <Link href="/auth/signup" className="text-pink-600 hover:text-pink-700 font-semibold">
                  登録する
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>

        <div className="mt-4 text-center">
          <Link href="/" className="text-pink-600 hover:text-pink-700 text-sm font-semibold">
            ← ホームに戻る
          </Link>
        </div>
      </div>
    </div>
  )
}
