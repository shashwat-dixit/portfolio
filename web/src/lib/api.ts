import { CONFIG } from "@/data/config";

const BASE_URL = CONFIG.api.baseUrl;

export interface PostSummary {
  slug: string;
  title: string;
  description: string;
  cover: string;
  readingTime: number;
  tags: string[];
  date: string;
}

export interface Post {
  slug: string;
  title: string;
  description: string;
  contentHtml: string;
  tags: string[];
  date: string;
  updated: string;
  cover: string;
  readingTime: number;
  author: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface PostListResponse {
  posts: PostSummary[];
  pagination: Pagination;
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
  count: number;
}

export interface TagsResponse {
  tags: Tag[];
}

async function fetchAPI<T>(path: string): Promise<T> {
  const url = `${BASE_URL}${path}`;
  const res = await fetch(url, {
    headers: { "Accept": "application/json" },
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText} for ${url}`);
  }

  return res.json() as Promise<T>;
}

export async function getPosts(
  options: { tag?: string; page?: number; limit?: number } = {}
): Promise<PostListResponse> {
  const params = new URLSearchParams();
  if (options.tag) params.set("tag", options.tag);
  if (options.page) params.set("page", String(options.page));
  if (options.limit) params.set("limit", String(options.limit));

  const query = params.toString();
  return fetchAPI<PostListResponse>(`/api/posts${query ? `?${query}` : ""}`);
}

export async function getPostBySlug(slug: string): Promise<Post> {
  return fetchAPI<Post>(`/api/posts/${slug}`);
}

export async function getTags(): Promise<TagsResponse> {
  return fetchAPI<TagsResponse>("/api/tags");
}

export async function getAllPosts(): Promise<PostSummary[]> {
  const resp = await getPosts({ limit: 100 });
  return resp.posts;
}
