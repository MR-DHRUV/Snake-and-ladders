import type { AxiosInstance } from "axios";

export async function createGame({
    max_players,
    max_winners,
    dice_count,
    axiosInsatnce,
}: {
    max_players: number;
    max_winners: number;
    dice_count: number;
    axiosInsatnce: AxiosInstance;
}) {
    return axiosInsatnce.post('/game', {
        max_players,
        max_winners,
        dice_count
    });
}