import SideBar from "./sidebar";
import './globals.css';

export default function Router({ children }) {
    return (
        <html lang="en">
            <body>
                <div className="layout">
                    <SideBar />
                    <main className="mainContent">{children}</main>
                </div>
            </body>
        </html>
    )
}