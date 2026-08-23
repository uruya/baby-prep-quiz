"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { Baby, Heart, Stethoscope, Utensils, Brain, Home, Lock, Crown } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "~/components/ui/card"
import { getSubscriptionStatus } from "~/lib/api"

const FREE_CATEGORIES = ["pregnancy", "birth"]

export default function Categories() {
  const [isPremium, setIsPremium] = useState(false)
  const [showUpgradeModal, setShowUpgradeModal] = useState(false)

  useEffect(() => {
    getSubscriptionStatus().then((s) => setIsPremium(s.tier === "premium"))
  }, [])

  const categories = [
    {
      id: "pregnancy",
      title: "妊娠の基礎知識",
      description: "妊娠期間中の変化と注意点",
      icon: <Heart className="h-8 w-8 text-pink-500" />,
      color: "bg-pink-50 border-pink-200",
    },
    {
      id: "birth",
      title: "出産の準備",
      description: "出産時に必要な知識と心構え",
      icon: <Stethoscope className="h-8 w-8 text-blue-500" />,
      color: "bg-blue-50 border-blue-200",
    },
    {
      id: "baby-care",
      title: "赤ちゃんのお世話",
      description: "おむつ交換や授乳の基本",
      icon: <Baby className="h-8 w-8 text-purple-500" />,
      color: "bg-purple-50 border-purple-200",
    },
    {
      id: "nutrition",
      title: "栄養と食事",
      description: "ママと赤ちゃんの健康的な食事",
      icon: <Utensils className="h-8 w-8 text-green-500" />,
      color: "bg-green-50 border-green-200",
    },
    {
      id: "development",
      title: "赤ちゃんの発達",
      description: "成長段階と発達のサポート",
      icon: <Brain className="h-8 w-8 text-amber-500" />,
      color: "bg-amber-50 border-amber-200",
    },
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
          <h1 className="text-xl md:text-2xl font-bold text-pink-600">クイズカテゴリー</h1>
          <div className="w-20 shrink-0"></div>
        </header>

        <div className="max-w-md mx-auto">
          <p className="text-center text-gray-600 mb-6">挑戦したいカテゴリーを選んでください</p>

          <div className="grid gap-4">
            {categories.map((category) => {
              const isLocked = !FREE_CATEGORIES.includes(category.id) && !isPremium
              if (isLocked) {
                return (
                  <button
                    key={category.id}
                    className="block w-full text-left"
                    onClick={() => setShowUpgradeModal(true)}
                  >
                    <Card className={`shadow-sm hover:shadow-md transition-shadow ${category.color} opacity-75`}>
                      <CardContent className="p-4">
                        <div className="flex items-center gap-4">
                          <div className="p-2 rounded-full bg-white shadow-sm">{category.icon}</div>
                          <div className="flex-1">
                            <div className="flex items-center gap-2 flex-wrap">
                              <h2 className="font-semibold text-gray-800">{category.title}</h2>
                              <span className="inline-flex items-center gap-1 bg-amber-100 text-amber-700 text-xs font-medium px-2 py-0.5 rounded-full border border-amber-200">
                                <Lock className="h-3 w-3" />
                                プレミアム限定
                              </span>
                            </div>
                            <p className="text-sm text-gray-600">{category.description}</p>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </button>
                )
              }
              return (
                <Link key={category.id} href={`/quiz/${category.id}`} className="block">
                  <Card className={`shadow-sm hover:shadow-md transition-shadow ${category.color}`}>
                    <CardContent className="p-4">
                      <div className="flex items-center gap-4">
                        <div className="p-2 rounded-full bg-white shadow-sm">{category.icon}</div>
                        <div>
                          <h2 className="font-semibold text-gray-800">{category.title}</h2>
                          <p className="text-sm text-gray-600">{category.description}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </Link>
              )
            })}
          </div>

          {!isPremium && (
            <div className="mt-6 p-4 bg-white rounded-lg border border-pink-100 shadow-sm text-center">
              <p className="text-sm text-gray-600 mb-3">
                プレミアムプランで全カテゴリーを解放しよう
              </p>
              <Link href="/pricing">
                <Button className="bg-pink-600 hover:bg-pink-700">プランを見る</Button>
              </Link>
            </div>
          )}
        </div>
      </div>

      {showUpgradeModal && (
        <UpgradeModal onClose={() => setShowUpgradeModal(false)} />
      )}
    </div>
  )
}

function UpgradeModal({ onClose }: { onClose: () => void }) {
  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
      onClick={onClose}
    >
      <Card
        className="w-full max-w-sm shadow-xl border-pink-200"
        onClick={(e) => e.stopPropagation()}
      >
        <CardHeader className="text-center bg-pink-50 rounded-t-lg">
          <div className="mx-auto mb-2 bg-pink-100 p-3 rounded-full w-14 h-14 flex items-center justify-center">
            <Crown className="h-7 w-7 text-pink-600" />
          </div>
          <CardTitle className="text-pink-700">プレミアムプランにアップグレード</CardTitle>
        </CardHeader>
        <CardContent className="pt-4 text-center">
          <p className="text-gray-600 text-sm mb-3">
            このカテゴリーはプレミアム会員限定です。
          </p>
          <p className="text-gray-700 text-sm">
            月額 <span className="font-bold text-pink-600">¥480</span> で全5カテゴリーが使い放題になります。
          </p>
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Link href="/pricing" className="w-full">
            <Button className="w-full bg-pink-600 hover:bg-pink-700">
              今すぐアップグレード
            </Button>
          </Link>
          <Button variant="ghost" className="w-full text-gray-500" onClick={onClose}>
            キャンセル
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}
