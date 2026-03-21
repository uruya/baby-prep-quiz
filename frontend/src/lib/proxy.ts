import { NextRequest, NextResponse } from "next/server"

const BACKEND_URL = process.env.BACKEND_URL || "http://127.0.0.1:8080"

export async function proxyToBackend(
  request: NextRequest,
  backendPath: string,
): Promise<NextResponse> {
  const url = new URL(request.url)
  const backendUrl = `${BACKEND_URL}${backendPath}${url.search}`

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  }

  const cookieHeader = request.headers.get("cookie")
  if (cookieHeader) {
    headers["Cookie"] = cookieHeader
  }

  const method = request.method
  let body: string | undefined
  if (method !== "GET" && method !== "HEAD") {
    body = await request.text()
  }

  try {
    const backendResponse = await fetch(backendUrl, { method, headers, body })

    // 204 No Content はbody不可のため別処理
    if (backendResponse.status === 204) {
      const response = new NextResponse(null, { status: 204 })
      const setCookie = backendResponse.headers.get("set-cookie")
      if (setCookie) response.headers.set("Set-Cookie", setCookie)
      return response
    }

    const responseData = await backendResponse.text()
    const response = new NextResponse(responseData, {
      status: backendResponse.status,
      headers: { "Content-Type": "application/json" },
    })

    const setCookie = backendResponse.headers.get("set-cookie")
    if (setCookie) {
      response.headers.set("Set-Cookie", setCookie)
    }

    return response
  } catch (error) {
    console.error(`[proxy] ${method} ${backendUrl} failed:`, error)
    return new NextResponse(JSON.stringify({ message: "Backend unreachable" }), {
      status: 502,
      headers: { "Content-Type": "application/json" },
    })
  }
}
