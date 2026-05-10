import { useState, type FormEvent } from "react";
import { EyeIcon, EyeOffIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Signin, Signup } from "../../../wailsjs/go/auth/Auth";

type AuthMode = "signin" | "signup";
type cred = {
  email: string;
  password: string;
};

export function AuthScene({
  handleNav,
}: {
  handleNav: (p: "auth" | "app" | "payment") => void;
}) {
  const [mode, setMode] = useState<AuthMode>("signin");
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    try {
      if (mode === "signin") {
        await Signin({ email, password });
      } else {
        await Signup({ email, password });
      }

      handleNav("app");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed");
    }
  };

  return (
    <div className="flex min-h-svh items-center justify-center bg-[#f6f6f4] px-4 py-8">
      <Card className="w-full max-w-md rounded-2xl border border-black/5 bg-white shadow-[0_18px_40px_rgba(15,23,42,0.08)]">
        <CardHeader className="space-y-2 px-6 pt-6 pb-0">
          <p className="text-xs font-medium uppercase tracking-[0.22em] text-muted-foreground">
            Account
          </p>
          <CardTitle className="text-3xl font-semibold tracking-tight">
            {mode === "signin" ? "Sign In" : "Sign Up"}
          </CardTitle>
          <CardDescription className="max-w-[28ch] text-sm leading-6 text-muted-foreground">
            {mode === "signin"
              ? "Enter your email and password to continue."
              : "Create an account with your email and password."}
          </CardDescription>
        </CardHeader>
        <CardContent className="px-6 pb-6 pt-5">
          <form
            className="flex flex-col gap-4"
            onSubmit={handleSubmit}
          >
            <FieldGroup className="gap-3">
              <Field className="gap-2">
                <FieldLabel htmlFor={`${mode}-email`}>Email</FieldLabel>
                <Input
                  id={`${mode}-email`}
                  type="email"
                  inputSize="lg"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="h-11 rounded-xl border-black/10 bg-[#fafafa] px-3.5"
                />
              </Field>
              <Field className="gap-2">
                <FieldLabel htmlFor={`${mode}-password`}>Password</FieldLabel>
                <div className="relative">
                  <Input
                    id={`${mode}-password`}
                    type={showPassword ? "text" : "password"}
                    inputSize="lg"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="h-11 rounded-xl border-black/10 bg-[#fafafa] pr-12 px-3.5"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword((value) => !value)}
                    className="absolute inset-y-0 right-0 flex cursor-pointer items-center px-4 text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {showPassword ? (
                      <EyeOffIcon className="size-4" />
                    ) : (
                      <EyeIcon className="size-4" />
                    )}
                  </button>
                </div>
              </Field>
            </FieldGroup>

            <Button
              type="submit"
              size="lg"
              className="h-11 w-full rounded-xl cursor-pointer bg-black text-white shadow-none transition-colors hover:bg-zinc-800 active:scale-[0.98]"
            >
              {mode === "signin" ? "Sign In" : "Sign Up"}
            </Button>

            {error ? (
              <p className="text-sm text-destructive">{error}</p>
            ) : null}

            <div className="pt-1 text-sm text-muted-foreground">
              {mode === "signin" ? (
                <p>
                  Don&apos;t have an account?{" "}
                  <button
                    type="button"
                    onClick={() => setMode("signup")}
                    className="cursor-pointer font-medium text-foreground underline-offset-4 transition-colors hover:text-primary"
                  >
                    Sign up
                  </button>
                </p>
              ) : (
                <p>
                  Already have an account?{" "}
                  <button
                    type="button"
                    onClick={() => setMode("signin")}
                    className="cursor-pointer font-medium text-foreground underline-offset-4 transition-colors hover:text-primary"
                  >
                    Sign in
                  </button>
                </p>
              )}
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
