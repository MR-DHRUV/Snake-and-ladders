import { useState, useEffect, useRef } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Send } from "lucide-react";
import type { ChatMessage } from "@/types/chat";

interface ChatContainerProps {
    chatMessages: ChatMessage[];
    currentUser?: string;
    sendChatMessage: (message: string) => void;
}

export default function ChatContainer({
    chatMessages,
    currentUser,
    sendChatMessage,
}: ChatContainerProps) {
    const [newMessage, setNewMessage] = useState("");
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const handleSend = () => {
        if (!newMessage.trim()) return;
        sendChatMessage(newMessage);
        setNewMessage("");
    };

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [chatMessages]);

    return (
        <div className="flex flex-col w-full max-w-md rounded-lg border">
            {/* Chat messages */}
            <div className="flex-1 p-4 space-y-2 overflow-y-auto max-h-[710px]">
                {chatMessages.map((msg, idx) => {
                    const isCurrentUser = msg.player_id === currentUser;
                    const isSystem = msg.player_name === "System";
                    return (
                        <div
                            key={idx}
                            className={`flex ${isCurrentUser ? "justify-end" : "justify-start"}`}
                        >
                            <div
                                className={
                                    isSystem
                                        ? "text-[13px] italic text-gray-500 w-full text-center"
                                        : isCurrentUser
                                            ? "max-w-[70%] px-3 py-2 rounded-2xl text-sm bg-amber-400 text-white rounded-br-none"
                                            : "max-w-[70%] px-3 py-2 rounded-2xl text-sm bg-gray-200 text-gray-900 rounded-bl-none"
                                }
                            >
                                {!isCurrentUser && !isSystem && (
                                    <div className="text-xs font-medium text-gray-600 mb-1">
                                        {msg.player_name}
                                    </div>
                                )}
                                {msg.message}
                            </div>
                        </div>
                    );
                })}

                {/* Scroll Anchor */}
                <div ref={messagesEndRef} className="h-[1px]" />
            </div>

            {/* Input + Send button */}
            <div className="flex items-center gap-2 p-6">
                <Input
                    placeholder="Type a message..."
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleSend()}
                />
                <Button onClick={handleSend} size="icon" className="rounded-full">
                    <Send className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}
