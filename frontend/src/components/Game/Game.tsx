import { useParams } from 'react-router';
import Board from './Board';
import { useGame } from '../../hooks/useGame';
import Players from './Players';
import Dice from './Dice';
import { useNavigate } from "react-router";
import { Button } from '../ui/button';
import { useQuery } from '@tanstack/react-query';
import type { User } from '@/types/user';
import { Spinner } from '../ui/spinner';
import ChatContainer from '../ChatContainer/ChatContainer';
import WinningExperience from './WinningExperience';


export default function Game() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    const {
        gameState,
        chatMessages,
        startGame,
        restartGame,
        nextTurn,
        sendChatMessage
    } = useGame(id || '');
    const { data: user } = useQuery<User>({
        queryKey: ['user'],
    });

    if (!id) {
        navigate('/');

        return (
            <main className="py-12 px-4">
                <h1 className="text-2xl font-bold mb-6">Game ID is missing</h1>
            </main>
        );
    }

    if (!gameState) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Spinner size={'large'} show={true} />
                <p className="text-lg">Loading game...</p>
            </div>
        );
    }

    return (
        <main className="py-4 px-4 flex flex-row justify-around gap-4">
            <WinningExperience winners={gameState?.winners || []} currentUser={user} restartGame={restartGame} creator={gameState?.creator} />
            <div className="flex flex-col justify-between gap-20 min-w-[300px]">
                <Players players={gameState?.players || []} user={user} currentTurn={gameState?.current_turn} />
                <div className="flex flex-row gap-4 items-center justify-start">
                    {gameState?.creator._id === user?._id && gameState?.status === 0 && (
                        <Button onClick={startGame} className='w-max text-xl h-[72px]'>
                            Start Game
                        </Button>
                    )}

                    {gameState?.creator._id !== user?._id && gameState?.status === 0 && (
                        <h4 className='w-max text-xl h-[72px]'>
                            Waiting for the host to start the game...
                        </h4>
                    )}

                    {gameState?.status === 1 &&
                        <Dice nextTurn={nextTurn}
                            lastTurn={gameState?.last_turn}
                            currentTurn={gameState?.current_turn}
                            currentUser={user?._id}
                            gameState={gameState?.status}
                            diceCount={gameState?.dice_count}
                        />}
                </div>
            </div>
            <Board gameState={gameState} />
            <ChatContainer
                chatMessages={chatMessages}
                currentUser={user?._id}
                sendChatMessage={sendChatMessage}
            />
        </main>
    );
}