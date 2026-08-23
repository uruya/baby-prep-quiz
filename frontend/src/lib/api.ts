export type SubscriptionStatus = {
  tier: "free" | "premium"
  expiresAt?: string
}

export async function getSubscriptionStatus(): Promise<SubscriptionStatus> {
  const res = await fetch("/api/subscription/status")
  if (!res.ok) return { tier: "free" }
  return res.json()
}
