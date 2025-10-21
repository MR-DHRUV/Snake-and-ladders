/*
    For dynamic environment variable injection during deployment.
*/
window.__ENV__ = {
    VITE_GOOGLE_CLIENT_ID: "$VITE_GOOGLE_CLIENT_ID",
    VITE_GOOGLE_REDIRECT_URI: "$VITE_GOOGLE_REDIRECT_URI",
    VITE_API_URL: "$VITE_API_URL",
    VITE_WS_URL: "$VITE_WS_URL"
};
