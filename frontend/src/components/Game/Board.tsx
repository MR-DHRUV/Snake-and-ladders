import { useEffect, useRef } from "react";
import type { GameState } from "@/types/game";


export default function Board({ gameState }: { gameState: GameState | null }) {

    const canvasRef = useRef<HTMLCanvasElement>(null);

    const boardSize = Math.sqrt(gameState?.game_board.total_cells || 100);
    const cellSize = 70;
    const startingAreaHeight = 70; // Extra space for starting position
    const canvasWidth = cellSize * boardSize;
    const canvasHeight = cellSize * boardSize + startingAreaHeight;

    useEffect(() => {
        if (gameState) {
            const canvas = canvasRef.current;
            if (!canvas) return;

            const ctx = canvas.getContext('2d');
            if (!ctx) return;

            // Clear canvas
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            // Draw the board cells
            drawBoard(ctx);

            // Draw snakes and ladders
            drawSnakesAndLadders(ctx);

            // Draw player positions
            drawPlayerPositions(ctx);
        }
    }, [gameState, gameState?._id]);

    if (!gameState) {
        return null;
    }

    const drawBoard = (ctx: CanvasRenderingContext2D) => {
        ctx.lineWidth = 2;

        for (let row = 0; row < boardSize; row++) {
            for (let col = 0; col < boardSize; col++) {
                const x = col * cellSize;
                const y = (boardSize - row - 1) * cellSize; // Remove offset - board starts at top

                // Determine cell number based on row (alternate left-to-right and right-to-left)
                let cellNum;
                if (row % 2 === 0) {
                    // Left to right
                    cellNum = row * boardSize + col + 1;
                } else {
                    // Right to left
                    cellNum = (row + 1) * boardSize - col;
                }

                // Create vibrant alternating colors with gradients
                // const gradient = ctx.createLinearGradient(x, y, x + cellSize, y + cellSize);

                if ((row + col) % 2 === 0) {
                    // Light vibrant cells
                    // gradient.addColorStop(0, '#ff9999'); // Light pink
                    // gradient.addColorStop(1, '#ff4d4d'); // Peach
                    ctx.fillStyle = '#ffffff';
                } else {
                    // Darker vibrant cells
                    // gradient.addColorStop(0, '#fffc00'); // Light green
                    // gradient.addColorStop(1, '#f3f9a7'); // Sky blue
                    ctx.fillStyle = '#fff8e6';
                }

                // Draw cell background with gradient
                ctx.fillRect(x, y, cellSize, cellSize);

                // Add a colorful border
                ctx.strokeStyle = '#000000';
                ctx.strokeRect(x, y, cellSize, cellSize);

                // Add inner highlight for depth
                ctx.strokeStyle = '#FFFFFF';
                ctx.lineWidth = 1;
                ctx.strokeRect(x + 2, y + 2, cellSize - 4, cellSize - 4);

                // Draw cell number with shadow effect
                // Shadow
                ctx.fillStyle = 'rgba(0, 0, 0, 0.3)';
                ctx.font = 'bold 16px Arial';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillText(cellNum.toString(), x + cellSize / 2 + 1, y + cellSize / 2 + 1);

                // Main text
                ctx.fillStyle = '#2C3E50';
                ctx.fillText(cellNum.toString(), x + cellSize / 2, y + cellSize / 2);

                // Reset line width for next iteration
                ctx.lineWidth = 2;
            }
        }

        // Draw starting area
        drawStartingArea(ctx);
    };

    const drawSnakesAndLadders = (ctx: CanvasRenderingContext2D) => {
        const drawnConnections = new Set<string>();

        gameState.game_board.cells.forEach(cell => {
            if (cell.pos !== cell.next_pos) {
                const connectionId = `${Math.min(cell.pos, cell.next_pos)}-${Math.max(cell.pos, cell.next_pos)}`;
                if (drawnConnections.has(connectionId)) return;
                drawnConnections.add(connectionId);

                const startPos = getCellPosition(cell.pos);
                const endPos = getCellPosition(cell.next_pos);

                if (cell.next_pos > cell.pos) {
                    ctx.strokeStyle = '#4CAF50';
                    drawLadder(ctx, startPos, endPos);
                } else {
                    ctx.strokeStyle = '#F44336';
                    drawSnake(ctx, startPos, endPos);
                }
            }
        });
    };

    const drawLadder = (
        ctx: CanvasRenderingContext2D,
        start: { x: number, y: number },
        end: { x: number, y: number }
    ) => {
        const ladderWidth = 24;
        const sideThickness = 6;

        // Direction vector
        const dx = end.x - start.x;
        const dy = end.y - start.y;
        const distance = Math.sqrt(dx * dx + dy * dy);

        // Normalize
        const ux = dx / distance;
        const uy = dy / distance;

        // Perpendicular vector for ladder width
        const px = -uy * ladderWidth;
        const py = ux * ladderWidth;

        // === Cyan gradient for futuristic ladder ===
        const gradient = ctx.createLinearGradient(start.x, start.y, end.x, end.y);
        gradient.addColorStop(0, "#34adc1"); // Cyan
        gradient.addColorStop(1, "#2b8fa3"); // Deeper teal
        // gradient.addColorStop(1, "#1f6b7a"); // Dark teal

        ctx.lineCap = "round";
        ctx.lineJoin = "round";
        ctx.strokeStyle = gradient;
        ctx.fillStyle = gradient;
        ctx.shadowColor = "#34adc1";
        ctx.shadowBlur = 6;

        // === Sides ===
        ctx.lineWidth = sideThickness;
        ctx.beginPath();
        ctx.moveTo(start.x, start.y);
        ctx.lineTo(end.x, end.y);
        ctx.stroke();

        ctx.beginPath();
        ctx.moveTo(start.x + px, start.y + py);
        ctx.lineTo(end.x + px, end.y + py);
        ctx.stroke();

        // === Rungs ===
        const rungs = Math.floor(distance / 28);
        for (let i = 1; i < rungs; i++) {
            const rx = start.x + i * (dx / rungs);
            const ry = start.y + i * (dy / rungs);

            ctx.lineWidth = sideThickness - 2;
            ctx.beginPath();
            ctx.moveTo(rx, ry);
            ctx.lineTo(rx + px, ry + py);
            ctx.stroke();
        }

        ctx.shadowBlur = 0;

        // === Metallic highlight ===
        const highlight = ctx.createLinearGradient(start.x, start.y, end.x, end.y);
        highlight.addColorStop(0, "rgba(255,255,255,0.5)");
        highlight.addColorStop(0.5, "rgba(255,255,255,0.2)");
        highlight.addColorStop(1, "rgba(255,255,255,0.4)");

        ctx.lineWidth = 2;
        ctx.strokeStyle = highlight;

        // Left highlight
        ctx.beginPath();
        ctx.moveTo(start.x - 2, start.y - 2);
        ctx.lineTo(end.x - 2, end.y - 2);
        ctx.stroke();

        // Right highlight
        ctx.beginPath();
        ctx.moveTo(start.x + px - 2, start.y + py - 2);
        ctx.lineTo(end.x + px - 2, end.y + py - 2);
        ctx.stroke();
    };

    const drawSnake = (
        ctx: CanvasRenderingContext2D,
        start: { x: number, y: number },
        end: { x: number, y: number }
    ) => {
        const snakeWidth = 14;
        const midX = (start.x + end.x) / 2 + 40;

        // Gradient body
        const gradient = ctx.createLinearGradient(start.x, start.y, end.x, end.y);
        gradient.addColorStop(0, '#228B22');   // Forest green
        gradient.addColorStop(0.5, '#32CD32'); // Lime green
        gradient.addColorStop(1, '#006400');   // Dark green

        ctx.lineWidth = snakeWidth;
        ctx.strokeStyle = gradient;
        ctx.shadowColor = '#32CD32';
        ctx.shadowBlur = 12;

        // Snake body
        ctx.beginPath();
        ctx.moveTo(start.x, start.y);
        ctx.bezierCurveTo(midX, start.y, midX, end.y, end.x, end.y);
        ctx.stroke();

        // Reset shadow for details
        ctx.shadowBlur = 0;

        // Scales
        const segments = 12;
        ctx.fillStyle = 'rgba(255,255,255,0.2)';
        for (let i = 1; i < segments; i++) {
            const t = i / segments;
            const x = Math.pow(1 - t, 3) * start.x +
                3 * Math.pow(1 - t, 2) * t * midX +
                3 * (1 - t) * Math.pow(t, 2) * midX +
                Math.pow(t, 3) * end.x;
            const y = Math.pow(1 - t, 3) * start.y +
                3 * Math.pow(1 - t, 2) * t * start.y +
                3 * (1 - t) * Math.pow(t, 2) * end.y +
                Math.pow(t, 3) * end.y;

            ctx.save();
            ctx.translate(x, y);
            ctx.rotate(Math.PI / 4);
            ctx.fillRect(-3, -3, 6, 6);
            ctx.restore();
        }

        // === Snake head (oval) ===
        ctx.beginPath();
        ctx.ellipse(start.x, start.y, snakeWidth * 1.2, snakeWidth, 0, 0, 2 * Math.PI);
        ctx.fillStyle = gradient;
        ctx.fill();
        ctx.strokeStyle = '#004d00';
        ctx.lineWidth = 2;
        ctx.stroke();

        // === Eyes with slit pupils ===
        const eyeOffsetX = 6;
        const eyeOffsetY = -3;

        // Eye whites
        ctx.beginPath();
        ctx.ellipse(start.x - eyeOffsetX, start.y + eyeOffsetY, 3, 4, 0, 0, 2 * Math.PI);
        ctx.ellipse(start.x + eyeOffsetX, start.y + eyeOffsetY, 3, 4, 0, 0, 2 * Math.PI);
        ctx.fillStyle = 'white';
        ctx.fill();

        // Vertical slit pupils
        ctx.beginPath();
        ctx.moveTo(start.x - eyeOffsetX, start.y + eyeOffsetY - 3);
        ctx.lineTo(start.x - eyeOffsetX, start.y + eyeOffsetY + 3);
        ctx.moveTo(start.x + eyeOffsetX, start.y + eyeOffsetY - 3);
        ctx.lineTo(start.x + eyeOffsetX, start.y + eyeOffsetY + 3);
        ctx.strokeStyle = 'black';
        ctx.lineWidth = 1.5;
        ctx.stroke();

        // === Nostrils ===
        ctx.beginPath();
        ctx.arc(start.x - 3, start.y + 5, 1.2, 0, 2 * Math.PI);
        ctx.arc(start.x + 3, start.y + 5, 1.2, 0, 2 * Math.PI);
        ctx.fillStyle = '#111';
        ctx.fill();

        // === Forked tongue (red) ===
        ctx.beginPath();
        ctx.moveTo(start.x, start.y - 5 + snakeWidth);
        ctx.lineTo(start.x, start.y - 5 + snakeWidth + 10);
        ctx.moveTo(start.x, start.y - 5 + snakeWidth + 10);
        ctx.lineTo(start.x - 2, start.y - 5 + snakeWidth + 7);
        ctx.moveTo(start.x, start.y - 5 + snakeWidth + 10);
        ctx.lineTo(start.x + 2, start.y - 5 + snakeWidth + 7);
        ctx.strokeStyle = 'red';
        ctx.lineWidth = 2;
        ctx.stroke();
    };

    const drawPlayerPositions = (ctx: CanvasRenderingContext2D) => {
        // Draw each player's token at their current position
        gameState.players.forEach((player, index) => {
            drawPlayerToken(ctx, player.position, index);
        });
    };

    const drawPlayerToken = (ctx: CanvasRenderingContext2D, position: number, playerIndex: number) => {

        let x, y;

        // initial position (starting area at bottom)
        if (position == 0) {
            x = canvasWidth / 2; // Center horizontally in starting area
            y = boardSize * cellSize + startingAreaHeight / 2; // Position in bottom starting area
        } else {
            const cellPos = getCellPosition(position);
            x = cellPos.x;
            y = cellPos.y;
        }

        const tokenX = x + (playerIndex % 2) * 20 + 10;
        const tokenY = y + Math.floor(playerIndex / 2) * 20 + 10;
        const radius = 12;

        const playerImageSrc = gameState.players[playerIndex].picture; // Base64 string
        const img = new Image();
        img.src = playerImageSrc;

        img.onload = () => {
            // Draw glow effect (optional)
            ctx.shadowColor = '#FFFFFF'; // or any glow color
            ctx.shadowBlur = 15;
            ctx.beginPath();
            ctx.arc(tokenX, tokenY, radius + 2, 0, 2 * Math.PI);
            ctx.fill();

            // Reset shadow for image
            ctx.shadowBlur = 0;

            // Draw the image clipped in a circle
            ctx.save();
            ctx.beginPath();
            ctx.arc(tokenX, tokenY, radius, 0, 2 * Math.PI);
            ctx.closePath();
            ctx.clip();

            ctx.drawImage(img, tokenX - radius, tokenY - radius, radius * 2, radius * 2);
            ctx.restore();

            // Draw border around the image
            ctx.strokeStyle = '#FFFFFF';
            ctx.lineWidth = 2;
            ctx.beginPath();
            ctx.arc(tokenX, tokenY, radius, 0, 2 * Math.PI);
            ctx.stroke();
        };
    };


    const getCellPosition = (position: number): { x: number, y: number } => {
        if (position <= 0 || position > gameState.game_board.total_cells) {
            return { x: 0, y: 0 };
        }

        const row = boardSize - Math.ceil(position / boardSize);
        let col;

        // Determine if row is left-to-right or right-to-left
        if (Math.ceil(position / boardSize) % 2 === 1) {
            // Left to right
            col = (position - 1) % boardSize;
        } else {
            // Right to left
            col = boardSize - 1 - ((position - 1) % boardSize);
        }

        return {
            x: (col * cellSize + cellSize / 2),
            y: (row * cellSize + cellSize / 2) - 7 // Remove starting area offset
        };
    };

    const drawStartingArea = (ctx: CanvasRenderingContext2D) => {
        // Position starting area at the bottom
        const startY = boardSize * cellSize;

        // Draw starting area background
        const gradient = ctx.createLinearGradient(0, startY, canvasWidth, startY + startingAreaHeight);
        gradient.addColorStop(0, '#ffffff'); // Peach
        gradient.addColorStop(1, '#fff8e6'); // Light pink

        ctx.fillStyle = gradient;
        ctx.fillRect(0, startY, canvasWidth, startingAreaHeight);

        // Add border
        ctx.strokeStyle = '#000000';
        ctx.lineWidth = 2;
        ctx.strokeRect(0, startY, canvasWidth, startingAreaHeight);

        // Add "START" text
        ctx.fillStyle = '#2C3E50';
        ctx.font = 'bold 24px Arial';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText('START', canvasWidth / 2, startY + startingAreaHeight / 2);
    };

    return (
        <div className="flex flex-col items-center justify-center">
            <canvas
                ref={canvasRef}
                width={canvasWidth}
                height={canvasHeight}
                className="border-2 border-gray-300"
            />
        </div>
    );
}