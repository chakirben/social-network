import MyGroupsPage from "./mygroups";
import Groupbar from "@/components/groupbar"
import NotMyGroupsPage from "./notmygroups"
import SideBar from "@/components/sidebar";
import "./groups.css"

export default function JustMyGroupsPage() {
    return (
        <div className="home">
            <SideBar />
            <div className="divallGroups">
              <Groupbar />
              <div className="groupsmn">
                <NotMyGroupsPage />
                <MyGroupsPage />
              </div>
            </div>
        </div>
    );
}