"use client"

import type React from "react"

import { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card"
import { Input } from "~/components/ui/input"
import { Label } from "~/components/ui/label"

export default function SignUpPage() {
  const router = useRouter()
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  })
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    // バリデーション
    if (formData.password !== formData.confirmPassword) {
      setError("パスワードが一致しません")
      setLoading(false)
      return
    }

    if (formData.password.length < 6) {
      setError("パスワードは6文字以上である必要があります")
      setLoading(false)
      return
    }

    try {
      // TODO: API呼び出しで登録処理を実装
      const response = await fetch(`/api/auth/signup`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: formData.name,
          email: formData.email,
          password: formData.password,
        }),
      })

      if (!response.ok) {
        const data = await response.json()
        setError(data.message || "登録に失敗しました")
        return
      }

      router.push("/auth/login")
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
            <CardTitle className="text-2xl text-pink-700">新規登録</CardTitle>
            <CardDescription>パパクイズにアカウントを作成する</CardDescription>
          </CardHeader>
          <CardContent className="pt-6">
            <form onSubmit={handleSignUp} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name" className="text-gray-700">
                  お名前
                </Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="田中太郎"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                  className="border-gray-300"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="email" className="text-gray-700">
                  メールアドレス
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="example@email.com"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
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
                  placeholder="6文字以上のパスワード"
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  required
                  className="border-gray-300"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirmPassword" className="text-gray-700">
                  パスワード（確認）
                </Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="パスワードを再度入力"
                  name="confirmPassword"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  required
                  className="border-gray-300"
                />
              </div>

              {error && (
                <div className="p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">{error}</div>
              )}

              <Button type="submit" disabled={loading} className="w-full bg-pink-600 hover:bg-pink-700">
                {loading ? "登録中..." : "登録する"}
              </Button>
            </form>

            <div className="mt-6 pt-6 border-t border-gray-200">
              <p className="text-center text-gray-600 text-sm">
                既にアカウントをお持ちですか？{" "}
                <Link href="/auth/login" className="text-pink-600 hover:text-pink-700 font-semibold">
                  ログイン
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
