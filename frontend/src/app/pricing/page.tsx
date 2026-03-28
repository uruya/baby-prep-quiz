"use client"

import { useState } from "react"
import Link from "next/link"
import { Home, Check, Lock } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card"

export default function PricingPage() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const handleCheckout = async () => {
    setLoading(true)
    setError("")
    try {
      const res = await fetch("/api/billing/checkout", { method: "POST" })
      if (!res.ok) {
        const data = await res.json().catch(() => ({}))
        if (res.status === 401) {
          window.location.href = "/auth/login"
          return
        }
        setError(data.message || "エラーが発生しました。もう一度お試しください。")
        return
      }
      const { url } = await res.json()
      if (url) window.location.href = url
    } catch {
      setError("エラーが発生しました。もう一度お試しください。")
    } finally {
      setLoading(false)
    }
  }

  const freeFeatures = [
    "妊娠の基礎知識クイズ",
    "出産の準備クイズ",
    "学習リソース閲覧",
    "クイズ結果の保存（ログイン時）",
  ]

  const premiumFeatures = [
    "無料プランのすべての機能",
    "赤ちゃんのお世話クイズ",
    "栄養と食事クイズ",
    "赤ちゃんの発達クイズ",
    "全カテゴリの達成状況トラッキング",
  ]

  return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50">
      <div className="container mx-auto px-4 py-8">
        <header className="flex items-center justify-between mb-8">
          <Link href="/">
            <Button variant="ghost" size="sm" className="flex items-center gap-1 shrink-0">
              <Home className="h-4 w-4" />
              ホーム
            </Button>
          </Link>
          <h1 className="text-xl md:text-2xl font-bold text-pink-600">料金プラン</h1>
          <div className="w-20 shrink-0" />
        </header>

        <div className="max-w-2xl mx-auto">
          <p className="text-center text-gray-600 mb-8">
            パパとしての準備を全力でサポートします
          </p>

          {error && (
            <div className="mb-6 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm text-center">
              {error}
            </div>
          )}

          <div className="grid md:grid-cols-2 gap-6">
            {/* 無料プラン */}
            <Card className="shadow-md border-gray-200">
              <CardHeader className="text-center bg-gray-50 rounded-t-lg">
                <CardTitle className="text-xl text-gray-700">無料プラン</CardTitle>
                <CardDescription>基本的な学習機能</CardDescription>
                <div className="mt-2">
                  <span className="text-3xl font-bold text-gray-800">¥0</span>
                  <span className="text-gray-500 text-sm"> / 月</span>
                </div>
              </CardHeader>
              <CardContent className="pt-6">
                <ul className="space-y-3">
                  {freeFeatures.map((feature) => (
                    <li key={feature} className="flex items-start gap-2 text-sm text-gray-700">
                      <Check className="h-4 w-4 text-green-500 shrink-0 mt-0.5" />
                      {feature}
                    </li>
                  ))}
                  <li className="flex items-start gap-2 text-sm text-gray-400">
                    <Lock className="h-4 w-4 shrink-0 mt-0.5" />
                    プレミアム限定カテゴリ（3つ）
                  </li>
                </ul>
              </CardContent>
              <CardFooter>
                <Link href="/categories" className="w-full">
                  <Button variant="outline" className="w-full border-gray-300 text-gray-700">
                    無料で始める
                  </Button>
                </Link>
              </CardFooter>
            </Card>

            {/* プレミアムプラン */}
            <Card className="shadow-lg border-pink-200 relative">
              <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                <span className="bg-pink-600 text-white text-xs font-bold px-3 py-1 rounded-full">
                  おすすめ
                </span>
              </div>
              <CardHeader className="text-center bg-pink-50 rounded-t-lg">
                <CardTitle className="text-xl text-pink-700">プレミアムプラン</CardTitle>
                <CardDescription>すべての機能が使い放題</CardDescription>
                <div className="mt-2">
                  <span className="text-3xl font-bold text-pink-700">¥480</span>
                  <span className="text-gray-500 text-sm"> / 月</span>
                </div>
              </CardHeader>
              <CardContent className="pt-6">
                <ul className="space-y-3">
                  {premiumFeatures.map((feature) => (
                    <li key={feature} className="flex items-start gap-2 text-sm text-gray-700">
                      <Check className="h-4 w-4 text-pink-500 shrink-0 mt-0.5" />
                      {feature}
                    </li>
                  ))}
                </ul>
              </CardContent>
              <CardFooter>
                <Button
                  onClick={handleCheckout}
                  disabled={loading}
                  className="w-full bg-pink-600 hover:bg-pink-700"
                >
                  {loading ? "処理中..." : "プレミアムにアップグレード"}
                </Button>
              </CardFooter>
            </Card>
          </div>

          <p className="text-center text-xs text-gray-400 mt-6">
            いつでもキャンセル可能。決済はStripeで安全に処理されます。
          </p>
        </div>
      </div>
    </div>
  )
}
