import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./user/LoginForm";
import PublicListingSearch from "./explore/PublicListingSearch";
import WithNavbar from "./components/navigation/WithNavbar";
import { Home } from "lucide-react";
import { UserSession } from "./components/restriction/UserSession";
import AccountPage from "./user/Account";

function App() {
    return (
        <UserSession>
            <div className="absolute inset-0 max-h-screen max-w-screen bg-slate-50 dark:bg-gray-950">
                <Router>
                    <Routes>
                        <Route path="/login" element={<LoginForm />} />
                        <Route path="/explore/listing" element={<WithNavbar><PublicListingSearch /></WithNavbar>} />
                        <Route path="/" element={<WithNavbar><Home /></WithNavbar>} />
                        <Route path="/account" element={<AccountPage></AccountPage>} />
                    </Routes>
                </Router>
            </div>
        </UserSession>
    );
}

export default App;