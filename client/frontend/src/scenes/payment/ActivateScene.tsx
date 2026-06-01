import { Key } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ActivateLicense } from "../../../wailsjs/go/license/License";

interface ActivateSceneProps {
  onActivated: () => void;
  onBack: () => void;
}

export function ActivateScene({ onActivated, onBack }: ActivateSceneProps) {
  const [key, setKey] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleActivate = async () => {
    if (!key.trim()) {
      setError("Please enter a license key");
      return;
    }

    setLoading(true);
    setError("");

    try {
      await ActivateLicense(key.trim());
      onActivated();
    } catch (err) {
      setError(String(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-svh bg-[#f6f6f4] px-4 py-8">
      <div className="mx-auto flex w-full max-w-md flex-col gap-6">
        <button
          type="button"
          onClick={onBack}
          className="inline-flex w-fit items-center gap-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
        >
          Back
        </button>

        <div className="text-center">
          <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-green-50">
            <Key className="size-6 text-green-600" />
          </div>

          <h1 className="mt-4 text-2xl font-semibold tracking-tight text-foreground">
            Activate License
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Enter your license key to activate the application.
          </p>
        </div>

        <div className="flex flex-col gap-4">
          <Input
            inputSize="lg"
            className="h-12"
            placeholder="Paste your license key here"
            value={key}
            onChange={(e) => {
              setKey(e.target.value);
              setError("");
            }}
          />

          {error && <p className="text-sm text-red-500">{error}</p>}

          <Button
            type="button"
            className="h-12 w-full rounded-xl bg-black text-white shadow-none hover:bg-zinc-800"
            disabled={loading}
            onClick={handleActivate}
          >
            {loading ? "Activating..." : "Activate"}
          </Button>
        </div>
      </div>
    </div>
  );
}
