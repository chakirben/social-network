"use client";
import { useEffect } from "react";
import InitWs from "./websocket";
export let ws;
export default function WebSocketProvider() {

    useEffect(() => {

     InitWs();
    }, []);

}