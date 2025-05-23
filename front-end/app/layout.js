import { UserProvider } from "@/components/userContext";
import "../styles/global.css";
import InitWs from "./websocket/websocket";

export default async function Main({ children }) {
    if (!InitWs()) {
        ws = await InitWs();
        ws.send()
    }
    return (
        <html lang="en">
            <body>
                <UserProvider>
                    <div className="layout">
                        <main className="mainContent">{children}</main>
                    </div>
                </UserProvider>
            </body>
        </html>
    );
}
