import Link from "next/link"
import { Book, Video, FileText, Home } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent } from "~/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs"

export default function Learning() {
  const articles = [
    {
      title: "妊娠中のパートナーへのサポート方法",
      description: "妊娠中の奥さんをどのようにサポートすべきか、実践的なアドバイス",
      icon: <FileText className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "出産の立ち会いについて知っておくべきこと",
      description: "出産に立ち会う際の心構えと準備",
      icon: <FileText className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "新生児の抱き方と寝かしつけのコツ",
      description: "赤ちゃんを安全に抱く方法と効果的な寝かしつけテクニック",
      icon: <FileText className="h-5 w-5 text-pink-500" />,
    },
  ]

  const videos = [
    {
      title: "おむつ交換の基本",
      description: "ステップバイステップで学ぶおむつ交換の方法",
      icon: <Video className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "赤ちゃんのお風呂の入れ方",
      description: "安全に赤ちゃんをお風呂に入れる方法",
      icon: <Video className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "ミルクの作り方と与え方",
      description: "正しいミルクの調乳方法と授乳のポイント",
      icon: <Video className="h-5 w-5 text-pink-500" />,
    },
  ]

  const books = [
    {
      title: "はじめてのパパの教科書",
      description: "出産前から育児までのパパの役割を解説",
      icon: <Book className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "赤ちゃんとの絆の作り方",
      description: "父親と赤ちゃんの絆を深める方法",
      icon: <Book className="h-5 w-5 text-pink-500" />,
    },
    {
      title: "パパの育児バイブル",
      description: "現代のお父さんのための実践的な育児ガイド",
      icon: <Book className="h-5 w-5 text-pink-500" />,
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
          <h1 className="text-2xl font-bold text-pink-600">学習リソース</h1>
          <div className="w-20"></div> {/* スペーサー */}
        </header>

        <div className="max-w-md mx-auto">
          <Tabs defaultValue="articles" className="w-full">
            <TabsList className="grid w-full grid-cols-3 mb-6">
              <TabsTrigger value="articles">記事</TabsTrigger>
              <TabsTrigger value="videos">動画</TabsTrigger>
              <TabsTrigger value="books">書籍</TabsTrigger>
            </TabsList>

            <TabsContent value="articles">
              <div className="grid gap-4">
                {articles.map((item, index) => (
                  <Card key={index} className="shadow-sm hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="p-2 rounded-full bg-pink-50">{item.icon}</div>
                        <div>
                          <h3 className="font-medium text-gray-800">{item.title}</h3>
                          <p className="text-sm text-gray-600">{item.description}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>

            <TabsContent value="videos">
              <div className="grid gap-4">
                {videos.map((item, index) => (
                  <Card key={index} className="shadow-sm hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="p-2 rounded-full bg-pink-50">{item.icon}</div>
                        <div>
                          <h3 className="font-medium text-gray-800">{item.title}</h3>
                          <p className="text-sm text-gray-600">{item.description}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>

            <TabsContent value="books">
              <div className="grid gap-4">
                {books.map((item, index) => (
                  <Card key={index} className="shadow-sm hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3">
                        <div className="p-2 rounded-full bg-pink-50">{item.icon}</div>
                        <div>
                          <h3 className="font-medium text-gray-800">{item.title}</h3>
                          <p className="text-sm text-gray-600">{item.description}</p>
                        </div>
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
