import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Player } from '@/types/game';
import type { User } from '@/types/user';
import { handleSuccessEvent, handleInGameError } from "@/lib/utils";
import { CopyIcon, Dices } from "lucide-react";
import "./animation.css";

export default function Players({
    players,
    user,
    currentTurn
}: {
    players: Player[];
    user?: User;
    currentTurn?: string;
}) {
    function copyGameLink(): void {
        const gameLink = window.location.href;
        navigator.clipboard.writeText(gameLink).then(() => {
            handleSuccessEvent({
                message: 'Game link copied to clipboard!',
            });
        }).catch(() => {
            handleInGameError({
                message: 'Failed to copy game link.',
                description: 'Please try copying the browser url manually.',
            });
        });
    }

    return (
        <div className="flex flex-col gap-4">
            <h2 className="text-2xl font-semibold">Players </h2>
            <p onClick={copyGameLink} className=" cursor-pointer text-blue-500 mt-[-15px] hover:text-black transition-colors">
                <CopyIcon className="inline-block mr-1 mb-1" size={16} />
                Copy your game link to invite others
            </p>
            <div className="*:data-[slot=avatar]:ring-background flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:grayscale flex-col justify-center gap-1">
                {players.map((player) => (
                    <div key={player._id} className="flex flex-row items-center gap-2">
                        <Avatar key={player._id} className='size-12'>
                            <AvatarImage src={player.picture} alt={player.name} />
                            <AvatarFallback>{player.name.charAt(0).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <p className='w-max'>
                            {player._id === user?._id ? 'You' : player.name}
                            {currentTurn === player._id && (
                                <span className="inline-block ml-2 text-amber-400 pulsate-fwd" title="Current Turn">
                                    <Dices />
                                </span>
                            )}
                        </p>
                    </div>
                ))}
            </div>
        </div>
    )
}
