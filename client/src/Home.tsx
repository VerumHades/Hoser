import { Link } from "react-router-dom";

export default function Home() {
    return (
        <div>
            <h2>Welcome Home!</h2>
            <p>This is the home page.</p>
            <Link to="/login" style={{ marginRight: 10 }}>Login</Link>
        </div>
    );
}
