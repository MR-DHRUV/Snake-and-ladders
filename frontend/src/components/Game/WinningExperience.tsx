
import { useEffect, useRef, useState } from 'react';
import type { User } from '@/types/user';
import Confetti from 'react-confetti';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useNavigate } from 'react-router';


function WinningExperience({
    winners,
    currentUser,
    creator,
    restartGame
}: {
    winners: User[],
    currentUser?: User,
    creator?: User,
    restartGame: () => void
}) {
    const [showConfetti, setShowConfetti] = useState(true);
    const [animate, setAnimate] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);
    const navigate = useNavigate();

    useEffect(() => {
        setAnimate(true);
        setShowConfetti(true);
        const confettiTimeout = setTimeout(() => setShowConfetti(false), 10000);
        return () => clearTimeout(confettiTimeout);
    }, [winners]);

    if (winners.length === 0) {
        return null;
    }

    return (
        <div
            ref={containerRef}
            className="fixed inset-0 z-[1000] flex flex-col items-center justify-center w-screen h-screen bg-black/30 backdrop-blur-md"
        >
            {showConfetti && (
                <Confetti
                    width={window.innerWidth}
                    height={window.innerHeight}
                    numberOfPieces={350}
                    recycle={false}
                />
            )}
            <div className="flex flex-col items-center gap-8">
                <div className="flex flex-row items-end justify-center gap-8 mb-8">
                    {winners.map((winner) => (
                        <div key={winner._id} className="flex flex-col items-center gap-6">
                            <Avatar
                                key={winner._id}
                                className={`rounded-full shadow-2xl border-4 border-white transition-all duration-1000 ${animate ? 'size-70 text-9xl' : 'size-10'}`}
                                style={{ boxShadow: '0 0 40px 10px #fff8' }}
                            >
                                <AvatarImage src={winner.picture} alt={winner.name} />
                                <AvatarFallback>{winner.name.charAt(0).toUpperCase()}</AvatarFallback>
                            </Avatar>
                            <h2 className="text-white text-5xl font-bold mt-4 drop-shadow-lg">
                                {winner._id === currentUser?._id ? "You have won the game!" : `${winner.name.toLocaleLowerCase().charAt(0).toUpperCase() + winner.name.toLocaleLowerCase().slice(1)} wins!`}
                            </h2>
                        </div>
                    ))}
                </div>
                <div className="flex flex-row gap-6">
                    {currentUser?._id === creator?._id && <Button size="lg" onClick={restartGame}>Restart</Button>}
                    <Button size="lg" variant="secondary" onClick={() => navigate('/')}>Create New Game</Button>
                </div>
            </div>
        </div>
    );
}

export default WinningExperience
