"use client"

import { useState } from "react"
import Link from "next/link"
import { Home, Trophy, BookOpen } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card"
import { Progress } from "~/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs"

export default function Profile() {
  // 仮のユーザーデータ
  const [userData] = useState({
    name: "山田太郎",
    completedQuizzes: 8,
    totalQuizzes: 15,
    score: 24,
    maxScore: 45,
    badges: [
      { name: "初心者パパ", description: "最初のクイズを完了", icon: <Trophy className="h-5 w-5 text-amber-500" /> },
      { name: "知識収集家", description: "5つのクイズを完了", icon: <BookOpen className="h-5 w-5 text-blue-500" /> },
    ],
    achievements: [
      { name: "妊娠の基礎知識マスター", progress: 66 },
      { name: "出産の準備エキスパート", progress: 40 },
      { name: "赤ちゃんのお世話の達人", progress: 20 },
    ],
  })

  const totalProgress = Math.round((userData.score / userData.maxScore) * 100)

  return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50">
      <div className="container mx-auto px-4 py-8">
        <header className="flex items-center justify-between mb-8">
          <Link href="/">
            <Button variant="ghost" size="sm" className="flex items-center gap-1">
              <Home className="h-4 w-4" />
              ホーム
            </Button>
          </Link>
          <h1 className="text-2xl font-bold text-pink-600">マイページ</h1>
          <div className="w-20"></div> {/* スペーサー */}
        </header>

        <div className="max-w-md mx-auto">
          <Card className="shadow-lg border-pink-100 mb-6">
            <CardHeader className="text-center bg-pink-50 rounded-t-lg pb-2">
              <div className="mx-auto mb-2 bg-white p-2 rounded-full w-20 h-20 flex items-center justify-center">
                <span className="text-3xl font-bold text-pink-600">{userData.name.charAt(0)}</span>
              </div>
              <CardTitle>{userData.name}</CardTitle>
              <CardDescription>もうすぐパパになる準備中</CardDescription>
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
                    <div className="text-2xl font-bold text-pink-600">{userData.completedQuizzes}</div>
                    <div className="text-xs text-gray-600">完了したクイズ</div>
                  </div>
                  <div className="bg-blue-50 rounded-lg p-3">
                    <div className="text-2xl font-bold text-blue-600">{userData.score}</div>
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
                {userData.achievements.map((achievement, index) => (
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
                {userData.badges.map((badge, index) => (
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
