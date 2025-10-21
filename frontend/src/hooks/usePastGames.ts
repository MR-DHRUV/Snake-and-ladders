import type { PastGamesResponse } from '@/types/game';
import { useQuery } from '@tanstack/react-query';
import useAxios from '@/hooks/useAxios';
import type { AxiosInstance } from 'axios';

const fetchPastGames = async (
    axiosInstance: AxiosInstance,
    page: number,
    limit: number
) => {
    const response = await axiosInstance.get(`/past-games?page=${page}&limit=${limit}`);
    const data: PastGamesResponse = response.data.data;
    return data;
};

const defaultLimit = 10;

export function usePastGames(page: number = 1, limit: number = defaultLimit) {
    const axiosInstance = useAxios();

    return useQuery<PastGamesResponse>({
        queryKey: ['past-games', page, limit],
        queryFn: () => fetchPastGames(axiosInstance, page, limit),
        staleTime: 1000 * 60 * 20,
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        retry: 2,
        retryDelay: 2000,
    });
}