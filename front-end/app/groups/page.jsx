"use client";

import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import CreateGroupForm from "@/components/groups/creatGroup";
import fetchCreateGroup from "./api_and_funcs/fetch_creat_gp";
import MyGroupsPage from "./api_and_funcs/get_groups";

import "./css/groups1.css";
import "./css/creatgroup.css";
import "./../home/home.css";
import Header from "@/components/header";

export default function JustMyGroupsPage() {
  const [showCreateGroupForm, setShowCreateGroupForm] = useState(false);
  const [myGroups, setMyGroups] = useState([]);
  const [otherGroups, setOtherGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const [myRes, otherRes] = await Promise.all([
          fetch("http://localhost:8080/api/MyGroups", { credentials: "include" }),
          fetch("http://localhost:8080/api/NotMyGroups", { credentials: "include" }),
        ]);
        const myData = await myRes.json();
        const otherData = await otherRes.json();
        setMyGroups(myData || []);
        setOtherGroups(otherData || []);
      } catch (error) {
        console.error("Error fetching groups:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchGroups();
  }, []);

  const handleCreateGroupClick = () => {
    setShowCreateGroupForm(true);
  };

  const handleCreateGroup = async (groupData) => {
    const newGroup = await fetchCreateGroup(groupData.title, groupData.description);
    if (newGroup) {
      setMyGroups((prev) => [...prev, newGroup]);
      setShowCreateGroupForm(false);
    }
  };

  const renderGroupContent = () => {
    if (myGroups.length === 0 && otherGroups.length === 0) {
      return <div>No groups found...</div>;
    }

    return (
      <MyGroupsPage
        myGroups={myGroups}
        notMyGroups={otherGroups}
      />
    );
  };

  return (
    <div className="home">
      <SideBar />
      <div className="divallGroups">
        <Header
          pageName="Groups"
          ele={
            <button onClick={handleCreateGroupClick} className="create-group-btn">
              + Create Group
            </button>
          }
        />

        {showCreateGroupForm && (
    
            <CreateGroupForm onCreate={handleCreateGroup} onSkip={() => setShowCreateGroupForm(false)} />
        
        )}

        {renderGroupContent()}

      </div>
    </div>
  )
}