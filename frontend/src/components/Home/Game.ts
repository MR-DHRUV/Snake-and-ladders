import type { AxiosInstance } from "axios";

export async function createGame({
    total_cells,
    max_snakes,
    max_ladders,
    max_players,
    max_winners,
    dice_count,
    axiosInsatnce,
}: {
    total_cells: number;
    max_snakes: number;
    max_ladders: number;
    max_players: number;
    max_winners: number;
    dice_count: number;
    axiosInsatnce: AxiosInstance;
}) {
    return axiosInsatnce.post('/game/create', {
        total_cells,
        max_snakes,
        max_ladders,
        max_players,
        max_winners,
        dice_count
    });
}

export async function joinGame({
    gameId,
    axiosInstance
}: {
    gameId: string;
    axiosInstance: AxiosInstance;
}) {
    return axiosInstance.post(`/game/join`, {
        game_id: gameId
    });
}

export async function startGame({
    gameId,
    axiosInstance
}: {
    gameId: string;
    axiosInstance: AxiosInstance;
}) {
    return axiosInstance.post(`/game/start`, {
        game_id: gameId
    });
}