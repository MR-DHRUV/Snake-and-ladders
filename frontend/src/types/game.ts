import type { ChatMessage } from "./chat";
import { type User } from "./user";

export type PastGame = {
    _id: string;
    date: string;
    creator: User;
    status: number;
    players: User[];
    winners: User[];
}

export type PastGamesResponse = {
    games: PastGame[];
    total: number;
}

export type BoardCell = {
    pos: number;
    next_pos: number;
}

export type Board = {
    cells: BoardCell[];
    total_cells: number;
}

export type LastTurn = {
    player_id: string;
    dice_rolls: number[];
    previous_position: number;
    next_position: number;
}

export type Player = {
    _id: string;
    name: string;
    picture: string;
    position: number;
}

export type Game = {
    _id: string;
    date: string;
    creator: User;
    total_cells: number;
    max_snakes: number;
    max_ladders: number;
    max_players: number;
    max_winners: number;
    dice_count: number;
    game_board: Board;
    status: number;
    current_turn: number;
    players: Player[];
    winners: Player[];
    last_turn?: LastTurn;
}

export type GameState = {
    _id: string;
    current_turn: string;
    players: Player[];
    winners: Player[];
    status: number;
    game_board: Board;
    last_turn?: LastTurn | null;
    messages: ChatMessage[];
    creator: User;
    dice_count: number;
}

export type GameErrorResponse = {
    status: number;
    message: string;
}