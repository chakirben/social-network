import { UserProvider } from "@/components/userContext";
import "../styles/global.css";

export default function Main({ children }) {
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
