import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./user/LoginForm";
import Home from "./user/Home";
import Navbar from "./Navbar";
import AccountPage from "./user/Account";
import DashboardApp from "./developer/DeveloperDashboard";
import DeveloperJoin from "./developer/DeveloperJoin";
import PublicListingSearch from "./explore/PublicListingSearch";
import WithNavbar from "./components/prefabs/WithNavbar";

function App() {
    return (
        <div className="absolute inset-0">
            <Router>
                <Routes>
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/explore/listings" element={<WithNavbar><PublicListingSearch /></WithNavbar>} />
                    <Route path="/account" element={<AccountPage />} />
                    <Route path="/developer/dashboard" element={<DashboardApp />} />
                    <Route path="/developer/join" element={<DeveloperJoin />} />
                    <Route path="/" element={<WithNavbar><Home /></WithNavbar>} />
                </Routes>
            </Router>
        </div>
    );
}

export default App;