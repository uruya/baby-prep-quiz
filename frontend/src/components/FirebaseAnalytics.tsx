"use client";

import { useEffect } from "react";
import { analytics } from "~/lib/firebase";

export function FirebaseAnalytics() {
  useEffect(() => {
    // This will initialize analytics and log a page_view event
    if (typeof window !== "undefined") {
      // By simply importing and using analytics, Firebase initializes it.
      // You can log custom events here if needed, for example:
      // logEvent(analytics, 'page_view', { page_path: window.location.pathname });
      console.log("Firebase Analytics initialized for:", window.location.pathname);
    }
  }, []);

  return null;
}
