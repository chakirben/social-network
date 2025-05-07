"use client"
import Groupbar from "@/components/groupbar"
import MyGroupsPage from "./api_and_funcs/get_groups"
import SideBar from "@/components/sidebar";
import "./css/groups1.css"

export default function JustMyGroupsPage() {
    const handleBTNCreateGroup = () => {
      DataToCreatGroup()
    }
    return (
        <div className="home">
            <SideBar />
            <div className="divallGroups">
            <Groupbar onCreateGroup={handleBTNCreateGroup} />
              <div className="groupsmn">
                <MyGroupsPage />
              </div>
            </div>
        </div>
    );
}
