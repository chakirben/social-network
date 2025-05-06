import Groupbar from "@/components/groupbar"
import MyGroupsPage from "./groups"
import SideBar from "@/components/sidebar";
import "./groups.css"

export default function JustMyGroupsPage() {
    return (
        <div className="home">
            <SideBar />
            <div className="divallGroups">
              <Groupbar />
              <div className="groupsmn">
                <MyGroupsPage />
              </div>
            </div>
        </div>
    );
}