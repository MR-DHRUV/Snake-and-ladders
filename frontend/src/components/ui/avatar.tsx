"use client"

import * as React from "react"
import * as AvatarPrimitive from "@radix-ui/react-avatar"

import { cn } from "@/lib/utils"

function Avatar({
    className,
    ...props
}: React.ComponentProps<typeof AvatarPrimitive.Root>) {
    return (
        <AvatarPrimitive.Root
            data-slot="avatar"
            className={cn(
                "relative flex size-8 shrink-0 overflow-hidden rounded-full",
                className
            )}
            {...props}
        />
    )
}

function AvatarImage({
    className,
    ...props
}: React.ComponentProps<typeof AvatarPrimitive.Image>) {
    return (
        <AvatarPrimitive.Image
            data-slot="avatar-image"
            className={cn("aspect-square size-full", className)}
            {...props}
        />
    )
}

function AvatarFallback({
    className,
    ...props
}: React.ComponentProps<typeof AvatarPrimitive.Fallback>) {
    const colors = [
        "bg-red-100",
        "bg-blue-100",
        "bg-green-100",
        "bg-yellow-100",
        "bg-purple-100",
        "bg-pink-100",
        "bg-orange-100",
        "bg-teal-100",
        "bg-indigo-100",
        "bg-lime-100",
        "bg-amber-100",
        "bg-cyan-100",
        "bg-emerald-100",
        "bg-rose-100",
        "bg-violet-100",
    ];

    const randomColor = colors[Math.floor(Math.random() * colors.length)];

return (
    <AvatarPrimitive.Fallback
        data-slot="avatar-fallback"
        className={cn(
            "flex size-full items-center justify-center rounded-full grayscale-0 text-black",
            randomColor,
            className
        )}
        {...props}
    />
)}

export { Avatar, AvatarImage, AvatarFallback }
