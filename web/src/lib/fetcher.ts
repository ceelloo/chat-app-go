import { BASE_URL } from "./constant";

type FetchOptions = RequestInit & {
  json?: boolean;
};

export type ErrorResponse = {
  error: boolean;
  message: string;
  code: number;
};

export async function apiFetch<T = any>(
  url: string,
  options: FetchOptions = {}
): Promise<T> {
  const { json = true, headers, ...rest } = options;

  const res = await fetch(`${BASE_URL}${url}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
    ...rest,
  });

  if (!res.ok) {
    let errorData: ErrorResponse;
    try {
      errorData = await res.json();
    } catch (_) {
      throw new Error(res.statusText);
    }
    throw new Error(errorData.message || "Request failed");
  }

  return json ? res.json() : (res as any);
}
