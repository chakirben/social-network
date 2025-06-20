"use client";

import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import CreateGroupForm from "@/components/groups/creatGroup";
import fetchCreateGroup from "./api_and_funcs/fetch_creat_gp";
import MyGroupsPage from "./api_and_funcs/get_groups";

import "./css/groups1.css";
import "./css/creatgroup.css";
import "./../home/home.css";
import Header from "@/components/Header/header";

export default function JustMyGroupsPage() {
  const [showCreateGroupForm, setShowCreateGroupForm] = useState(false);
  const [myGroups, setMyGroups] = useState([]);
  const [otherGroups, setOtherGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/MyGroups`, {
          credentials: "include",
        });
        const data = await res.json();
        setMyGroups(data || []);
      } catch (error) {
        console.error("Error fetching groups:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchGroups();
    const fetchNotGroups = async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/NotMyGroups`, {
          credentials: "include",
        });
        const data = await res.json();
        setOtherGroups(data || []);
      } catch (error) {
        console.error("Error fetching groups:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchNotGroups()
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
      return <div className="msgNoGroups">No groups found... ): creat a one (:</div>;
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
              + Create
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