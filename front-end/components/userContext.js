"use client";
import { createContext, useContext, useState, useEffect } from "react";

const UserContext = createContext();

export const UserProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const savedUser = localStorage.getItem("user");
        if (savedUser) {
            setUser(JSON.parse(savedUser));
            setLoading(false);
        } else {
            fetch("http://localhost:8080/api/getUserData",{ credentials: "include" })
                .then((response) => response.json())
                .then((data) => {
                    setUser(data);
                    localStorage.setItem("user", JSON.stringify(data));
                })
                .finally(() => {
                    setLoading(false);
                });//GetCurrentUserData
        }
    }, []);

    return (
        <UserContext.Provider value={{ user, setUser }}>
            {!loading && children}
        </UserContext.Provider>
    );
};

export const useUser = () => useContext(UserContext);
