import { useState } from "react";
import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import {
    Menubar,
    MenubarContent,
    MenubarItem,
    MenubarMenu,
    MenubarShortcut,
    MenubarTrigger,
} from "@/components/ui/menubar"
import { useUser } from "@/hooks/useUser";
import { useNavigate } from "react-router";

export default function Navbar() {
    const [darkMode, setDarkMode] = useState(false);
    const navigate = useNavigate();

    const toggleDarkMode = () => {
        setDarkMode(!darkMode);
        document.documentElement.classList.toggle("dark", !darkMode);
    };

    const signOut = () => { 
        localStorage.clear();
        navigate("/login");
    }

    const {
        data,
        isError,
        error,
    } = useUser();

    if (isError) {
        console.error("Error fetching user data:", error);
    }

    return (
        <nav className="flex justify-between items-center p-4">
            <div className="flex items-center space-x-4">
                <img src="/logo.gif" alt="Logo" className="h-12" />
            </div>

            <div className="flex flex-row items-center">
                <Button onClick={toggleDarkMode} variant="ghost" >
                    {darkMode ? (
                        <Moon className="size-6" />
                    ) : (
                        <Sun className="size-6" />
                    )}
                </Button>
                <Menubar>
                    <MenubarMenu>
                        <MenubarTrigger>
                            <Avatar>
                                <AvatarImage src={data?.picture} alt="user" />
                                <AvatarFallback>{data?.name.charAt(0).toUpperCase()}</AvatarFallback>
                            </Avatar>
                        </MenubarTrigger>
                        <MenubarContent>
                            <MenubarItem onClick={signOut}>
                                Sign Out <MenubarShortcut>⌘S</MenubarShortcut>
                            </MenubarItem>
                        </MenubarContent>
                    </MenubarMenu>
                </Menubar>
            </div>
        </nav>
    );
}
