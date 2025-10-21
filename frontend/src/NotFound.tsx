import React, { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router";
import { Spinner } from "./components/ui/spinner";

const NotFound: React.FC = () => {
    const navigate = useNavigate();

    useEffect(() => {
        const timer = setTimeout(() => {
            navigate("/");
        }, 1200);
        return () => clearTimeout(timer);
    }, [navigate]);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground">
            <h1 className="text-9xl font-bold mb-4">404</h1>
            <p className="text-2xl mb-6">Page not found.</p>
            <p className="mb-8 text-muted-foregroun flex flex-row gap-2">
                <Spinner size="small" className="inline-block" />
                Redirecting to home...
            </p>
            <Button onClick={() => navigate("/")}>Go to Home</Button>
        </div>
    );
};

export default NotFound;