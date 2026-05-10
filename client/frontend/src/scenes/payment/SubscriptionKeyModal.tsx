import { CheckIcon, CopyIcon, XIcon } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";

interface SubscriptionKeyModalProps {
  subscriptionKey: string;
  onClose: () => void;
}

export function SubscriptionKeyModal({
  subscriptionKey,
  onClose,
}: SubscriptionKeyModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="relative w-full max-w-lg rounded-2xl bg-white p-6 shadow-[0_18px_40px_rgba(15,23,42,0.2)]">
        <button
          type="button"
          onClick={onClose}
          className="absolute right-4 top-4 rounded-lg p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <XIcon className="size-4" />
        </button>

        <div className="flex flex-col items-center gap-4 text-center">
          <div className="flex size-12 items-center justify-center rounded-full bg-green-50">
            <CheckIcon className="size-6 text-green-600" />
          </div>

          <h2 className="text-xl font-semibold tracking-tight text-foreground">
            Payment Successful
          </h2>

          <p className="text-sm text-muted-foreground">
            Your subscription has been activated. Save your subscription key
            below.
          </p>

          <div className="w-full space-y-3">
            <div className="flex items-center gap-2 rounded-xl border border-black/5 bg-muted/50 px-4 py-3">
              <code className="flex-1 text-sm font-mono text-foreground">
                {subscriptionKey}
              </code>
            </div>
          </div>

          <Button
            type="button"
            onClick={onClose}
            className="bg-black px-8 text-white shadow-none hover:bg-zinc-800"
          >
            Activate License
          </Button>
        </div>
      </div>
    </div>
  );
}
