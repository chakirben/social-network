"use client";
import MyGroup from "@/components/groups/mygroups";
import NoMyGroup from "@/components/groups/notmygroups";
import { useState } from "react";

export default function MyGroupsPage({ myGroups, notMyGroups, onJoin }) {
  if (myGroups.length === 0 && notMyGroups.length === 0) {
    return <div>No groups found...</div>;
  }

  const [ActiveTab ,  setActiveTab] = useState("mygroups")

  return (
    <>
    <div className='filterGroups'>
 
             <span className="MyGroupsSpam" 
              onClick={() => setActiveTab("mygroups")}
             >
                    MyGroups
                </span>
                <span className="GroupsNotJoinYetSpam">
                  NotJoinYet
                </span>
      
    </div>
    <div className="groupsmn">
      {ActiveTab === "mygroups" && (
        <>
          {myGroups.map((group) => (
             <MyGroup key={group.Id} group={group} />
          ))}
        </>

      ) 
      }
      
      {notMyGroups.map((group) => (
        <NoMyGroup key={group.Id} group={group} onJoin={onJoin} />
      ))}
    </div>
    </>
  );
}
