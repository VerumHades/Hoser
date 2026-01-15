import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginForm from "./user/LoginForm";
import PublicListingSearch from "./explore/PublicListingSearch";
import WithNavbar from "./components/navigation/WithNavbar";
import { UserSession } from "./components/restriction/UserSession";
import AccountPage from "./user/Account";
import { PublicListingView } from "./explore/PublicListingView";
import PublicLandingPage from "./explore/PublicLandingPage";

function App() {
    return (
        <UserSession>
            <div className="absolute inset-0 max-h-screen max-w-screen  bg-slate-50 dark:bg-gray-950">
                <Router>
                    <Routes>
                        <Route path="/login" element={<LoginForm />} />
                        <Route path="/search" element={<WithNavbar><PublicListingSearch /></WithNavbar>} />
                        <Route path="/" element={<WithNavbar><PublicLandingPage/></WithNavbar>} />
                        <Route path="/dashboard/*" element={<AccountPage></AccountPage>} />
                        <Route path="/listing/:id" element={<WithNavbar><PublicListingView></PublicListingView></WithNavbar>} />
                    </Routes>
                </Router>
            </div>
        </UserSession>
    );
}

export default App;