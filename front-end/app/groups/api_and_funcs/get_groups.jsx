"use client";
import MyGroup from "@/components/groups/mygroups";
import NoMyGroup from "@/components/groups/notmygroups";
import { useState } from "react";

export default function MyGroupsPage({ myGroups, notMyGroups, onJoin }) {
  if (myGroups.length === 0 && notMyGroups.length === 0) {
      return <div>No groups found... ): creat a one (:</div>;
  }

  const [ActiveTab, setActiveTab] = useState("mygroups")

  return (
    <>
      <div className='filterGroups'>

        <span className={`MyGroupsSpam ${ActiveTab === "mygroups" ? "activeGroups" : ""}`}
          onClick={() => setActiveTab("mygroups")
          }>
          MyGroups
        </span>
        <span className={`GroupsNotJoinYetSpam ${ActiveTab === "NotJoinYet" ? "activeGroups" : ""}`}
          onClick={() => setActiveTab("NotJoinYet")
          }>
          NotJoinYet
        </span>

      </div>

      <div className="groupsmn">
        {ActiveTab === "mygroups" ? (
          <>
            {myGroups.length == 0 ? (
              <div className="msgNoGroups">
                There is no group for you ): creat a one... (:
              </div>
            ) : (
              <>
                {myGroups.map((group) => (
                  <MyGroup key={group.Id} group={group} />
                ))}
              </>
            )}
          </>

        ) : (
          <>
            {notMyGroups.length == 0 ? (
              <div className="msgNoGroups">
                There are no groups to join ): Create one... (:
              </div>
            ) : (
              <>
                {notMyGroups.map((group) => (
                  <NoMyGroup key={group.Id} group={group} onJoin={onJoin} />
                ))}
              </>
            )}
          </>
        )}


      </div >
    </>
  );
}
