import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./LoginForm";
import Home from "./Home";
import Navbar from "./Navbar";
import Search from "./Search";
import AccountPage from "./Account";
import DashboardApp from "./DeveloperDashboard";

function App() {
    return (
        <div className="absolute flex flex-col inset-0 justify-start">
            <div className="h-16"></div>
            <Router>
                <Navbar></Navbar>
                <Routes>
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/explore" element={<Search />} />
                    <Route path="/account" element={<AccountPage />} />
                    <Route path="/developer/dashboard" element={<DashboardApp />} />
                    <Route path="/" element={<Home />} />
                </Routes>
            </Router>
        </div>
    );
}

export default App;