import { AuthScene } from "@/scenes/auth/AuthScene";
import { useEffect, useState } from "react";
import { IsLoggedIn } from "../wailsjs/go/main/App";

function App() {
  const [page, setPage] = useState<string>("auth");

  const navigate = (p: string) => setPage(p);

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
      return <>hello this is app</>;
    default:
      return <AuthScene handleNav={(p) => navigate(p)} />;
  }
}

export default App;
