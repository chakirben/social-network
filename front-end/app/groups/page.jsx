"use client"
import { useState } from "react"
import Groupbar from "@/components/groups/groupbar"
import MyGroupsPage from "./api_and_funcs/get_groups"
import SideBar from "@/components/sidebar";
import DataToCreatGroup from "@/components/groups/creatGroup"
import FetchCreatGroup from "./api_and_funcs/fetch_creat_gp"
import FetchJoinToGroup from "./api_and_funcs/fetch_req_join_gp"
import "./css/groups1.css"
import "./css/creatgroup.css"

export default function JustMyGroupsPage() {
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [isGroupCreated, setIsGroupCreated] = useState(false);
  const [viewingGroup, setViewingGroup] = useState(null);

  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true)
  }
  const creatgroupGroup = (groupData) => {
    FetchCreatGroup(groupData.title, groupData.description)
    setIsGroupCreated(true);
    setShowCreateGroup(false);
  }

  const Jointogroup = (groupId) => {
    FetchJoinToGroup(groupId)
  }

  const viewGroup = (groupId) => {
    setViewingGroup(groupId)
  }

  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Groupbar onCreateGroup={handleBTNCreateGroup} />
        {showCreateGroup && <DataToCreatGroup onCreate={creatgroupGroup} onSkip={() => setShowCreateGroup(false)} />}

        {!viewingGroup ? (
            <MyGroupsPage onJoin={Jointogroup} onView={viewGroup} />
        ) : (
          // <GroupDetails group={viewingGroup} onBack={() => setViewingGroup(null)} />
          <div>
            hiiiiiiiiiiiiiii
          </div>
        )}  

      </div>
    </div>
  );
}
