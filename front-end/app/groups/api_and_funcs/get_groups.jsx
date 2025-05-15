"use client";
import MyGroup from "@/components/groups/mygroups";
import NoMyGroup from "@/components/groups/notmygroups";

export default function MyGroupsPage({ myGroups, notMyGroups, onJoin }) {
  if (myGroups.length === 0 && notMyGroups.length === 0) {
    return <div>No groups found...</div>;
  }

  return (
    <div className="groupsmn">
      {myGroups.map((group) => (
        <MyGroup key={group.Id} group={group} />
      ))}
      {notMyGroups.map((group) => (
        <NoMyGroup key={group.Id} group={group} onJoin={onJoin} />
      ))}
    </div>
  );
}
