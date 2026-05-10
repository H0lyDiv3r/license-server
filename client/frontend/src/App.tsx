import { AuthScene } from "@/scenes/auth/AuthScene";
import {
  SubscriptionScene,
  type Plan,
} from "@/scenes/subscription/SubscriptionScene";
import { PaymentScene } from "@/scenes/payment/PaymentScene";
import { useEffect, useState } from "react";
import { IsLoggedIn } from "../wailsjs/go/main/App";

function App() {
  const [page, setPage] = useState<"auth" | "app" | "payment">("auth");
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);

  const navigate = (p: "auth" | "app" | "payment") => setPage(p);

  useEffect(() => {
    IsLoggedIn().then((loggedIn) => {
      if (loggedIn) {
        navigate("app");
      } else {
        navigate("auth");
      }
    });
  }, []);

  switch (page) {
    case "auth":
      return <AuthScene handleNav={(p) => navigate(p)} />;
    case "app":
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
        <PaymentScene
          plan={selectedPlan}
          onBack={() => navigate("app")}
        />
      );
    default:
      return <AuthScene handleNav={(p) => navigate(p)} />;
  }
}

export default App;
