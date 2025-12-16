/* eslint-disable react-hooks/exhaustive-deps */
import { useQuery } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import type { User } from '@/types/user';
import type { GameState } from '@/types/game';
import { handleUnauthorized, handleInGameError } from '@/lib/utils';
import { useNavigate } from 'react-router';
import { ActionTypes, ResponseTypes } from '@/constants/game';
import type { ChatMessage } from '@/types/chat';


export function useGame(gameId: string) {
    const { data: user } = useQuery<User>({ queryKey: ['user'] });
    const navigate = useNavigate();

    const [gameState, setGameState] = useState<GameState | null>(null);
    const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
    const [ws, setWs] = useState<WebSocket | null>(null);

    const wsBaseUrl =
        (window as any).__ENV__?.VITE_WS_URL || import.meta.env.VITE_WS_URL;

    useEffect(() => {
        if (!user) return;
        let reconnectTimeout: number | null = null;
        let socket: WebSocket | null = null;

        const connect = () => {
            const wsUrl = `${wsBaseUrl}/game?gameId=${gameId}&userId=${user._id}`;
            socket = new WebSocket(wsUrl);
            setWs(socket);

            socket.onopen = () => {
                console.log("WebSocket connected");
                socket!.send(JSON.stringify({ action: "joinGame" }));
            };

            socket.onmessage = (event) => {
                const resp = JSON.parse(event.data);

                if (resp?.status === 500) {
                    handleInGameError({
                        message: `Something went wrong.`,
                        description: `The game may not exist or is currently unavailable.`,
                        duration: 5000,
                    });
                    navigate("/");
                    return;
                }

                switch (resp?.data?.status) {
                    case 401:
                        handleUnauthorized(navigate);
                        return;
                    case 403:
                        handleInGameError({
                            message: `Not your turn`,
                            description: `Wait for your turn.`,
                        });
                        return;
                    case 404:
                        handleInGameError({
                            message: `Game not found`,
                            description: `Check your game ID`,
                        });
                        navigate("/");
                        return;
                    case 500:
                        handleInGameError({
                            message: `Something went wrong.`,
                            description: `It's not you it's us. We're trying to make it alright.`,
                        });
                        return;
                }

                switch (resp.type) {
                    case ResponseTypes.GameState:
                        setGameState(resp.data);
                    setChatMessages((prev) => [...prev, ...resp.data.messages]);
                        break;
                    case ResponseTypes.ChatMessage:
                    setChatMessages((prev) => [...prev, resp.data]);
                        break;
                    case ResponseTypes.NextGame:
                        navigate(`/game/${resp.data.game_id}`);
                        break;
                }
            };

            socket.onerror = () => {
                console.log("WebSocket error — will attempt reconnect");
            };

            socket.onclose = () => {
                console.log("WebSocket closed — reconnecting in 1s...");
                reconnectTimeout = window.setTimeout(() => {
                    connect();
                }, 1000);
            };
        };

        connect();

        return () => {
            if (socket) socket.close();
            if (reconnectTimeout) clearTimeout(reconnectTimeout);
        };
    }, [gameId, user]);

    function sendMessage(message: { action: string;[key: string]: any }) {
        if (!ws || ws.readyState !== WebSocket.OPEN) return;
        ws.send(JSON.stringify(message));
    }

    return {
        gameState,
        chatMessages,
        sendChatMessage: (msg: string) =>
            sendMessage({ action: ActionTypes.ChatMessage, message: msg }),
        startGame: () => sendMessage({ action: ActionTypes.StartGame }),
        restartGame: () => sendMessage({ action: ActionTypes.RestartGame }),
        joinGame: () => sendMessage({ action: ActionTypes.JoinGame }),
        nextTurn: () => sendMessage({ action: ActionTypes.NextTurn }),
    };
}
