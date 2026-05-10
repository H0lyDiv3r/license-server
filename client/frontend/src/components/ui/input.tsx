import * as React from "react"

import { cn } from "@/lib/utils"

function Input({ className, type, inputSize = "default", ...props }: React.ComponentProps<"input"> & { inputSize?: "default" | "lg" }) {
  return (
    <input
      type={type}
      data-slot="input"
      data-size={inputSize}
      className={cn(
        "flex h-8 w-full min-w-0 items-center rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 data-[size=lg]:h-10 data-[size=lg]:text-base data-[size=lg]:[&::-webkit-inner-spin-button,input[type=number]::-webkit-inner-spin-button]:h-10 data-[size=lg]:[&::-webkit-outer-spin-button,input[type=number]::-webkit-outer-spin-button]:h-10",
        className
      )}
      {...props}
    />
  )
}

export { Input }
