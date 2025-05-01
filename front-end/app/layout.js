import SideBar from "./sidebar";
import './globals.css';

export default function Router({ children }) {
    return (
        <html lang="en">
            <body>
                <SideBar />
                <main>{children}</main>
            </body>
        </html>
    )
}