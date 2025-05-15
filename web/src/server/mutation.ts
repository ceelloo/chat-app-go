import { BASE_URL } from "@/lib/constant";
import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { toast } from "sonner";

interface ApiError extends Error {
  error: boolean;
  message: string;
  code: number;
}

async function postData<TResponse, TData>(endpoint: string, data: TData): Promise<TResponse> {
  const response = await fetch(`${BASE_URL}${endpoint}`, {
    method: "POST",
    body: JSON.stringify(data),
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const text = await response.text();
  let jsonResponse: any;

  try {
    jsonResponse = JSON.parse(text);
  } catch (error) {
    if (response.ok) {
      console.warn("Response was OK but not valid JSON:", text);
    }
    throw new Error("Invalid JSON response from server");
  }

  if (!response.ok) {
    const errorMessage = (jsonResponse as { message?: string })?.message || `Request failed for endpoint: ${endpoint}`;
    const customError = new Error(errorMessage) as ApiError;
    throw customError;
  }

  return jsonResponse as TResponse;
}

interface AuthPayload {
  email: string;
  password: string;
}

interface AuthResponse {
  token?: string;
  userId?: string;
  message?: string;
}

interface CreateAuthMutationOptions<TResponse, TData, TError = ApiError> {
  mutationKey: string[];
  endpoint: string;
  successMessage?: string | ((data: TResponse, variables: TData) => string);
  errorMessage?: string | ((error: TError, variables: TData) => string);
  mutationOptions?: Omit<UseMutationOptions<TResponse, TError, TData>, 'mutationFn' | 'mutationKey'>;
}

function createAuthMutation<
  TResponse extends { message?: string },
  TData extends AuthPayload,
  TError extends ApiError = ApiError
>(options: CreateAuthMutationOptions<TResponse, TData, TError>) {
  return useMutation<TResponse, TError, TData>({
    mutationKey: options.mutationKey,
    mutationFn: (data: TData) => postData<TResponse, TData>(options.endpoint, data),
    onSuccess: (data, variables) => {
      let message: string;
      if (typeof options.successMessage === 'function') {
        message = options.successMessage(data, variables);
      } else {
        message = options.successMessage || data.message || "Operation successful!";
      }
      toast.success(message);

      setTimeout(() => {
        window.location.href = "/"
      }, 4000);

      if (options.mutationOptions?.onSuccess) {
        options.mutationOptions.onSuccess(data, variables, undefined)
      }
    },
    onError: (error, variables) => {
      let message: string;
      if (typeof options.errorMessage === 'function') {
        message = options.errorMessage(error, variables);
      } else {
        message = options.errorMessage || error.message || "An error occurred.";
      }
      toast.error(message);

      if (options.mutationOptions?.onError) {
        options.mutationOptions.onError(error, variables, undefined)
      }
    },
    onSettled(data, error, variables) {
      setTimeout(() => {
        toast.dismiss()
      }, 3000);
      if (options.mutationOptions?.onSettled) {
        options.mutationOptions.onSettled(data, error, variables, undefined);
      }
    },
    ...(options.mutationOptions || {}),
  });
}

export const useLoginMutationRefactored = () =>
  createAuthMutation<AuthResponse, AuthPayload>({
    mutationKey: ["auth.login"],
    endpoint: "/authentication/login",
    successMessage: (data) => data.message || "Login successful!",
    errorMessage: (error) => error.message || "Login failed. Please try again.",
  });

export const useRegisterMutationRefactored = () =>
  createAuthMutation<AuthResponse, AuthPayload>({
    mutationKey: ["auth.register"],
    endpoint: "/authentication/register",
    successMessage: "Registration successful! You can now log in.",
    errorMessage: (error) => error.message || "Registration failed. Please check your input.",
  });