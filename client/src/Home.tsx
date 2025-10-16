import { useEffect } from "react";
import { Link } from "react-router-dom";

export default function Home() {
    useEffect(() => {
        document.title = "GHost";
    }, []);
    
    return (
        <div>
            
        </div>
    );
}
