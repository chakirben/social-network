"use client";

import SideBar from "@/components/sidebar";
import { ws } from "../websocket/WebSocketProvider";

export default function Chat() {
    const handleSendMessage = () => {
        console.log("Message sent!");
        const msg = {
            type: "chat",
            data: {
                chat: 1,
                content: "Hello, this is a test message!",
            },
        };
        ws.send(JSON.stringify(msg));
    };

    return (
        <div>
            <SideBar />
            <main style={{ padding: "1rem" }}>
                <div style={{ border: "1px solid #ccc", borderRadius: "8px", padding: "16px", maxWidth: "600px", margin: "0 auto" }}>
                    <div style={{ minHeight: "200px", marginBottom: "16px" }}>
                        {/* Chat messages will appear here */}
                    </div>
                    <form style={{ display: "flex", gap: "8px" }}>
                        <input
                            type="text"
                            placeholder="Type your message..."
                            style={{ flex: 1, padding: "8px", borderRadius: "4px", border: "1px solid #ccc" }}
                        />
                        <Sendmsg onClick={handleSendMessage}>
                            <button
                                type="button"
                                style={{ padding: "8px 16px", borderRadius: "4px", background: "#0070f3", color: "#fff", border: "none" }}
                            >
                                Send
                            </button>
                        </Sendmsg>
                    </form>
                </div>
            </main>
        </div>
    );
}

function Sendmsg({ onClick, children }) {
    return (
        <div onClick={onClick}>
            {children}
        </div>
    );
}
