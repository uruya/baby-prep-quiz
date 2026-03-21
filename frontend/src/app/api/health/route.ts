import { NextResponse } from "next/server"

export async function GET() {
  const backendUrl = process.env.BACKEND_URL ?? "NOT SET"

  let fetchResult = "unknown"
  let errorDetail = ""
  try {
    const res = await fetch(`${backendUrl}/api/auth/me`)
    fetchResult = `${res.status}`
  } catch (e) {
    fetchResult = "fetch failed"
    errorDetail = String(e)
  }

  return NextResponse.json({
    backendUrl,
    fetchResult,
    errorDetail,
  })
}
