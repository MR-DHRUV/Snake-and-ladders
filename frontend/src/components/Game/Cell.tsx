import React from 'react';
import { cn } from '@/lib/utils';
import type { BoardCell } from '@/types/game';

// Define the types
interface CellProps {
    cell: BoardCell;
    boardSize: number;
    cellSize: number;
}

// Helper function to convert position to SVG coordinates
const positionToCoordinates = (
    pos: number,
    boardSize: number,
    cellSize: number
): { x: number; y: number } => {
    const row = Math.floor((boardSize * boardSize - pos) / boardSize);
    const isEvenRow = row % 2 === 0;

    let col;
    if (isEvenRow) {
        col = boardSize - 1 - ((boardSize * boardSize - pos) % boardSize);
    } else {
        col = (boardSize * boardSize - pos) % boardSize;
    }

    return {
        x: col * cellSize + cellSize / 2,
        y: row * cellSize + cellSize / 2
    };
};

const Cell: React.FC<CellProps> = ({ cell, boardSize, cellSize }) => {

    // Only render snake or ladder visuals if this is the starting position
    const renderSnakeLadder = cell.pos !== cell.next_pos;
    const isLadder = cell.next_pos > cell.pos;
    const isSnake = cell.next_pos < cell.pos;

    return (
        <div
            className={cn(
                "w-18 h-18 flex items-center justify-center border-2 text-sm font-semibold relative",
                cell.pos === cell.next_pos
                    ? "bg-white"
                    : cell.next_pos > cell.pos
                        ? "bg-green-400" // Ladder
                        : "bg-red-400" // Snake
            )}
        >
            {cell.next_pos > cell.pos ? (
                <span className="text-green-800">{cell.pos} → {cell.next_pos}</span>
            ) : cell.next_pos < cell.pos ? (
                <span className="text-red-800">{cell.pos} → {cell.next_pos}</span>
            ) : (
                <span>{cell.pos}</span>
            )}

            {renderSnakeLadder && (
                <div className="absolute top-0 left-0 pointer-events-none">
                    {isLadder && <LadderVisual start={cell.pos} end={cell.next_pos} boardSize={boardSize} cellSize={cellSize} />}
                    {isSnake && <SnakeVisual start={cell.pos} end={cell.next_pos} boardSize={boardSize} cellSize={cellSize} />}
                </div>
            )}
        </div>
    );
};

interface VisualProps {
    start: number;
    end: number;
    boardSize: number;
    cellSize: number;
}

const LadderVisual: React.FC<VisualProps> = ({ start, end, boardSize, cellSize }) => {
    const startCoords = positionToCoordinates(start, boardSize, cellSize);
    const endCoords = positionToCoordinates(end, boardSize, cellSize);

    // Calculate distance and direction
    const deltaX = endCoords.x - startCoords.x;
    const deltaY = endCoords.y - startCoords.y;
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);

    // Ladder width
    const ladderWidth = 20;
    const numRungs = Math.max(2, Math.floor(distance / 40));

    return (
        <svg
            className="absolute z-10"
            width={boardSize * cellSize}
            height={boardSize * cellSize}
            style={{ pointerEvents: 'none' }}
        >
            {/* Left rail */}
            <line
                x1={startCoords.x - ladderWidth / 2}
                y1={startCoords.y}
                x2={endCoords.x - ladderWidth / 2}
                y2={endCoords.y}
                stroke="#8B4513"
                strokeWidth={3}
            />

            {/* Right rail */}
            <line
                x1={startCoords.x + ladderWidth / 2}
                y1={startCoords.y}
                x2={endCoords.x + ladderWidth / 2}
                y2={endCoords.y}
                stroke="#8B4513"
                strokeWidth={3}
            />

            {/* Rungs */}
            {Array.from({ length: numRungs }).map((_, i) => {
                const t = i / (numRungs - 1);
                const x1 = startCoords.x - ladderWidth / 2 + deltaX * t;
                const y1 = startCoords.y + deltaY * t;
                const x2 = x1 + ladderWidth;

                return (
                    <line
                        key={`rung-${i}`}
                        x1={x1}
                        y1={y1}
                        x2={x2}
                        y2={y1}
                        stroke="#8B4513"
                        strokeWidth={2}
                    />
                );
            })}
        </svg>
    );
};

const SnakeVisual: React.FC<VisualProps> = ({ start, end, boardSize, cellSize }) => {
    const startCoords = positionToCoordinates(start, boardSize, cellSize);
    const endCoords = positionToCoordinates(end, boardSize, cellSize);

    // Create a curved snake path
    const midX = (startCoords.x + endCoords.x) / 2;
    const midY = (startCoords.y + endCoords.y) / 2;
    const controlX1 = startCoords.x + (midX - startCoords.x) * 0.5 + 30;
    const controlY1 = startCoords.y + (midY - startCoords.y) * 0.5;
    const controlX2 = midX + (endCoords.x - midX) * 0.5 - 30;
    const controlY2 = midY + (endCoords.y - midY) * 0.5;

    const pathD = `
    M ${startCoords.x} ${startCoords.y}
    C ${controlX1} ${controlY1}, ${controlX2} ${controlY2}, ${endCoords.x} ${endCoords.y}
  `;

    return (
        <svg
            className="absolute z-10"
            width={boardSize * cellSize}
            height={boardSize * cellSize}
            style={{ pointerEvents: 'none' }}
        >
            <path
                d={pathD}
                fill="none"
                stroke="darkgreen"
                strokeWidth={8}
                strokeLinecap="round"
            />

            {/* Snake head */}
            <circle
                cx={endCoords.x}
                cy={endCoords.y}
                r={10}
                fill="darkgreen"
            />

            {/* Snake eyes */}
            <circle cx={endCoords.x - 3} cy={endCoords.y - 3} r={2} fill="white" />
            <circle cx={endCoords.x + 3} cy={endCoords.y - 3} r={2} fill="white" />
        </svg>
    );
};

export default Cell;