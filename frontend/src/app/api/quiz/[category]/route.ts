import { NextRequest } from "next/server"
import { proxyToBackend } from "~/lib/proxy"

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ category: string }> },
) {
  const { category } = await params
  return proxyToBackend(request, `/api/quiz/${category}`)
}
