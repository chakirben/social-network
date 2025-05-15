"use client"
import { useEffect, useState } from "react"
import Groupbar from "@/components/groups/groupbar"
import SideBar from "@/components/sidebar"
import DataToCreatGroup from "@/components/groups/creatGroup"
import FetchCreatGroup from "./api_and_funcs/fetch_creat_gp"
import MyGroupsPage from "./api_and_funcs/get_groups"

import "./css/groups1.css"
import "./css/creatgroup.css"
import "./../home/home.css"

export default function JustMyGroupsPage() {
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [myGroups, setMyGroups] = useState([]);
  const [notMyGroups, setNotMyGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true); 

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const [myGroupsRes, notMyGroupsRes] = await Promise.all([
          fetch("http://localhost:8080/api/MyGroups", { credentials: "include" }),
          fetch("http://localhost:8080/api/NotMyGroups", { credentials: "include" }),
        ]);
        const myGroupsData = await myGroupsRes.json();
        const notMyGroupsData = await notMyGroupsRes.json();
        setMyGroups(myGroupsData || []);
        setNotMyGroups(notMyGroupsData || []);
      } catch (error) {
        console.error("Error fetching groups:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchGroups();
  }, []);

  const handleBTNCreateGroup = () => {
    setShowCreateGroup(true);
  }

  const creatgroupGroup = async (groupData) => {
    const newGroup = await FetchCreatGroup(groupData.title, groupData.description);
    if (newGroup) {
      setMyGroups((prev) => [...prev, newGroup]);
      setShowCreateGroup(false);
    }
  }


  if (isLoading) {
    return <div className="loading">Loading groups...</div>;
  }

 
  if (myGroups.length === 0 && notMyGroups.length === 0) {
    return (
      <div className="home">
        <SideBar />
        <div className="divallGroups">
          <Groupbar onCreateGroup={handleBTNCreateGroup} />
          <div>No groups found...</div>
          {showCreateGroup && (
            <DataToCreatGroup
              onCreate={creatgroupGroup}
              onSkip={() => setShowCreateGroup(false)}
            />
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">

        <Groupbar onCreateGroup={handleBTNCreateGroup} />

        {showCreateGroup &&
          <DataToCreatGroup
            onCreate={creatgroupGroup}
            onSkip={() => setShowCreateGroup(false)}
          />
        }

        <MyGroupsPage
          myGroups={myGroups}
          notMyGroups={notMyGroups}
        />
      </div>
    </div>
  );
}
