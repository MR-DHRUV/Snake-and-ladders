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
    const { data: user } = useQuery<User>({
        queryKey: ['user'],
    });

    const navigate = useNavigate();
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
    const [ws, setWs] = useState<WebSocket | null>(null);
    const wsBaseUrl = (window as any).__ENV__?.VITE_WS_URL || import.meta.env.VITE_WS_URL;

    function sendChatMessage(message: string) {
        sendMessage({
            action: ActionTypes.ChatMessage,
            message: message,
        });
    }

    function startGame() {
        sendMessage({
            action: ActionTypes.StartGame,
        });
    }

    function restartGame() {
        sendMessage({
            action: ActionTypes.RestartGame,
        });
    }

    function joinGame() {
        sendMessage({
            action: ActionTypes.JoinGame,
        });
    }

    function nextTurn() {
        sendMessage({
            action: ActionTypes.NextTurn,
        });
    }

    function sendMessage(message: { action: string;[key: string]: string }) {
        if (!ws) {
            console.log('WebSocket is null or undefined.');
            return;
        }

        if (ws && ws.readyState === WebSocket.OPEN) {
            console.log('Sending message to WebSocket:', message);
            ws.send(JSON.stringify(message));
        }
    }

    useEffect(() => {
        if (!user) return;

        const wsUrl = `${wsBaseUrl}/game?gameId=${gameId}&userId=${user._id}`;

        // Initialize WebSocket connection
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => {
            console.log('WebSocket connection established on', wsUrl);
            socket.send(JSON.stringify({ action: 'joinGame' }));
        };

        // Listen for messages from WebSocket (GameState type messages)
        socket.onmessage = (event) => {
            console.log(event.data)
            const resp = JSON.parse(event.data);
            console.log('received data:', resp);

            if (resp?.status === 500) {
                handleInGameError({
                    message: `Something went wrong.`,
                    description: `The game you are trying to join may not exist or is currently unavailable.`,
                    duration: 5000,
                });
                navigate('/');
                return;
            }

            switch (resp.data?.status) {
                case 401:
                    handleUnauthorized(navigate);
                    return;
                case 403:
                    handleInGameError({
                        message: `It looks like it's not your turn yet.`,
                        description: `Please wait for your turn to play.`,
                    });
                    return;
                case 404:
                    handleInGameError({
                        message: `Game not found`,
                        description: `Please double check the game ID.`,
                    });
                    // TODO: Navigate to a 404 game page or similar
                    navigate('/');
                    return;
                case 500:
                    handleInGameError({
                        message: `Something went wrong.`,
                        description: `It's not you it's us. We're trying to make it alright.`,
                    });
                    return;
                default:
                    break;
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
                default:
                    break;
            }
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        setWs(socket);

        // Cleanup WebSocket connection when component unmounts or hook reruns
        return () => {
            if (socket) {
                socket.close();
            }
        };
    }, [gameId, user]);

    return {
        gameState,
        chatMessages,
        sendChatMessage,
        startGame,
        restartGame,
        joinGame,
        nextTurn,
    }
}