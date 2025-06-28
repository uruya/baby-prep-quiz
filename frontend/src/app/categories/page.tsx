import Link from "next/link"
import { Baby, Heart, Stethoscope, Utensils, Brain, Home } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent } from "~/components/ui/card"

export default function Categories() {
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
            <Button variant="ghost" size="sm" className="flex items-center gap-1">
              <Home className="h-4 w-4" />
              ホーム
            </Button>
          </Link>
          <h1 className="text-2xl font-bold text-pink-600">クイズカテゴリー</h1>
          <div className="w-20"></div> {/* スペーサー */}
        </header>

        <div className="max-w-md mx-auto">
          <p className="text-center text-gray-600 mb-6">挑戦したいカテゴリーを選んでください</p>

          <div className="grid gap-4">
            {categories.map((category) => (
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
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
