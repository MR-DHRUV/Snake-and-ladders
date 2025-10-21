/* eslint-disable @typescript-eslint/no-unused-vars */
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import PaginationComponent from "../ui/PaginationComponent";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { usePastGames } from "@/hooks/usePastGames";
import { Spinner } from "../ui/spinner";
import config from "@/config/config";
import React from "react";

const PastGames = () => {

    const [page, setPage] = React.useState(1);

    const {
        data,
        isLoading,
        isError,
        error,
    } = usePastGames(page);

    if (isLoading) {
        return (
            <div className='mt-30 flex justify-center'>
                <Spinner size="large" show={isLoading} />
            </div>
        )
    }

    if (isError && error) {
        return (
            <div className='mt-10 flex justify-center'>
                <p className="text-red-500 text-center">Error fetching past games: {error.message}.</p>
            </div>
        )
    }

    if ((!data || !data.games || !data.total) || (data && data.games.length === 0)) {
        return (
            <div className='mt-10 flex justify-center'>
                <p className="text-muted-foreground text-center">
                    Looks like you haven't played any games yet! Get started and have some fun!
                </p>
            </div>
        )
    }

    return (
        <div className="flex flex-col gap-5 max-w-[1400px] mx-auto mt-20">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="font-semibold ">Game</TableHead>
                        <TableHead className="font-semibold ">Date</TableHead>
                        <TableHead className="font-semibold ">Status</TableHead>
                        <TableHead className="font-semibold ">Players</TableHead>
                        <TableHead className="font-semibold text-right">Winners</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {data.games.map((game) => (
                        <TableRow key={game._id}>
                            <TableCell>{game.creator.name + "'s game"}</TableCell>
                            <TableCell>{new Date(game.date).toLocaleDateString()}</TableCell>
                            <TableCell>{game.status === 2 ? 'Finished' : 'Abondaned'}</TableCell>
                            <TableCell>
                                <div className="*:data-[slot=avatar]:ring-background flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:grayscale">
                                    {game.players.map((player) => (
                                        <Avatar key={player._id}>
                                            <AvatarImage src={player.picture} alt={player.name} />
                                            <AvatarFallback>{player.name.charAt(0).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                    ))}
                                </div>
                            </TableCell>
                            <TableCell className="text-right">
                                {game.winners.length > 0 ? (
                                    <div className="*:data-[slot=avatar]:ring-background flex -space-x-2 *:data-[slot=avatar]:ring-2 *:data-[slot=avatar]:grayscale flex-row justify-end">
                                        {game.winners.map((winner) => (
                                            <Avatar key={winner._id}>
                                                <AvatarImage src={winner.picture} alt={winner.name} />
                                                <AvatarFallback>{winner.name.charAt(0).toUpperCase()}</AvatarFallback>
                                            </Avatar>
                                        ))}
                                    </div>
                                ) : 'None'}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            <PaginationComponent
                currentPage={page}
                setCurrentPage={setPage}
                totalPages={Math.ceil(data.total / config.limit)}
            />
        </div>
    )
}

export default PastGames

/**
 export type PastGame = {
     _id: string;
     date: string;
     creator: User;
     players: User[];
     winners: User[];
 }
 
 export type PastGamesResponse = {
     games: PastGame[];
     total: number;
 }
 */