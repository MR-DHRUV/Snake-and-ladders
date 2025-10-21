import axios from "axios";

async function generatePKCE(): Promise<{ code_verifier: string, code_challenge: string }> {
    const encoder = new TextEncoder();
    const verifier = Array.from(crypto.getRandomValues(new Uint8Array(32)))
        .map(x => ('0' + x.toString(16)).slice(-2))
        .join('');

    const data = encoder.encode(verifier);
    const digest = await crypto.subtle.digest('SHA-256', data);
    const base64url = btoa(String.fromCharCode(...new Uint8Array(digest)))
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=+$/, '');

    return { code_verifier: verifier, code_challenge: base64url };
}

export async function initGoogleAuth(): Promise<void> {
    const { code_verifier, code_challenge } = await generatePKCE();
    localStorage.setItem("code_verifier", code_verifier);

    const params = new URLSearchParams({
        client_id: (window as any).__ENV__?.VITE_GOOGLE_CLIENT_ID || import.meta.env.VITE_GOOGLE_CLIENT_ID,
        redirect_uri: (window as any).__ENV__?.VITE_GOOGLE_REDIRECT_URI || import.meta.env.VITE_GOOGLE_REDIRECT_URI,
        response_type: "code",
        scope: "openid email profile",
        code_challenge,
        code_challenge_method: "S256",
    });

    window.location.href = `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
}

export async function handleGoogleAuthCallback(): Promise<void> {
    const params = new URLSearchParams(window.location.search);
    const code = params.get("code");
    const code_verifier = localStorage.getItem("code_verifier");
    const baseUrl = (window as any).__ENV__?.VITE_API_URL || import.meta.env.VITE_API_URL;

    if (!code || !code_verifier) return Promise.reject("Missing code or code_verifier");

    await axios.post(
        baseUrl + "/auth/google",
        {
            code,
            code_verifier,
        },
        {
            withCredentials: true,
        }
    );
};
