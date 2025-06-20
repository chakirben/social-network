"use client";
import { createContext, useContext, useState, useEffect } from "react";

const UserContext = createContext();

export const UserProvider = ({ children }) => {
    console.log("UserProvider rendered");

    const [user, setUser] = useState(null);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const rep = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getUserData`, {
                    credentials: "include",
                });

                if (!rep.ok) {
                    throw new Error(`HTTP error! Status: ${rep.status}`);
                }

                const data = await rep.json();
                setUser(data);
            } catch (err) {
                console.log("Error fetching user:", err);
            }

        };
        fetchUser();
    }, []);
    return (
        <UserContext.Provider value={{ user, setUser }}>
            {children}
        </UserContext.Provider>
    );
};

export const useUser = () => useContext(UserContext);


