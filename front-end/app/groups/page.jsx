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
  const [isGroupCreated, setIsGroupCreated] = useState(false);

  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true)
  }
  const creatgroupGroup = (groupData) => {
    console.log("Group created:", groupData)
    setIsGroupCreated(true);
    setShowCreateGroup(false);
  };
  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Groupbar onCreateGroup={handleBTNCreateGroup} />
        {showCreateGroup && <DataToCreatGroup onCreate={creatgroupGroup}  />}
          <div className="groupsmn">
          <MyGroupsPage />
        </div>
      </div>
    </div>
  );
}
