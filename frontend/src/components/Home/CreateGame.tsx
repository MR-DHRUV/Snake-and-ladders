/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState } from "react";
import { Dialog, DialogContent, DialogClose } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { createGame } from "./createGameRequest";
import { z } from "zod";
import useAxios from "@/hooks/useAxios";
import { useNavigate } from "react-router";

const formSchema = z.object({
    maxPlayers: z.number().min(1, { message: "Max players must be between 1 and 10." }).max(10, { message: "Max players must be between 1 and 10." }),
    maxWinners: z.number().min(1, { message: "Max winners must be between 1 and 2." }).max(2, { message: "Max winners must be between 1 and 2." }),
    diceCount: z.number().min(1, { message: "Dice count must be between 1 and 2." }).max(2, { message: "Dice count must be between 1 and 2." }),
});

type FormData = z.infer<typeof formSchema>;

const defaultFormData: FormData = {
    maxPlayers: 2,
    maxWinners: 1,
    diceCount: 1,
};

const CreateGameModal = (
    { open, setOpen }: {
        open: boolean;
        setOpen: (open: boolean) => void;
    }
) => {
    const [formData, setFormData] = useState<FormData>(defaultFormData);
    const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({});
    const [_submitting, setSubmitting] = useState(false);
    const axiosInstance = useAxios();
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: parseInt(value, 10),
        }));
    };

    const validate = (data: FormData) => {
        const result = formSchema.safeParse(data);
        if (result.success) {
            setErrors({});
            return true;
        } else {
            const fieldErrors: Partial<Record<keyof FormData, string>> = {};
            result.error.errors.forEach((err) => {
                fieldErrors[err.path[0] as keyof FormData] = err.message;
            });
            setErrors(fieldErrors);
            return false;
        }
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (validate(formData)) {
            resetForm();
            setSubmitting(true);

            createGame({
                max_players: formData.maxPlayers,
                max_winners: formData.maxWinners,
                dice_count: formData.diceCount,
                axiosInsatnce: axiosInstance,
            })
                .then((response) => {
                    setSubmitting(false);
                    setOpen(false);
                    navigate(`/game/${response.data.data}`);
                })
        }
    };

    const resetForm = () => {
        setFormData(defaultFormData);
        setErrors({});
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <div className="flex flex-col mb-3">
                    <h2 className="text-xl font-semibold">Game Settings</h2>
                    <p >
                        Create your own game with custom settings.
                    </p>
                </div>
                <form onSubmit={handleSubmit}>
                    <div className="space-y-4">
                        <div className="flex flex-col">
                            <Label htmlFor="maxPlayers">Max Players</Label>
                            <Input
                                type="number"
                                id="maxPlayers"
                                name="maxPlayers"
                                value={formData.maxPlayers}
                                onChange={handleChange}
                                placeholder="Max players (1-10)"
                                className="mt-1"
                            />
                            {errors.maxPlayers && <p className="text-sm text-red-500">{errors.maxPlayers}</p>}
                        </div>

                        <div className="flex flex-col">
                            <Label htmlFor="maxWinners">Max Winners</Label>
                            <Input
                                type="number"
                                id="maxWinners"
                                name="maxWinners"
                                value={formData.maxWinners}
                                onChange={handleChange}
                                placeholder="Max winners (1-2)"
                                className="mt-1"
                            />
                            {errors.maxWinners && <p className="text-sm text-red-500">{errors.maxWinners}</p>}
                        </div>

                        <div className="flex flex-col">
                            <Label htmlFor="diceCount">Dice Count</Label>
                            <Input
                                type="number"
                                id="diceCount"
                                name="diceCount"
                                value={formData.diceCount}
                                onChange={handleChange}
                                placeholder="Dice count (1-2)"
                                className="mt-1"
                            />
                            {errors.diceCount && <p className="text-sm text-red-500">{errors.diceCount}</p>}
                        </div>

                        <div className="mt-6 gap-3 flex justify-end">
                            <DialogClose asChild>
                                <Button variant="secondary" onClick={resetForm}>
                                    Cancel
                                </Button>
                            </DialogClose>
                            <Button type="submit" variant="default">
                                Create Game
                            </Button>
                        </div>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateGameModal;
