"use client";

import { useEffect } from "react";
import { analytics } from "~/lib/firebase";

export function FirebaseAnalytics() {
  useEffect(() => {
    if (typeof window !== "undefined") {
      void analytics; // analyticsをインポートすることで初期化される
      console.log("Firebase Analytics initialized for:", window.location.pathname);
    }
  }, []);

  return null;
}
