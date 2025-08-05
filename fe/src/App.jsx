import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { HeroUIProvider } from "@heroui/system";
import LogIn from "./pages/LogIn";

function App() {
  return (
    <HeroUIProvider>
      <Router>
        <Routes>
          <Route path="/" element={<LogIn />} />

          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </Router>
    </HeroUIProvider>
  );
}

export default App;
