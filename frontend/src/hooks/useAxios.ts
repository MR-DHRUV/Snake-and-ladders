import { useNavigate } from "react-router";
import { toast } from "sonner";
import axios from "axios";
import { handleUnauthorized } from "@/lib/utils";

// Custom hook for Axios instance
const useAxios = () => {
    const navigate = useNavigate();

    const axiosInstance = axios.create({
        baseURL: (window as any).__ENV__?.VITE_API_URL || import.meta.env.VITE_API_URL,
        withCredentials: true,
    });

    // Set up response interceptor
    axiosInstance.interceptors.response.use(
        (response) => response,
        (error) => {
            if (error.response && error.response.status == 401) {
                handleUnauthorized(navigate);
            } else {
                toast.error(error.response?.data?.message || "An error occurred. Please try again.");
            }
            return Promise.reject(error);
        }
    );

    return axiosInstance;
};

export default useAxios;
