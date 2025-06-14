import { UserProvider } from "@/components/context/userContext";
import "../styles/global.css";
import WebSocketProvider from "./websocket/WebSocketProvider";
import NotifPopUp from "@/components/notificationPopUp";

export default function Main({ children }) {

    return (
        <html lang="en">
            <body>
                {/* <NotifPopUp /> */}
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
