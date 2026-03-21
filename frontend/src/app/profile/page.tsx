"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Home, Trophy, BookOpen, LogOut } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card"
import { Progress } from "~/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs"

type UserData = {
  name: string
  email: string
}

type CategoryStat = {
  bestScore: number
  total: number
}

type Stats = {
  completedQuizzes: number
  totalScore: number
  totalPossible: number
  categories: Record<string, CategoryStat>
}

type Badge = {
  name: string
  description: string
  icon: React.ReactNode
}

type Achievement = {
  name: string
  category: string
  progress: number
}

export default function Profile() {
  const router = useRouter()
  const [user, setUser] = useState<UserData | null>(null)
  const [stats, setStats] = useState<Stats | null>(null)

  useEffect(() => {
    Promise.all([
      fetch(`/api/auth/me`),
      fetch(`/api/quiz/stats`),
    ]).then(async ([meRes, statsRes]) => {
      if (!meRes.ok) {
        router.push("/auth/login")
        return
      }
      setUser(await meRes.json())
      if (statsRes.ok) setStats(await statsRes.json())
    })
  }, [router])

  const handleLogout = async () => {
    await fetch(`/api/auth/logout`, { method: "POST" })
    router.push("/")
  }

  const badges: Badge[] = [
    { name: "初心者パパ", description: "最初のクイズを完了", icon: <Trophy className="h-5 w-5 text-amber-500" /> },
    { name: "知識収集家", description: "5つのクイズを完了", icon: <BookOpen className="h-5 w-5 text-blue-500" /> },
  ]

  const achievements: Achievement[] = [
    { name: "妊娠の基礎知識マスター", category: "pregnancy", progress: 0 },
    { name: "出産の準備エキスパート", category: "birth", progress: 0 },
    { name: "赤ちゃんのお世話の達人", category: "baby-care", progress: 0 },
    { name: "栄養と食事の専門家", category: "nutrition", progress: 0 },
    { name: "発達サポーター", category: "development", progress: 0 },
  ].map((a) => {
    const cat = stats?.categories[a.category]
    return { ...a, progress: cat ? Math.round((cat.bestScore / cat.total) * 100) : 0 }
  })

  const totalProgress = stats && stats.totalPossible > 0
    ? Math.round((stats.totalScore / stats.totalPossible) * 100)
    : 0

  if (user === null) return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50 flex items-center justify-center">
      <p className="text-gray-500">読み込み中...</p>
    </div>
  )

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
          <h1 className="text-xl md:text-2xl font-bold text-pink-600">マイページ</h1>
          <Button variant="ghost" size="sm" className="flex items-center gap-1 text-gray-500 shrink-0" onClick={handleLogout}>
            <LogOut className="h-4 w-4" />
            ログアウト
          </Button>
        </header>

        <div className="max-w-md mx-auto">
          <Card className="shadow-lg border-pink-100 mb-6">
            <CardHeader className="text-center bg-pink-50 rounded-t-lg pb-2">
              <div className="mx-auto mb-2 bg-white p-2 rounded-full w-20 h-20 flex items-center justify-center">
                <span className="text-3xl font-bold text-pink-600">{user.name.charAt(0)}</span>
              </div>
              <CardTitle>{user.name}</CardTitle>
              <CardDescription>{user.email}</CardDescription>
            </CardHeader>
            <CardContent className="pt-4">
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>総合進捗</span>
                    <span>{totalProgress}%</span>
                  </div>
                  <Progress value={totalProgress} className="h-2" />
                </div>

                <div className="grid grid-cols-2 gap-4 text-center">
                  <div className="bg-pink-50 rounded-lg p-3">
                    <div className="text-2xl font-bold text-pink-600">{stats?.completedQuizzes ?? 0}</div>
                    <div className="text-xs text-gray-600">完了したクイズ</div>
                  </div>
                  <div className="bg-blue-50 rounded-lg p-3">
                    <div className="text-2xl font-bold text-blue-600">{stats?.totalScore ?? 0}</div>
                    <div className="text-xs text-gray-600">獲得ポイント</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Tabs defaultValue="achievements" className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-6">
              <TabsTrigger value="achievements">達成状況</TabsTrigger>
              <TabsTrigger value="badges">バッジ</TabsTrigger>
            </TabsList>

            <TabsContent value="achievements">
              <div className="space-y-4">
                {achievements.map((achievement, index) => (
                  <div key={index} className="bg-white rounded-lg p-4 shadow-sm">
                    <div className="flex justify-between text-sm mb-1">
                      <span className="font-medium">{achievement.name}</span>
                      <span>{achievement.progress}%</span>
                    </div>
                    <Progress value={achievement.progress} className="h-2" />
                  </div>
                ))}
              </div>
            </TabsContent>

            <TabsContent value="badges">
              <div className="grid grid-cols-2 gap-4">
                {badges.map((badge, index) => (
                  <Card key={index} className="shadow-sm border-pink-100">
                    <CardContent className="p-4 flex items-center gap-3">
                      <div className="p-2 rounded-full bg-pink-50">{badge.icon}</div>
                      <div>
                        <h3 className="font-medium text-gray-800">{badge.name}</h3>
                        <p className="text-xs text-gray-600">{badge.description}</p>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}
