"use client"
import { useState } from "react"
import Groupbar from "@/components/groups/groupbar"
import MyGroupsPage from "./api_and_funcs/get_groups"
import SideBar from "@/components/sidebar";
import DataToCreatGroup from "@/components/groups/creatGroup"
import FetchCreatGroup from "./api_and_funcs/fetch_creat_gp"
import "./css/groups1.css"
import "./css/creatgroup.css"

export default function JustMyGroupsPage() {
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [isGroupCreated, setIsGroupCreated] = useState(false);

  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true)
  }
  const creatgroupGroup = (groupData) => {
    FetchCreatGroup(groupData.title , groupData.description)
    setIsGroupCreated(true);
    setShowCreateGroup(false);
  };
  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Groupbar onCreateGroup={handleBTNCreateGroup} />
        {showCreateGroup && <DataToCreatGroup onCreate={creatgroupGroup}  onSkip={() => setShowCreateGroup(false)} />}
          <div className="groupsmn">
          <MyGroupsPage />
        </div>
      </div>
    </div>
  );
}
