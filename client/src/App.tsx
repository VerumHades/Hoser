import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./user/LoginForm";
import Home from "./user/Home";
import Navbar from "./Navbar";
import AccountPage from "./user/Account";
import DashboardApp from "./developer/DeveloperDashboard";
import DeveloperJoin from "./developer/DeveloperJoin";
import PublicListingSearch from "./explore/PublicListingSearch";

function App() {
    return (
        <div className="absolute flex flex-col inset-0 justify-start">
            <Router>
                <Navbar></Navbar>
                <div className="absolute inset-0 top-16 flex flex-col items-center">
                <Routes>
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/explore/listings" element={<PublicListingSearch/>} />
                    <Route path="/account" element={<AccountPage />} />
                    <Route path="/developer/dashboard" element={<DashboardApp />} />
                    <Route path="/developer/join" element={<DeveloperJoin />} />
                    <Route path="/" element={<Home />} />
                </Routes>
                </div>

            </Router>
        </div>
    );
}

export default App;