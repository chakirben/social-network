import { UserProvider } from "@/components/context/userContext";
import "../styles/global.css";
import { WebSocketProvider } from "@/components/context/wsContext";
import { PopupProvider } from "@/components/context/popUp";

export default function Layout({ children }) {
    console.log("Main layout rendered");
    return (
        <html lang="en">
            <body>
                <UserProvider>
                    <PopupProvider>
                        <WebSocketProvider>
                            <div className="layout">
                                <main className="mainContent">{children}</main>
                            </div>
                        </WebSocketProvider>
                    </PopupProvider>
                </UserProvider>
            </body>
        </html>
    );
}
