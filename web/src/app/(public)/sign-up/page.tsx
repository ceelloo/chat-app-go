"use client";

import { FormInput } from "@/components/form-input";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useRegisterMutationRefactored } from "@/server/mutation";
import Link from "next/link";
import { useFormStatus } from "react-dom";

export default function SignUp() {
  const { mutate: register, isPending } = useRegisterMutationRefactored();
  
    const action = (form: FormData) => {
      const data = Object.fromEntries(form.entries()) as {
        email: string;
        password: string;
      };
  
      register(data);
    };

  return (
    <form action={action} className="w-96 border rounded-lg shadow-xl">
      <div className="px-10 py-6 flex flex-col gap-8">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Create account</h1>
          <p className="text-sm text-zinc-500">
            Let's get started. Fill in the details below to create your account.
          </p>
        </div>

        <div className="border border-dashed" />

        <FormInput
          label="Name"
          name="name"
          placeholder="yourname"
          autoComplete="name"
          required
        />

        <FormInput
          label="Email"
          name="email"
          type="email"
          placeholder="your@email.com"
          autoComplete="email"
          required
        />

        <FormInput
          label="Password"
          name="password"
          type="password"
          placeholder="yourpassword"
          required
        />

        <Button disabled={isPending}>{isPending ? "..." : "Create account"}</Button>
      </div>

      <div className="bg-muted border rounded-lg p-3">
        <p className="text-accent-foreground text-center text-sm">
          Already have an account ?
          <Link
            href="/sign-in"
            className={cn(
              buttonVariants({ variant: "link" }),
              "px-2 text-blue-500"
            )}
          >
            Sign in
          </Link>
        </p>
      </div>
    </form>
  );
}
