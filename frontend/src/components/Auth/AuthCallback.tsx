import { handleGoogleAuthCallback } from "./GoogleAuth";
import { useNavigate } from "react-router";
import { Spinner } from "../ui/spinner";
import { useEffect } from "react";

export default function AuthCallback() {

    const navigate = useNavigate();
    const handleCallback = async () => {
        console.log("Handling Google Auth Callback...");
        try {
            await handleGoogleAuthCallback();
            navigate("/");
        } catch (error) {
            console.error("Authentication failed:", error);
            navigate("/login");
        }
    };

    useEffect(() => {
        handleCallback();
    }, []);

    return <div className="flex flex-col gap-5 items-center justify-center h-screen p-4">
        <div className="flex flex-row gap-5 items-center justify-center">
            <Spinner size="large" show={true} className="text-primary" />
            <h2 className="text-2xl font-semibold text-primary">
                Authenticating...
            </h2>
        </div>
        <p className="text-muted-foreground text-lg text-center">
            Just a moment while we log you in. If you are not redirected automatically, please wait or try again.
        </p>
    </div>
}
