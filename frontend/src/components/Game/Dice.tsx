import type { LastTurn } from '@/types/game';
import "./animation.css";

export default function Dice({
    nextTurn,
    lastTurn,
    currentTurn,
    currentUser,
    gameState,
    diceCount
}: {
    nextTurn: () => void;
    lastTurn?: LastTurn | null;
    currentTurn: string | null;
    currentUser?: string;
    gameState: number;
    diceCount: number;
}) {
    return (
        <button
            onClick={nextTurn}
            className={`flex flex-row items-end cursor-pointer gap-2 ${currentTurn !== currentUser || gameState !== 1 ? 'opacity-50 cursor-not-allowed' : ' dice-animation'}`}
            disabled={currentTurn !== currentUser || gameState != 1}
        >
            {lastTurn?.dice_rolls.map((roll, index) => (
                <img key={index} src={`/dice-${roll}.webp`} alt={`Dice showing ${roll}`} className='h-18 w-18 rounded' />
            ))}
            {!lastTurn || lastTurn.dice_rolls.length === 0 ? (
                Array.from({ length: diceCount }).map((_, index) => (
                    <img key={index} src={`/dice-${index + 1}.webp`} alt={`Dice showing ${index + 1}`} className='h-18 w-18 rounded' />
                ))
            ) : null}
        </button>
    )
}
