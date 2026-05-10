import { useState } from "react";
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

export function AuthScene({ handleNav }: { handleNav: (p: string) => void }) {
  const [mode, setMode] = useState<AuthMode>("signin");
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSignin = (cred: cred) => {
    Signin(cred).then(() => {
      handleNav("app");
    });
  };

  const handleSignup = (cred: cred) => {
    Signup(cred).then(() => {
      handleNav("app");
    });
  };

  return (
    <div className="flex min-h-svh items-center justify-center bg-background px-4 py-8">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>{mode === "signin" ? "Sign In" : "Sign Up"}</CardTitle>
          <CardDescription>
            {mode === "signin"
              ? "Enter your email and password to continue."
              : "Create an account with your email and password."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            className="flex flex-col gap-5"
            onSubmit={(e) => {
              e.preventDefault();
            }}
          >
            <FieldGroup className="gap-4">
              <Field className="gap-2">
                <FieldLabel htmlFor={`${mode}-email`}>Email</FieldLabel>
                <Input
                  id={`${mode}-email`}
                  type="email"
                  inputSize="lg"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
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
                    className="pr-12"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword((value) => !value)}
                    className="absolute inset-y-0 right-0 flex items-center px-4 text-muted-foreground transition-colors hover:text-foreground"
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
              onClick={() => {
                if (mode === "signin") {
                  handleSignin({ email, password });
                } else {
                  handleSignup({ email, password });
                }
              }}
              type="submit"
              size="lg"
              className="w-full cursor-pointer transition-colors hover:bg-primary/90 active:scale-[0.98]"
            >
              {mode === "signin" ? "Sign In" : "Sign Up"}
            </Button>

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
                    className="cursor-pointer p-4 font-medium text-foreground underline-offset-4 transition-colors hover:text-primary"
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
