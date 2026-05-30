import { ShoppingCart, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { OpenCheckoutBrowser } from "../../../wailsjs/go/payment/Payment";

interface CheckoutReadyModalProps {
  url: string;
  onClose: () => void;
}

export function CheckoutReadyModal({ onClose, url }: CheckoutReadyModalProps) {
  return (
    <div
      className="absolute flex justify-center items-center w-full h-full top-0 left-0 bg-[rgba(0,0,0,0.6)] backdrop-blur-md"
      onClick={onClose}
    >
      <div className="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-[0_18px_40px_rgba(15,23,42,0.2)]">
        <button
          type="button"
          onClick={onClose}
          className="absolute right-4 top-4 rounded-lg p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <XIcon className="size-4" />
        </button>

        <div className="flex flex-col items-center gap-4 text-center">
          <div className="flex size-12 items-center justify-center rounded-full bg-green-50">
            <ShoppingCart className="size-6 text-green-600" />
          </div>

          <h2 className="text-xl font-semibold tracking-tight text-foreground">
            Checkout Ready
          </h2>

          <Button
            type="button"
            className="mt-2 w-full bg-black text-white shadow-none hover:bg-zinc-800"
            onClick={() => {
              OpenCheckoutBrowser(url);
            }}
          >
            Proceed to Checkout
          </Button>
        </div>
      </div>
    </div>
  );
}
