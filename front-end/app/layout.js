// app/layout.js (or app/mainLayout.js)
import SideBar from "@/components/sidebar";
import "../styles/global.css"
export default function Main({ children }) {
    return (
        <html lang="en">
            <body>
                <div className="layout">
                    <main className="mainContent">{children}</main>
                </div>
            </body>
        </html>
    );
}
