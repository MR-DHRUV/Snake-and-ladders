import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import PastGames from "./PastGames";
import CreateGameModal from "./CreateGame";
import { useNavigate } from "react-router";
import { handleInGameError } from "@/lib/utils";

export default function Home() {
    const [gameId, setGameId] = useState("");
    const [showCreateGameModal, setShowCreateGameModal] = useState(false);
    const navigate = useNavigate();

    const handleJoinGame = () => {
        if (gameId.trim() === "") {
            handleInGameError({
                message: "Please enter a valid Game ID",
            });
            return;
        }

        navigate(`/game/${gameId}`);
    };

    return (
        <main className="py-12 px-4">
            <div className="flex flex-col md:flex-row items-center gap-4 justify-center">
                <Button className="w-full md:w-auto min-w-[220px]" size={'lg'} onClick={() => { setShowCreateGameModal(true) }} >Create Game</Button>
                <Button
                    className="w-full md:w-auto min-w-[220px]"
                    size={'lg'}
                    variant={'outline'}
                    onClick={handleJoinGame}
                >
                    Join Game
                </Button>
                <Input
                    placeholder="Enter Game ID"
                    className="w-full md:w-100 h-10 text-center md:text-start"
                    value={gameId}
                    onChange={(e) => setGameId(e.target.value)}
                />
            </div>
            <PastGames />
            {showCreateGameModal && (
                <CreateGameModal setOpen={setShowCreateGameModal} open={showCreateGameModal} />
            )}
        </main>
    );
}
