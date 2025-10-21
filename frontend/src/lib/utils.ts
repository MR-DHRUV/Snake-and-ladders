import { clsx, type ClassValue } from "clsx"
import type { useNavigate } from "react-router";
import { toast } from "sonner";
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export const handleUnauthorized = (navigate: ReturnType<typeof useNavigate>) => {
    toast.info("Session expired, please log in again.");
    localStorage.clear();
    navigate("/login");
}

export const handleInGameError = ({
    message,
    description,
    duration = 3000,
    position = "bottom-center",
}: {
    message: string;
    description?: string;
    duration?: number;
    position?: "top-left" | "top-center" | "top-right" | "bottom-left" | "bottom-center" | "bottom-right";
}) => {
    toast.error(message, {
        description: description,
        duration: duration,
        position: position,
    });
}

export const handleSuccessEvent = ({
    message,
    description,
    duration = 3000,
    position = "bottom-center",
}: {
    message: string;
    description?: string;
    duration?: number;
    position?: "top-left" | "top-center" | "top-right" | "bottom-left" | "bottom-center" | "bottom-right";
}) => {
    toast.success(message, {
        description: description,
        duration: duration,
        position: position,
    });
}