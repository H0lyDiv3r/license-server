import { ArrowRightIcon, CheckIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

export type Plan = {
  name: string;
  price: string;
  period: string;
  description: string;
  features: readonly string[];
  cta: string;
  featured?: boolean;
};

export const plans: readonly Plan[] = [
  {
    name: "1 Month",
    price: "$9",
    period: "/month",
    description: "Flexible access for trying the app or short-term use.",
    features: ["Full access", "Cancel anytime", "Email support"],
    cta: "Choose Monthly",
  },
  {
    name: "6 Months",
    price: "$45",
    period: "/6 months",
    description: "Best balance between commitment and savings.",
    features: ["Full access", "Priority support", "Best value"],
    cta: "Choose 6 Months",
    featured: true,
  },
  {
    name: "1 Year",
    price: "$79",
    period: "/year",
    description: "Best value for long-term access and stability.",
    features: ["Full access", "Priority support", "Lowest monthly cost"],
    cta: "Choose Yearly",
  },
];

export function SubscriptionScene({
  onSelectPlan,
}: {
  onSelectPlan: (plan: Plan) => void;
}) {
  return (
    <div className="min-h-svh bg-[#f6f6f4] px-4 py-8">
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-5">
        <div className="max-w-2xl">
          <p className="text-xs font-medium uppercase tracking-[0.22em] text-muted-foreground">
            Subscription
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-foreground sm:text-4xl">
            Choose your subscription
          </h1>
          <p className="mt-2 text-base text-muted-foreground">
            Pick the plan that matches how long you want access.
          </p>
        </div>

        <div className="flex gap-4 overflow-x-auto pb-2 lg:flex-nowrap">
          {plans.map((plan) => (
            <Card
              key={plan.name}
              className={
                plan.featured
                  ? "group min-h-[520px] min-w-[320px] flex-1 rounded-[22px] border border-black/10 bg-white p-0 shadow-[0_18px_40px_rgba(15,23,42,0.10)] ring-1 ring-black/5 transition-all duration-200 hover:-translate-y-1 hover:shadow-[0_24px_48px_rgba(15,23,42,0.14)]"
                  : "group min-h-[520px] min-w-[320px] flex-1 rounded-[22px] border border-black/5 bg-white p-0 shadow-[0_14px_30px_rgba(15,23,42,0.07)] ring-1 ring-black/5 transition-all duration-200 hover:-translate-y-1 hover:shadow-[0_22px_44px_rgba(15,23,42,0.11)]"
              }
            >
              <CardHeader className="space-y-3 border-b border-black/5 px-5 py-5">
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <CardTitle className="text-3xl font-semibold tracking-tight text-foreground">
                      {plan.name}
                    </CardTitle>
                    <CardDescription className="mt-2 max-w-[24ch] text-sm leading-6 text-muted-foreground">
                      {plan.description}
                    </CardDescription>
                  </div>
                  {plan.featured ? (
                    <span className="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-700">
                      Most popular
                    </span>
                  ) : null}
                </div>

                <div className="flex items-end gap-1 pt-1">
                  <span className="text-5xl font-semibold tracking-tight">
                    {plan.price}
                  </span>
                  <span className="pb-1 text-sm text-muted-foreground">
                    {plan.period}
                  </span>
                </div>
              </CardHeader>

              <CardContent className="flex min-h-[340px] flex-1 flex-col justify-between px-5 py-5">
                <div className="space-y-5">
                  <Separator className="bg-black/5" />

                  <ul className="space-y-3 text-sm text-muted-foreground">
                    {plan.features.map((feature) => (
                      <li key={feature} className="my-7 flex items-center gap-3">
                        <span className="flex size-5 items-center justify-center rounded-full bg-black text-white">
                          <CheckIcon className="size-3.5" />
                        </span>
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <Button
                  type="button"
                  variant={plan.featured ? "default" : "outline"}
                  onClick={() => onSelectPlan(plan)}
                  className={
                    plan.featured
                      ? "mt-8 h-12 w-full rounded-xl border-0 bg-black px-4 text-white shadow-none transition-colors hover:bg-zinc-800"
                      : "mt-8 h-12 w-full rounded-xl border-black/10 bg-white px-4 text-black shadow-none transition-colors hover:bg-black hover:text-white"
                  }
                >
                  <span className="flex w-full items-center justify-between gap-3">
                    <span className="text-sm font-medium">Get started for</span>
                    <span className="flex items-center gap-2 text-sm font-semibold">
                      <span>{plan.price}</span>
                      <ArrowRightIcon className="size-4" />
                    </span>
                  </span>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}
