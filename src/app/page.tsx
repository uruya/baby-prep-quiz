import Link from "next/link"
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card"

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50">
      <div className="container mx-auto px-4 py-8">
        <header className="text-center mb-12">
          <h1 className="text-3xl md:text-4xl font-bold text-pink-600 mb-2">パパクイズ</h1>
          <p className="text-gray-600">赤ちゃんが生まれる前に知っておきたいこと</p>
        </header>

        <div className="max-w-md mx-auto">
          <Card className="shadow-lg border-pink-100 mb-6">
            <CardHeader className="text-center bg-pink-50 rounded-t-lg">
              <CardTitle className="text-2xl text-pink-700">ようこそ！</CardTitle>
              <CardDescription>これから父親になる方のための準備クイズアプリです</CardDescription>
            </CardHeader>
            <CardContent className="pt-6">
              <p className="text-gray-700 mb-4">
                赤ちゃんの誕生は人生で最も素晴らしい瞬間の一つです。このアプリでは、出産前に知っておくべき知識をクイズ形式で学べます。
              </p>
              <div className="flex justify-center">
                <img
                  src="/placeholder.svg?height=150&width=150"
                  alt="赤ちゃんのイラスト"
                  className="rounded-full border-4 border-pink-100"
                />
              </div>
            </CardContent>
            <CardFooter className="flex justify-center pb-6">
              <Link href="/categories" className="w-full">
                <Button className="w-full bg-pink-600 hover:bg-pink-700">クイズを始める</Button>
              </Link>
            </CardFooter>
          </Card>

          <div className="grid grid-cols-2 gap-4">
            <Link href="/learning" className="w-full">
              <Button variant="outline" className="w-full border-pink-200 text-pink-700 hover:bg-pink-50">
                学習リソース
              </Button>
            </Link>
            <Link href="/profile" className="w-full">
              <Button variant="outline" className="w-full border-pink-200 text-pink-700 hover:bg-pink-50">
                マイページ
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
