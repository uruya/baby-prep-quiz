"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { ArrowLeft, Home } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "~/components/ui/card"
import { Progress } from "~/components/ui/progress"
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group"
import { Label } from "~/components/ui/label"

type Question = {
  id: number
  category: string
  question: string
  options: string[]
  correctAnswer: number
  explanation: string
}

export default function Quiz({ params }: { params: { category: string } }) {
  const category = params.category
  const [questions, setQuestions] = useState<Question[]>([])
  const [loading, setLoading] = useState(true)
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0)
  const [selectedOption, setSelectedOption] = useState<number | null>(null)
  const [showAnswer, setShowAnswer] = useState(false)
  const [score, setScore] = useState(0)
  const [quizCompleted, setQuizCompleted] = useState(false)

  useEffect(() => {
    const base = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"
    fetch(`${base}/quiz/${category}`)
      .then((res) => res.json())
      .then((data: Question[]) => {
        setQuestions(data)
        setLoading(false)
      })
      .catch((err
      ) => {
        console.error(err)
        setLoading(false)
      })
  }, [category])

  // カテゴリー名の日本語マッピング
  const categoryNames: { [key: string]: string } = {
    pregnancy: "妊娠の基礎知識",
    birth: "出産の準備",
    "baby-care": "赤ちゃんのお世話",
    nutrition: "栄養と食事",
    development: "赤ちゃんの発達",
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-pink-50 to-blue-50">
      <p>Loading...</p>
    </div>
    )
  }

  if (questions.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-pink-50 to-blue-50">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-center text-red-600">エラー</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-center">このカテゴリーのクイズは現在準備中です。</p>
          </CardContent>
          <CardFooter className="flex justify-center">
            <Link href="/categories">
              <Button>カテゴリー選択に戻る</Button>
            </Link>
          </CardFooter>
        </Card>
      </div>
    )
  }

  const currentQuestion = questions[currentQuestionIndex]
  const progress = ((currentQuestionIndex + 1) / questions.length) * 100

  const handleOptionSelect = (index: number) => {
    if (!showAnswer) {
      setSelectedOption(index)
    }
  }

  const handleCheckAnswer = () => {
    if (selectedOption === currentQuestion.correctAnswer) {
      setScore(score + 1)
    }
    setShowAnswer(true)
  }

  const handleNextQuestion = () => {
    if (currentQuestionIndex < questions.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1)
      setSelectedOption(null)
      setShowAnswer(false)
    } else {
      setQuizCompleted(true)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-pink-50 to-blue-50">
      <div className="container mx-auto px-4 py-8">
        <header className="flex items-center justify-between mb-8">
          <Link href="/categories">
            <Button variant="ghost" size="sm" className="flex items-center gap-1">
              <ArrowLeft className="h-4 w-4" />
              戻る
            </Button>
          </Link>
          <h1 className="text-xl font-bold text-pink-600">{categoryNames[category] || category}</h1>
          <Link href="/">
            <Button variant="ghost" size="sm" className="flex items-center gap-1">
              <Home className="h-4 w-4" />
              ホーム
            </Button>
          </Link>
        </header>

        <div className="max-w-md mx-auto">
          {!quizCompleted ? (
            <Card className="shadow-lg">
              <CardHeader className="pb-2">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm text-gray-500">
                    問題 {currentQuestionIndex + 1}/{questions.length}
                  </span>
                  <span className="text-sm text-gray-500">スコア: {score}</span>
                </div>
                <Progress value={progress} className="h-2" />
                <CardTitle className="mt-4 text-lg">{currentQuestion.question}</CardTitle>
              </CardHeader>
              <CardContent>
                <RadioGroup value={selectedOption?.toString()} className="space-y-3">
                  {currentQuestion.options.map((option, index) => (
                    <div
                      key={index}
                      className={`flex items-center space-x-2 rounded-lg border p-3 cursor-pointer transition-colors ${
                        showAnswer
                          ? index === currentQuestion.correctAnswer
                            ? "bg-green-50 border-green-200"
                            : selectedOption === index
                              ? "bg-red-50 border-red-200"
                              : ""
                          : selectedOption === index
                            ? "bg-blue-50 border-blue-200"
                            : ""
                      }`}
                      onClick={() => handleOptionSelect(index)}
                    >
                      <RadioGroupItem value={index.toString()} id={`option-${index}`} disabled={showAnswer} />
                      <Label htmlFor={`option-${index}`} className="w-full cursor-pointer">
                        {option}
                      </Label>
                    </div>
                  ))}
                </RadioGroup>

                {showAnswer && (
                  <div className="mt-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
                    <p className="text-sm font-medium text-blue-800">解説:</p>
                    <p className="text-sm text-blue-700">{currentQuestion.explanation}</p>
                  </div>
                )}
              </CardContent>
              <CardFooter className="flex justify-center">
                {!showAnswer ? (
                  <Button
                    onClick={handleCheckAnswer}
                    disabled={selectedOption === null}
                    className="w-full bg-pink-600 hover:bg-pink-700"
                  >
                    回答を確認
                  </Button>
                ) : (
                  <Button onClick={handleNextQuestion} className="w-full bg-pink-600 hover:bg-pink-700">
                    {currentQuestionIndex < questions.length - 1 ? "次の問題へ" : "結果を見る"}
                  </Button>
                )}
              </CardFooter>
            </Card>
          ) : (
            <Card className="shadow-lg text-center">
              <CardHeader>
                <CardTitle className="text-2xl text-pink-600">クイズ完了！</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="mb-6">
                  <div className="text-5xl font-bold text-pink-600 mb-2">
                    {score}/{questions.length}
                  </div>
                  <p className="text-gray-600">
                    {score === questions.length
                      ? "素晴らしい！満点です！"
                      : score >= questions.length / 2
                        ? "よくできました！"
                        : "もう一度挑戦してみましょう！"}
                  </p>
                </div>
                <div className="flex flex-col gap-3">
                  <Link href="/categories">
                    <Button variant="outline" className="w-full">
                      他のカテゴリーに挑戦
                    </Button>
                  </Link>
                  <Link href={`/quiz/${category}`}>
                    <Button
                      onClick={() => {
                        setCurrentQuestionIndex(0)
                        setSelectedOption(null)
                        setShowAnswer(false)
                        setScore(0)
                        setQuizCompleted(false)
                      }}
                      className="w-full bg-pink-600 hover:bg-pink-700"
                    >
                      もう一度挑戦
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
