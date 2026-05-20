import { useState, useEffect, useCallback } from "react";

const STORAGE_KEY = "blog-read-progress";
const MAX_AGE_MS = 30 * 24 * 60 * 60 * 1000; // 30 days

interface ProgressEntry {
  scrollPercent: number;
  lastPosition: number;
  contentHeight: number;
  lastRead: number;
}

type ProgressStore = Record<string, ProgressEntry>;

function getStore(): ProgressStore {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return {};
    return JSON.parse(raw);
  } catch {
    return {};
  }
}

function setStore(store: ProgressStore) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
  } catch {
    // localStorage may be full or unavailable
  }
}

function cleanOldEntries(store: ProgressStore): ProgressStore {
  const now = Date.now();
  const cleaned: ProgressStore = {};
  for (const [slug, entry] of Object.entries(store)) {
    if (now - entry.lastRead < MAX_AGE_MS) {
      cleaned[slug] = entry;
    }
  }
  return cleaned;
}

export function useReadingProgress(slug: string) {
  const [progress, setProgress] = useState(0);

  const saveProgress = useCallback(() => {
    const scrollTop = window.scrollY;
    const docHeight = document.documentElement.scrollHeight - window.innerHeight;
    const percent = docHeight > 0 ? Math.round((scrollTop / docHeight) * 100) : 0;

    setProgress(percent);

    const store = getStore();
    store[slug] = {
      scrollPercent: percent,
      lastPosition: scrollTop,
      contentHeight: document.documentElement.scrollHeight,
      lastRead: Date.now(),
    };
    setStore(store);
  }, [slug]);

  useEffect(() => {
    const store = cleanOldEntries(getStore());
    setStore(store);

    const entry = store[slug];
    if (entry && entry.lastPosition > 0) {
      const heightRatio = document.documentElement.scrollHeight / entry.contentHeight;
      const adjustedPosition = Math.round(entry.lastPosition * heightRatio);
      window.scrollTo(0, adjustedPosition);
      setProgress(entry.scrollPercent);
    }
  }, [slug]);

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | null = null;

    const handleScroll = () => {
      if (timer) clearTimeout(timer);
      timer = setTimeout(saveProgress, 300);

      const scrollTop = window.scrollY;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;
      const percent = docHeight > 0 ? Math.round((scrollTop / docHeight) * 100) : 0;
      setProgress(percent);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
      if (timer) clearTimeout(timer);
    };
  }, [saveProgress]);

  return { progress };
}
