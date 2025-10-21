import type { User, UserResponse } from '@/types/user';
import { useQuery } from '@tanstack/react-query';
import useAxios from '@/hooks/useAxios';
import type { AxiosInstance } from 'axios';

const fetchUser = async (axiosInstance: AxiosInstance): Promise<User> => {
    const response = await axiosInstance.get('/auth/user');
    const data: UserResponse = await response.data.data;
    return data.user;
}

export function useUser() {
    const axiosInstance = useAxios();

    return useQuery<User>({
        queryKey: ['user'],
        queryFn: () => fetchUser(axiosInstance),
        staleTime: 0,
        retry: 2,
    });
}