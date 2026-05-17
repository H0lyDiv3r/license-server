import { AuthScene } from "@/scenes/auth/AuthScene";
import {
  SubscriptionScene,
  type Plan,
} from "@/scenes/subscription/SubscriptionScene";
import { PaymentScene } from "@/scenes/payment/PaymentScene";
import { useEffect, useState } from "react";
import { GetState, CheckLicense } from "../wailsjs/go/main/App";
import { MainPage } from "./scenes/main/MainPage";

function App() {
  const [page, setPage] = useState<
    "auth" | "app" | "payment" | "subscription" | null
  >(null);
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

  useEffect(() => {
    CheckLicense();
  }, []);

  switch (page) {
    case "auth":
      return <AuthScene handleNav={(p) => navigate(p)} />;
    case "app":
      return <MainPage />;
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
      return <>loading</>;
  }
}

export default App;
