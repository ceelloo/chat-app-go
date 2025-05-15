"use client";

import { FormInput } from "@/components/form-input";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useLoginMutationRefactored } from "@/server/mutation";
import Link from "next/link";

export default function SignIn() {
  const { mutate: login, isPending } = useLoginMutationRefactored();

  const action = (form: FormData) => {
    const data = Object.fromEntries(form.entries()) as {
      email: string;
      password: string;
    };

    login(data);
  };

  return (
    <form action={action} className="w-96 border rounded-lg shadow-xl">
      <div className="px-10 py-6 flex flex-col gap-8">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Sign in to chat.</h1>
          <p className="text-sm text-zinc-500">
            Welcome back! please sign in to continue
          </p>
        </div>

        <div className="border border-dashed" />

        <FormInput
          label="Email"
          name="email"
          type="email"
          placeholder="your@email.com"
          autoComplete="email"
        />

        <FormInput
          label="Password"
          name="password"
          type="password"
          placeholder="yourpassword"
        />

        <Button disabled={isPending}>{isPending ? "..." : "Sign in"}</Button>
      </div>
      <div className="bg-muted border rounded-lg p-3">
        <p className="text-accent-foreground text-center text-sm">
          Don't have an account ?
          <Link
            href="/sign-up"
            className={cn(
              buttonVariants({ variant: "link" }),
              "px-2 text-sky-500"
            )}
          >
            Create account
          </Link>
        </p>
      </div>
    </form>
  );
}
