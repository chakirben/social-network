import { UserProvider } from "@/components/userContext";
import "../styles/global.css";
import WebSocketProvider from "./websocket/WebSocketProvider";

export default function Main({ children }) {

    return (
        <html lang="en">
            <body>
                <UserProvider>
                    <WebSocketProvider />
                        <div className="layout">
                            <main className="mainContent">{children}</main>
                        </div>
                </UserProvider>
            </body>
        </html>
    );
}
