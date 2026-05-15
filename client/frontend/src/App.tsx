import { AuthScene } from "@/scenes/auth/AuthScene";
import {
  SubscriptionScene,
  type Plan,
} from "@/scenes/subscription/SubscriptionScene";
import { PaymentScene } from "@/scenes/payment/PaymentScene";
import { useEffect, useState } from "react";
import { GetState } from "../wailsjs/go/main/App";
import { state } from "../wailsjs/go/models";

function App() {
  const [page, setPage] = useState<"auth" | "app" | "payment" | "subscription">(
    "auth",
  );
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);

  const navigate = (p: "auth" | "app" | "payment" | "subscription") =>
    setPage(p);

  useEffect(() => {
    const getStateAsync = async () => {
      const state = await GetState();
      if (!state) {
        navigate("auth");
      }

      if (state.validLicense) {
        console.log("we have a valid license", state.validLicense);
        navigate("app");
      } else if (state.authToken.length > 0 && !state.validLicense) {
        console.log("we have a valid license", state.validLicense);
        navigate("subscription");
      } else {
        navigate("auth");
      }
    };

    getStateAsync();
  }, []);

  switch (page) {
    case "auth":
      return <AuthScene handleNav={(p) => navigate(p)} />;
    case "app":
      return <>main app</>;
    case "subscription":
      return (
        <SubscriptionScene
          onSelectPlan={(plan) => {
            setSelectedPlan(plan);
            navigate("payment");
          }}
        />
      );
    case "payment":
      if (!selectedPlan) {
        return (
          <SubscriptionScene
            onSelectPlan={(plan) => {
              setSelectedPlan(plan);
              navigate("payment");
            }}
          />
        );
      }
      return (
        <PaymentScene plan={selectedPlan} onBack={() => navigate("app")} />
      );
    default:
      return <AuthScene handleNav={(p) => navigate(p)} />;
  }
}

export default App;
