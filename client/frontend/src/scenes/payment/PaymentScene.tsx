import { ArrowLeftIcon, CheckIcon } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ConfirmPaymentModal } from "./ConfirmPaymentModal";
import { SubscriptionKeyModal } from "./SubscriptionKeyModal";
import type { Plan } from "@/scenes/subscription/SubscriptionScene";
import {
  GenerateLicense,
  ActivateLicense,
} from "../../../wailsjs/go/license/License";

import { CreateCheckout } from "../../../wailsjs/go/payment/Payment";

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-black/5 pb-3 last:border-b-0 last:pb-0">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium text-foreground">{value}</span>
    </div>
  );
}

function SummaryRow({
  label,
  value,
  strong = false,
}: {
  label: string;
  value: string;
  strong?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-4 text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span
        className={strong ? "font-semibold text-foreground" : "text-foreground"}
      >
        {value}
      </span>
    </div>
  );
}

export function PaymentScene({
  plan,
  onBack,
}: {
  plan: Plan;
  onBack: () => void;
}) {
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showKeyModal, setShowKeyModal] = useState(false);
  const [subscriptionKey, setSubscriptionKey] = useState("");

  const handlePayClick = () => {
    setShowConfirmModal(true);
  };

  const handleConfirmPayment = () => {
    CreateCheckout().then((res) => {
      console.log("payment started", res);
    });
    // GenerateLicense({ duration: plan.duration })
    //   .then((res) => {
    //     setShowConfirmModal(false);
    //     setShowKeyModal(true);
    //     setSubscriptionKey(res.Key);
    //     console.log("response", res);
    //   })
    //   .catch((err) => {
    //     console.log("error", err);
    //   });
  };

  const handleActivateLicense = () => {
    ActivateLicense(subscriptionKey)
      .then((res) => {
        setShowKeyModal(false);
        onBack();
      })
      .catch((err) => {
        console.log("error", err);
      });
  };

  return (
    <div className="min-h-svh bg-[#f6f6f4] px-4 py-8">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
        <button
          type="button"
          onClick={onBack}
          className="inline-flex w-fit items-center gap-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
        >
          <ArrowLeftIcon className="size-4" />
          Back to plans
        </button>

        <div className="max-w-2xl">
          <p className="text-xs font-medium uppercase tracking-[0.22em] text-muted-foreground">
            Payment simulation
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-foreground sm:text-4xl">
            Confirm your {plan.name} plan
          </h1>
          <p className="mt-2 text-base text-muted-foreground">
            This is a simulation checkout page using the plan you selected.
          </p>
        </div>

        <div className="grid gap-4 lg:grid-cols-[1.15fr_0.85fr]">
          <div className="flex flex-col gap-4">
            <div className="rounded-2xl border border-black/5 bg-white p-5 shadow-[0_18px_40px_rgba(15,23,42,0.08)]">
              <p className="text-xs font-medium uppercase tracking-[0.2em] text-muted-foreground">
                Selected plan
              </p>
              <div className="mt-3 flex items-end gap-2">
                <span className="text-4xl font-semibold tracking-tight">
                  {plan.price}
                </span>
                <span className="pb-1 text-sm text-muted-foreground">
                  {plan.period}
                </span>
              </div>
              <p className="mt-2 text-sm text-muted-foreground">
                {plan.description}
              </p>
            </div>

            <div className="rounded-2xl border border-black/5 bg-white p-5 shadow-[0_18px_40px_rgba(15,23,42,0.08)]">
              <p className="text-sm font-medium text-foreground">
                Default card information
              </p>
              <div className="mt-4 space-y-3 text-sm text-muted-foreground">
                <InfoRow label="Card holder" value="John Doe" />
                <InfoRow label="Card number" value="4242 4242 4242 4242" />
                <InfoRow label="Expiry" value="08 / 28" />
                <InfoRow label="CVC" value="123" />
              </div>
            </div>
          </div>

          <div className="flex flex-col justify-between rounded-2xl border border-black/5 bg-white p-5 shadow-[0_18px_40px_rgba(15,23,42,0.08)]">
            <div className="space-y-4">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-sm font-medium text-foreground">
                    Order summary
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Review the selected subscription
                  </p>
                </div>
                {plan.featured ? (
                  <span className="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-700">
                    Most popular
                  </span>
                ) : null}
              </div>

              <Separator className="bg-black/5" />

              <SummaryRow label="Plan" value={plan.name} />
              <SummaryRow label="Billing" value={plan.period} />
              <SummaryRow label="Amount due" value={plan.price} strong />
            </div>

            <div className="mt-6 space-y-3">
              <Button
                className="h-12 w-full rounded-xl bg-black text-white shadow-none transition-colors hover:bg-zinc-800"
                onClick={handlePayClick}
              >
                Pay {plan.price}
              </Button>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <CheckIcon className="size-3.5" />
                Simulation only, no real payment will be taken.
              </div>
            </div>
          </div>
        </div>
      </div>

      {showConfirmModal && (
        <ConfirmPaymentModal
          plan={plan}
          onConfirm={handleConfirmPayment}
          onCancel={() => setShowConfirmModal(false)}
        />
      )}

      {showKeyModal && (
        <SubscriptionKeyModal
          subscriptionKey={subscriptionKey}
          onClose={handleActivateLicense}
        />
      )}
    </div>
  );
}
