/* eslint-disable @typescript-eslint/no-empty-object-type */
import { Component, type ErrorInfo } from 'react';

interface State {
    hasError: boolean;
}

class ErrorBoundary extends Component<React.PropsWithChildren<{}>, State> {
    state: State = { hasError: false };

    // this lifecycle method is called when an error is thrown in a child component
    // it allows us to catch the error and update the state accordingly
    static getDerivedStateFromError() {
        return { hasError: true };
    }

    // this lifecycle method is called after an error has been thrown
    componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Error caught in Error Boundary:', error, errorInfo);
    }

    render() {
        if (this.state.hasError) {
            return (
                <div className='flex flex-col items-center justify-center h-screen overflow-hidden text-2xl font-semibold'>
                    Something went wrong. Please try again later.
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
