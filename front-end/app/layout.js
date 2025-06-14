import { UserProvider } from "@/components/context/userContext";
import "../styles/global.css";
import NotifPopUp from "@/components/notificationPopUp";
import { WebSocketProvider } from "@/components/context/wsContext";

export default function Main({ children }) {

    return (
        <html lang="en">
            <body>
                <UserProvider>
                    <WebSocketProvider >
                    <div className="layout">
                        <main className="mainContent">{children}</main>
                    </div>
                    
                    </WebSocketProvider>
                </UserProvider>
            </body>
        </html>
    );
}
