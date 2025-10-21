import {
    BrowserRouter,
    Routes,
    Route,
} from 'react-router';
import { Toaster } from './components/ui/sonner';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Auth } from './components/Auth/Auth';
import AuthCallback from './components/Auth/AuthCallback';
import Home from './components/Home/Home';
import ErrorBoundary from './ErrorBoundary';
import Game from './components/Game/Game';
import Layout from './Layout';
import NotFound from './NotFound';

const queryClient = new QueryClient()

function App() {
    const googleCallbackURL: string = (window as any).__ENV__?.VITE_GOOGLE_REDIRECT_URI || import.meta.env.VITE_GOOGLE_REDIRECT_URI;

    return (
        <ErrorBoundary>
            <QueryClientProvider client={queryClient}>
                <Toaster richColors position="top-right" duration={4000} />
                <BrowserRouter>
                    <Routes>
                        <Route path="/" element={<Layout />}>
                            <Route index element={<Home />} />
                            <Route path="/game/:id" element={<Game />} />
                        </Route>
                        <Route path="/login" element={<Auth />} />
                        <Route path={googleCallbackURL} element={<AuthCallback />} />
                        <Route path="*" element={<NotFound />} />
                    </Routes>
                </BrowserRouter>
            </QueryClientProvider>
        </ErrorBoundary>
    )
}

export default App
