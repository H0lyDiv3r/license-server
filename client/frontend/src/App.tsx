import "./App.css";
import { DecodeLicense } from "../wailsjs/go/license/License";

function App() {
  const handleDecode = (license: string) => {
    DecodeLicense(license);
  };

  return <div id="App"></div>;
}

export default App;
