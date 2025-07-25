import { BrowserRouter as Router, Routes } from "react-router-dom";
import { HeroUIProvider } from "@heroui/system";

function App() {
  return (
    <HeroUIProvider>
      <Router>
        <Routes />
      </Router>
    </HeroUIProvider>
  );
}

export default App;
