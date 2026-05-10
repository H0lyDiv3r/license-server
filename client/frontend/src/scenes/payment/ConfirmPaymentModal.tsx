import { CheckIcon, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import type { Plan } from "@/scenes/subscription/SubscriptionScene";

interface ConfirmPaymentModalProps {
  plan: Plan;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmPaymentModal({
  plan,
  onConfirm,
  onCancel,
}: ConfirmPaymentModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="relative w-full max-w-lg rounded-2xl bg-white p-6 shadow-[0_18px_40px_rgba(15,23,42,0.2)]">
        <button
          type="button"
          onClick={onCancel}
          className="absolute right-4 top-4 rounded-lg p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <XIcon className="size-4" />
        </button>

        <div className="flex flex-col items-center gap-4 text-center">
          <div className="flex size-12 items-center justify-center rounded-full bg-green-50">
            <CheckIcon className="size-6 text-green-600" />
          </div>

          <h2 className="text-xl font-semibold tracking-tight text-foreground">
            Confirm Payment
          </h2>

          <p className="text-sm text-muted-foreground">
            You are about to be charged for the following subscription.
          </p>

          <div className="w-full space-y-3 rounded-xl border border-black/5 bg-muted/50 p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Plan</span>
              <span className="text-sm font-medium text-foreground">
                {plan.name}
              </span>
            </div>
            <Separator className="bg-black/5" />
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Amount</span>
              <span className="text-sm font-semibold text-foreground">
                {plan.price}
                <span className="text-xs text-muted-foreground">
                  {" "}
                  {plan.period}
                </span>
              </span>
            </div>
          </div>

          <div className="flex w-full gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              className="flex-1"
            >
              Cancel
            </Button>
            <Button
              type="button"
              onClick={onConfirm}
              className="flex-1 bg-black text-white shadow-none hover:bg-zinc-800"
            >
              Confirm and Pay
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
