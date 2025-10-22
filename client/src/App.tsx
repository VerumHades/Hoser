import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./user/LoginForm";
import Home from "./user/Home";
import Navbar from "./Navbar";
import AccountPage from "./user/Account";
import DashboardApp from "./developer/DeveloperDashboard";
import DeveloperJoin from "./developer/DeveloperJoin";
import PublicRentalSearch from "./explore/PublicRentalSearch";

function App() {
    return (
        <div className="absolute flex flex-col inset-0 justify-start">
            <div className="h-16"></div>
            <Router>
                <Navbar></Navbar>
                <Routes>
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/explore/listings" element={<PublicRentalSearch/>} />
                    <Route path="/account" element={<AccountPage />} />
                    <Route path="/developer/dashboard" element={<DashboardApp />} />
                    <Route path="/developer/join" element={<DeveloperJoin />} />
                    <Route path="/" element={<Home />} />
                </Routes>
            </Router>
        </div>
    );
}

export default App;