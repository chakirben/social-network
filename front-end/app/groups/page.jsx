"use client"
import { useState } from "react"
import Groupbar from "@/components/groups/groupbar"
import MyGroupsPage from "./api_and_funcs/get_groups"
import SideBar from "@/components/sidebar";
import DataToCreatGroup from "@/components/groups/creatGroup"
import FetchCreatGroup from "./api_and_funcs/fetch_creat_gp"
import FetchJoinToGroup from "./api_and_funcs/fetch_req_join_gp"
import DisplyGroup from "./[id]/page"
import "./css/groups1.css"
import "./css/creatgroup.css"
import "./../home/home.css"

export default function JustMyGroupsPage() {
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [isGroupCreated, setIsGroupCreated] = useState(false)
  const [viewingGroup, setViewingGroup] = useState(null)

  // For the button to shoz the div of creat group...
  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true)
  }

  // For the button to creat group...
  const creatgroupGroup = (groupData) => {
    FetchCreatGroup(groupData.title, groupData.description)
    setIsGroupCreated(true);
    setShowCreateGroup(false);
  }

  // For the button to JOIN group...
  const Jointogroup = (groupId) => {
    FetchJoinToGroup(groupId)
  }

  // For the button to VIEW group...
  const viewGroup = (groupId) => {
    console.log("hawa", groupId);

    setViewingGroup(groupId)
  }

  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Groupbar onCreateGroup={handleBTNCreateGroup} />
        {showCreateGroup && <DataToCreatGroup onCreate={creatgroupGroup} onSkip={() => setShowCreateGroup(false)} />}

        <MyGroupsPage onJoin={Jointogroup} onView={viewGroup} />
        {viewingGroup && (
          <DisplyGroup groupId={viewingGroup} back={() => setViewingGroup(null)} />
        )}


      </div>
    </div>
  );
}
