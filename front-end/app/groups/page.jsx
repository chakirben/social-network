"use client"
import { useState } from "react"
import Groupbar from "@/components/groups/groupbar"
import MyGroupsPage from "./api_and_funcs/get_groups"
import SideBar from "@/components/sidebar";
import DataToCreatGroup from "./api_and_funcs/creat_group"
import "./css/groups1.css"
import "./css/creatgroup.css"

export default function JustMyGroupsPage() {
  const [showCreateGroup, setShowCreateGroup] = useState(false)

  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true)
  }
  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Groupbar onCreateGroup={handleBTNCreateGroup} />
        {showCreateGroup && <DataToCreatGroup />}
          <div className="groupsmn">
          <MyGroupsPage />
        </div>
      </div>
    </div>
  );
}
